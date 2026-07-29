package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrEventNotFound     = errors.New("event not found")
	ErrBookingNotFound   = errors.New("booking not found")
	ErrNoAvailableSeats  = errors.New("на это мероприятие нет свободных мест")
	ErrBookingExpired    = errors.New("срок бронирования истек")
	ErrEventExpired      = errors.New("срок ивента истек")
	ErrTotalSeats        = errors.New("недостаточно мест")
	ErrEmailAlreadyTaken = errors.New("email уже зарегистрирован")
)

// User — отражает таблицу users
type User struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	TelegramUse string    `json:"telegram_user"`
	CreatedAt   time.Time `json:"created_at"`
}

// Event — отражает таблицу events
type Event struct {
	ID                uuid.UUID `json:"id"`
	Title             string    `json:"title"`
	EventDate         time.Time `json:"event_date"`
	TotalSeats        int       `json:"total_seats"`
	AvailableSeats    int       `json:"available_seats"`
	BookingTTLMinutes int       `json:"booking_ttl_minutes"`
	CreatedAt         time.Time `json:"created_at"`
}

// BookingStatus — пользовательский тип для статусов
type BookingStatus string

const (
	StatusPending   BookingStatus = "pending"
	StatusConfirmed BookingStatus = "confirmed"
	StatusCancelled BookingStatus = "cancelled"
)

// Booking — отражает таблицу bookings
type Booking struct {
	ID        uuid.UUID     `json:"id"`
	EventID   uuid.UUID     `json:"event_id"`
	UserID    uuid.UUID     `json:"user_id"`
	Status    BookingStatus `json:"status"`
	ExpiresAt time.Time     `json:"expires_at"`
	CreatedAt time.Time     `json:"created_at"`
}

type CreateUserRequest struct {
	Email       string `json:"email"`
	TelegramUse string `json:"telegram_use"`
}

type CreateEventRequest struct {
	Title             string    `json:"title"`
	EventDate         time.Time `json:"event_date"`
	TotalSeats        int       `json:"total_seats"`
	BookingTTLMinutes int       `json:"booking_ttl_minutes,omitempty"`
}

type BookSeatRequest struct {
	UserID uuid.UUID `json:"user_id"`
}

type BookSeatResponse struct {
	BookingID uuid.UUID     `json:"booking_id"`
	Status    BookingStatus `json:"status"`
	ExpiresAt time.Time     `json:"expires_at"`
}

type EventDetailResponse struct {
	Event    Event     `json:"event"`
	Bookings []Booking `json:"bookings,omitempty"`
}
