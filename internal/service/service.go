package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"EventBooker/internal/model"

	"github.com/google/uuid"
)

// BookerRepository — интерфейс репозитория, который передается сервису.
type BookerRepository interface {
	// Users
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)

	// Events
	CreateEvent(ctx context.Context, event *model.Event) error
	GetEventByID(ctx context.Context, id uuid.UUID) (*model.Event, error)
	ListEvents(ctx context.Context) ([]*model.Event, error)

	// Bookings
	CreateBooking(ctx context.Context, booking *model.Booking) error
	GetBookingByID(ctx context.Context, id uuid.UUID) (*model.Booking, error)
	UpdateBookingStatus(ctx context.Context, bookingID uuid.UUID, status model.BookingStatus) error
	ListBookingsByEventID(ctx context.Context, eventID uuid.UUID) ([]*model.Booking, error)
	CancelExpiredBookings(ctx context.Context) (int64, error)
}

type Service struct {
	repo BookerRepository
}

func New(repo BookerRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RegisterUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	if req.Email == "" {
		return nil, errors.New("email is required")
	}

	existingUser, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, model.ErrUserNotFound) {
		return nil, fmt.Errorf("check email: %w", err)
	}

	if existingUser != nil {
		return nil, model.ErrEmailAlreadyTaken
	}

	user := &model.User{
		ID:          uuid.New(),
		Email:       req.Email,
		TelegramUse: req.TelegramUse,
		CreatedAt:   time.Now().UTC(),
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (s *Service) GetUserProfile(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, model.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user profile: %w", err)
	}

	return user, nil
}

func (s *Service) CreateEvent(ctx context.Context, req *model.CreateEventRequest) (*model.Event, error) {
	if req.EventDate.Before(time.Now()) {
		return nil, model.ErrEventExpired
	}

	if req.TotalSeats <= 0 {
		return nil, model.ErrTotalSeats
	}

	ttl := req.BookingTTLMinutes
	if ttl <= 0 {
		ttl = 15
	}

	event := &model.Event{
		ID:                uuid.New(),
		Title:             req.Title,
		EventDate:         req.EventDate,
		TotalSeats:        req.TotalSeats,
		AvailableSeats:    req.TotalSeats,
		BookingTTLMinutes: ttl,
		CreatedAt:         time.Now().UTC(),
	}

	if err := s.repo.CreateEvent(ctx, event); err != nil {
		return nil, fmt.Errorf("create event: %w", err)
	}

	return event, nil
}

func (s *Service) GetEventDetails(ctx context.Context, eventID uuid.UUID) (*model.EventDetailResponse, error) {
	event, err := s.repo.GetEventByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	bookings, err := s.repo.ListBookingsByEventID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("list bookings for event: %w", err)
	}

	var bookingList []model.Booking
	for _, b := range bookings {
		if b != nil {
			bookingList = append(bookingList, *b)
		}
	}

	return &model.EventDetailResponse{
		Event:    *event,
		Bookings: bookingList,
	}, nil
}

func (s *Service) ListUpcomingEvents(ctx context.Context) ([]*model.Event, error) {
	events, err := s.repo.ListEvents(ctx)
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}

	return events, nil
}

func (s *Service) BookSeat(ctx context.Context, eventID uuid.UUID, req *model.BookSeatRequest) (*model.BookSeatResponse, error) {
	if _, err := s.repo.GetUserByID(ctx, req.UserID); err != nil {
		return nil, err
	}

	event, err := s.repo.GetEventByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if event.EventDate.Before(time.Now()) {
		return nil, model.ErrEventExpired
	}

	if event.AvailableSeats <= 0 {
		return nil, model.ErrNoAvailableSeats
	}

	expiresAt := time.Now().UTC().Add(time.Duration(event.BookingTTLMinutes) * time.Minute)

	booking := &model.Booking{
		ID:        uuid.New(),
		EventID:   eventID,
		UserID:    req.UserID,
		Status:    model.StatusPending,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.repo.CreateBooking(ctx, booking); err != nil {
		return nil, fmt.Errorf("create booking: %w", err)
	}

	return &model.BookSeatResponse{
		BookingID: booking.ID,
		Status:    booking.Status,
		ExpiresAt: booking.ExpiresAt,
	}, nil
}

func (s *Service) ConfirmBooking(ctx context.Context, bookingID uuid.UUID) error {
	booking, err := s.repo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return err
	}

	if booking.Status == model.StatusCancelled {
		return errors.New("cannot confirm cancelled booking")
	}

	if time.Now().After(booking.ExpiresAt) {
		_ = s.repo.UpdateBookingStatus(ctx, bookingID, model.StatusCancelled)
		return model.ErrBookingExpired
	}

	return s.repo.UpdateBookingStatus(ctx, bookingID, model.StatusConfirmed)
}

func (s *Service) CancelBooking(ctx context.Context, bookingID uuid.UUID) error {
	booking, err := s.repo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return err
	}

	if booking.Status == model.StatusCancelled {
		return nil
	}

	return s.repo.UpdateBookingStatus(ctx, bookingID, model.StatusCancelled)
}

func (s *Service) GetUserBookings(ctx context.Context, userID uuid.UUID) ([]*model.Booking, error) {

	return []*model.Booking{}, nil
}

func (s *Service) ProcessExpiredBookings(ctx context.Context) (int64, error) {
	count, err := s.repo.CancelExpiredBookings(ctx)
	if err != nil {
		return 0, fmt.Errorf("process expired bookings: %w", err)
	}

	return count, nil
}
