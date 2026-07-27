package service

import (
	"EventBooker/internal/model"
	"context"
)

type Store interface {
	// Users
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id int64) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)

	// Events
	CreateEvent(ctx context.Context, event *model.Event) error
	GetEventByID(ctx context.Context, id int64) (*model.Event, error)
	ListEvents(ctx context.Context) ([]*model.Event, error)

	// Bookings
	CreateBookingWithTx(ctx context.Context, booking *model.Booking) error
	GetBookingByID(ctx context.Context, id int64) (*model.Booking, error)
	UpdateBookingStatus(ctx context.Context, bookingID int64, status model.BookingStatus) error
	ListBookingsByEventID(ctx context.Context, eventID int64) ([]*model.Booking, error)
	CancelExpiredBookings(ctx context.Context) (int64, error)
}
