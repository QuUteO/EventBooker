package handler

import (
	"context"
	"errors"
	"net/http"

	"EventBooker/internal/model"

	"github.com/wb-go/wbf/ginext" // Укажи актуальный путь к своему ginext

	"github.com/google/uuid"
)

type BookerService interface {
	RegisterUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error)
	GetUserProfile(ctx context.Context, userID uuid.UUID) (*model.User, error)

	CreateEvent(ctx context.Context, req *model.CreateEventRequest) (*model.Event, error)
	GetEventDetails(ctx context.Context, eventID uuid.UUID) (*model.EventDetailResponse, error)
	ListUpcomingEvents(ctx context.Context) ([]*model.Event, error)

	BookSeat(ctx context.Context, eventID uuid.UUID, req *model.BookSeatRequest) (*model.BookSeatResponse, error)
	ConfirmBooking(ctx context.Context, bookingID uuid.UUID) error
	CancelBooking(ctx context.Context, bookingID uuid.UUID) error
	GetUserBookings(ctx context.Context, userID uuid.UUID) ([]*model.Booking, error)

	ProcessExpiredBookings(ctx context.Context) (int64, error)
}

type Handler struct {
	service BookerService
}

func New(service BookerService) *Handler {
	return &Handler{service: service}
}

// InitRoutes настраивает роутинг с использованием ginext.Engine
func (h *Handler) InitRoutes(ginMode string) *ginext.Engine {
	engine := ginext.New(ginMode)

	// Подключаем стандартные мидлвары через ginext
	engine.Use(ginext.Logger(), ginext.Recovery())

	engine.StaticFile("/", "./web/index.html")

	v1 := engine.Group("/api/v1")
	{
		// Users (сохраняем для управления профилем и списком броней пользователя)
		users := v1.Group("/users")
		{
			users.POST("", h.RegisterUser)
			users.GET("/:id", h.GetUserProfile)
			users.GET("/:id/bookings", h.GetUserBookings)
		}

		// Events (согласно заданному контракту)
		events := v1.Group("/events")
		{
			events.POST("", h.CreateEvent)                // POST /events
			events.GET("", h.ListUpcomingEvents)          // GET /events (для списка на главной)
			events.GET("/:id", h.GetEventDetails)         // GET /events/{id}
			events.POST("/:id/book", h.BookSeat)          // POST /events/{id}/book
			events.POST("/:id/confirm", h.ConfirmBooking) // POST /events/{id}/confirm
		}

		// Bookings (для отмены конкретной брони)
		bookings := v1.Group("/bookings")
		{
			bookings.POST("/:id/cancel", h.CancelBooking)
		}
	}

	return engine
}

// renderError централизованно обрабатывает ошибки доменного слоя
func renderError(c *ginext.Context, err error) {
	status := http.StatusInternalServerError

	switch {
	case errors.Is(err, model.ErrUserNotFound), errors.Is(err, model.ErrEventNotFound), errors.Is(err, model.ErrBookingNotFound):
		status = http.StatusNotFound
	case errors.Is(err, model.ErrEmailAlreadyTaken):
		status = http.StatusConflict
	case errors.Is(err, model.ErrNoAvailableSeats), errors.Is(err, model.ErrBookingExpired), errors.Is(err, model.ErrEventExpired), errors.Is(err, model.ErrTotalSeats):
		status = http.StatusBadRequest
	}

	c.JSON(status, ginext.H{"error": err.Error()})
}
