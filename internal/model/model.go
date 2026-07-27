package model

import "time"

// Event — отражает таблицу events
type Event struct {
	ID                int64     `json:"id"`
	Title             string    `json:"title"`
	EventDate         time.Time `json:"event_date"`
	TotalSeats        int       `json:"total_seats"`
	AvailableSeats    int       `json:"available_seats"`
	BookingTTLMinutes int       `json:"booking_ttl_minutes"`
	CreatedAt         time.Time `json:"created_at"`
}

// BookingStatus — определяем пользовательский тип для статусов, чтобы избежать опечаток
type BookingStatus string

const (
	StatusPending   BookingStatus = "pending"
	StatusConfirmed BookingStatus = "confirmed"
	StatusCancelled BookingStatus = "cancelled"
)

// Booking — отражает таблицу bookings
type Booking struct {
	ID        int64         `json:"id"`
	EventID   int64         `json:"event_id"`
	UserID    int64         `json:"user_id"`
	Status    BookingStatus `json:"status"`
	ExpiresAt time.Time     `json:"expires_at"`
	CreatedAt time.Time     `json:"created_at"`
}

type CreateEventRequest struct {
	Title             string    `json:"title"`
	EventDate         time.Time `json:"event_date"`
	TotalSeats        int       `json:"total_seats"`
	BookingTTLMinutes int       `json:"booking_ttl_minutes,omitempty"` // если не передали — поставим дефолт
}

type BookSeatRequest struct {
	UserID int64 `json:"user_id"`
}

type BookSeatResponse struct {
	BookingID int64     `json:"booking_id"`
	Status    string    `json:"status"`
	ExpiresAt time.Time `json:"expires_at"`
}

type EventDetailResponse struct {
	Event    Event     `json:"event"`
	Bookings []Booking `json:"bookings,omitempty"` // список броней
}
