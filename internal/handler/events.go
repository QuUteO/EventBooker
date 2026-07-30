package handler

import (
	"net/http"

	"EventBooker/internal/model"

	"github.com/wb-go/wbf/ginext"

	"github.com/google/uuid"
)

// POST /api/v1/events
func (h *Handler) CreateEvent(c *ginext.Context) {
	var req model.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid JSON payload"})
		return
	}

	event, err := h.service.CreateEvent(c.Request.Context(), &req)
	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusCreated, event)
}

// GET /api/v1/events
func (h *Handler) ListUpcomingEvents(c *ginext.Context) {
	events, err := h.service.ListUpcomingEvents(c.Request.Context())
	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, events)
}

// GET /api/v1/events/:id
func (h *Handler) GetEventDetails(c *ginext.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid event UUID"})
		return
	}

	details, err := h.service.GetEventDetails(c.Request.Context(), eventID)
	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, details)
}
