package handler

import (
	"net/http"

	"EventBooker/internal/model"

	"github.com/wb-go/wbf/ginext"

	"github.com/google/uuid"
)

// POST /api/v1/users
func (h *Handler) RegisterUser(c *ginext.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid JSON payload"})
		return
	}

	user, err := h.service.RegisterUser(c.Request.Context(), &req)
	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GET /api/v1/users/:id
func (h *Handler) GetUserProfile(c *ginext.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid user UUID"})
		return
	}

	user, err := h.service.GetUserProfile(c.Request.Context(), userID)
	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}
