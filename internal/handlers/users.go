package handlers

import (
	"backend-api-go/internal/models"
	"backend-api-go/internal/services"
	"backend-api-go/pkg/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetProfile godoc
// GET /users/me
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	user, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateProfile godoc
// PATCH /users/me
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteAccount godoc
// DELETE /users/me
func (h *UserHandler) DeleteAccount(c *gin.Context) {
	userID := middleware.GetUserID(c)

	if err := h.userService.DeleteAccount(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "account deleted successfully"})
}
