package handler

import (
	"net/http"

	"EventBooker/internal/model"

	"github.com/wb-go/wbf/ginext"

	"github.com/google/uuid"
)

// POST /api/v1/events/:id/book
func (h *Handler) BookSeat(c *ginext.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid event UUID"})
		return
	}

	var req model.BookSeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid JSON payload"})
		return
	}

	resp, err := h.service.BookSeat(c.Request.Context(), eventID, &req)
	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// POST /api/v1/bookings/:id/confirm
func (h *Handler) ConfirmBooking(c *ginext.Context) {
	bookingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid booking UUID"})
		return
	}

	if err := h.service.ConfirmBooking(c.Request.Context(), bookingID); err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, ginext.H{"status": "confirmed"})
}

// POST /api/v1/bookings/:id/cancel
func (h *Handler) CancelBooking(c *ginext.Context) {
	bookingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid booking UUID"})
		return
	}

	if err := h.service.CancelBooking(c.Request.Context(), bookingID); err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, ginext.H{"status": "cancelled"})
}

// GET /api/v1/users/:id/bookings
func (h *Handler) GetUserBookings(c *ginext.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid user UUID"})
		return
	}

	bookings, err := h.service.GetUserBookings(c.Request.Context(), userID)
	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, bookings)
}
