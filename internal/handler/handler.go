package handler

import (
	"EventBooker/internal/model"
	"context"

	"github.com/google/uuid"
)

type BookerService interface {
	// RegisterUser регистрирует нового пользователя с базовой валидацией и проверкой уникальности email.
	RegisterUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error)
	// GetUserProfile возвращает профиль пользователя по его ID.
	GetUserProfile(ctx context.Context, userID uuid.UUID) (*model.User, error)

	// CreateEvent создает новое мероприятие и рассчитывает количество доступных мест.
	CreateEvent(ctx context.Context, req *model.CreateEventRequest) (*model.Event, error)
	// GetEventDetails возвращает подробную информацию о событии и списке его броней.
	GetEventDetails(ctx context.Context, eventID uuid.UUID) (*model.EventDetailResponse, error)
	// ListUpcomingEvents возвращает список всех актуальных мероприятий для главных страниц/каталога.
	ListUpcomingEvents(ctx context.Context) ([]*model.Event, error)

	// BookSeat резервирует место на мероприятие с проверкой свободных мест и таймаута брони.
	// Метод выполняется внутри транзакции.
	BookSeat(ctx context.Context, eventID uuid.UUID, req *model.BookSeatRequest) (*model.BookSeatResponse, error)
	// ConfirmBooking переводит статус бронирования в 'confirmed' (например, после оплаты).
	ConfirmBooking(ctx context.Context, bookingID uuid.UUID) error
	// CancelBooking досрочно отменяет бронь и возвращает место обратно в доступные события.
	CancelBooking(ctx context.Context, bookingID uuid.UUID) error
	// GetUserBookings возвращает все бронирования конкретного пользователя.
	GetUserBookings(ctx context.Context, userID uuid.UUID) ([]*model.Booking, error)

	// ProcessExpiredBookings отменяет все просроченные брони (у которых истек TTL)
	// и возвращает зарезервированные места в события. Вызывается по таймеру (Cron).
	ProcessExpiredBookings(ctx context.Context) (int64, error)
}
