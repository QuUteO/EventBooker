package service

import (
	"EventBooker/internal/model"
	"context"

	"github.com/google/uuid"
)

type Store interface {
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
