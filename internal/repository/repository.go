package repository

import (
	"EventBooker/internal/model"
	"context"
	"fmt"

	"github.com/google/uuid"
	pgxdriver "github.com/wb-go/wbf/dbpg/pgx-driver"
)

type EventRepository struct {
	conn *pgxdriver.Postgres
}

func New(conn *pgxdriver.Postgres) *EventRepository {
	return &EventRepository{conn: conn}
}

func (r *EventRepository) CreateUser(ctx context.Context, user *model.User) error {
	query := "INSERT INTO users (id, email, telegram_use, created_at) VALUES ($1, $2, $3, $4)"

	_, err := r.conn.Exec(ctx, query,
		user.ID,
		user.Email,
		user.TelegramUse,
		user.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("репоизторий: ошибка создания пользователя: %w", err)
	}

	return nil
}

func (r *EventRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User

	query := "SELECT id, email, telegram_use, created_at FROM users WHERE id = $1"

	if err := r.conn.QueryRow(ctx, query, id).Scan(&user.ID, &user.Email, &user.TelegramUse, &user.CreatedAt); err != nil {
		return nil, fmt.Errorf("репозиторий: ошибка получения пользователя: %w", err)
	}

	return &user, nil
}

func (r *EventRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	query := "SELECT id, email, telegram_use, created_at FROM users WHERE email = $1"

	if err := r.conn.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.TelegramUse, &user.CreatedAt); err != nil {
		return nil, fmt.Errorf("репозиторий: ошибка получения пользователя: %w", err)
	}

	return &user, nil
}

func (r *EventRepository) CreateEvent(ctx context.Context, event *model.Event) error {
	query := "INSERT INTO events (id, title, event_date, total_seats, available_seats, booking_ttl_minutes, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	_, err := r.conn.Exec(ctx, query,
		event.ID,
		event.Title,
		event.EventDate,
		event.TotalSeats,
		event.AvailableSeats,
		event.BookingTTLMinutes,
		event.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("репозиторий: ошибка сохранения ивента: %w", err)
	}

	return nil
}

func (r *EventRepository) GetEventByID(ctx context.Context, id uuid.UUID) (*model.Event, error) {
	var event model.Event

	query := "SELECT id, title, event_date, total_seats, available_seats, booking_ttl_minutes, created_at FROM events WHERE id = $1"

	if err := r.conn.QueryRow(ctx, query, id.String()).Scan(
		&event.ID,
		&event.Title,
		&event.EventDate,
		&event.TotalSeats,
		&event.AvailableSeats,
		&event.BookingTTLMinutes,
		&event.CreatedAt,
	); err != nil {
		return nil, fmt.Errorf("репозиторий: ошибка получения ивента: %w", err)
	}

	return &event, nil
}

func (r *EventRepository) ListEvents(ctx context.Context) ([]*model.Event, error) {
	events := make([]*model.Event, 0, 10)

	query := "SELECT id, title, event_date, total_seats, available_seats, booking_ttl_minutes, created_at FROM events LIMIT 10"

	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("репозиторий: ошибка получения ивентов: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var event model.Event

		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.EventDate,
			&event.TotalSeats,
			&event.AvailableSeats,
			&event.BookingTTLMinutes,
			&event.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("репозиторий: ошибка сохранения ивента: %w", err)
		}

		events = append(events, &event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (r *EventRepository) CreateBooking(ctx context.Context, booking *model.Booking) error {
	query := "INSERT INTO bookings (id, event_id, user_id, status, expires_at, created_at) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err := r.conn.Exec(ctx, query,
		&booking.ID,
		&booking.EventID,
		&booking.UserID,
		&booking.Status,
		&booking.ExpiresAt,
		&booking.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("репозиторий: ошибка сохранения брони: %w", err)
	}

	return nil
}

func (r *EventRepository) GetBookingByID(ctx context.Context, id uuid.UUID) (*model.Booking, error) {
	var booking model.Booking

	query := "SELECT id, event_id, user_id, status, expires_at, created_at FROM bookings WHERE id = $1"

	err := r.conn.QueryRow(ctx, query, id).Scan(
		&booking.ID,
		&booking.EventID,
		&booking.UserID,
		&booking.Status,
		&booking.ExpiresAt,
		&booking.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("репозиторий: ошибка получения брони: %w", err)
	}

	return &booking, nil
}

func (r *EventRepository) UpdateBookingStatus(ctx context.Context, bookingID uuid.UUID, status model.BookingStatus) error {
	query := "UPDATE bookings SET status = $1 WHERE id = $2"

	_, err := r.conn.Exec(ctx, query, status, bookingID)
	if err != nil {
		return fmt.Errorf("репозитрий: ошибка обновления статус брони: %w", err)
	}

	return nil
}

func (r *EventRepository) ListBookingsByEventID(ctx context.Context, eventID uuid.UUID) ([]*model.Booking, error) {
	query := "SELECT id, event_id, user_id, status, expires_at, created_at FROM bookings WHERE event_id = $1 LIMIT 10"

	bookings := make([]*model.Booking, 0, 10)

	rows, err := r.conn.Query(ctx, query, eventID)
	if err != nil {
		return nil, fmt.Errorf("репозиторий: ошибка получения списка брони: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var booking model.Booking
		err := rows.Scan(
			&booking.ID,
			&booking.EventID,
			&booking.UserID,
			&booking.Status,
			&booking.ExpiresAt,
			&booking.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("репозиторий: ошибка скана брони: %w", err)
		}

		bookings = append(bookings, &booking)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}
