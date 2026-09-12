package handlers

import (
	"backend-api-go/internal/apperr"
	"backend-api-go/internal/models"
	"backend-api-go/internal/services"
	"backend-api-go/pkg/middleware"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PurchaseHandler struct {
	purchaseService services.PurchaseService
}

// NewPurchaseHandler constructs a PurchaseHandler backed by the given PurchaseService.
func NewPurchaseHandler(purchaseService services.PurchaseService) *PurchaseHandler {
	return &PurchaseHandler{purchaseService: purchaseService}
}

// CreatePurchase handles POST /purchases: it validates the request and creates a
// purchase for the authenticated user, decrementing product stock.
func (h *PurchaseHandler) CreatePurchase(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req models.CreatePurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	purchase, err := h.purchaseService.CreatePurchase(c.Request.Context(), userID, &req)
	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, apperr.ErrProductNotFound):
			status = http.StatusNotFound
		case errors.Is(err, apperr.ErrInsufficientStock):
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, purchase)
}

// GetPurchaseHistory handles GET /purchases: it returns a paginated list of the
// authenticated user's purchases.
func (h *PurchaseHandler) GetPurchaseHistory(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	resp, err := h.purchaseService.GetPurchaseHistory(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetPurchase handles GET /purchases/:id: it returns a single purchase owned by
// the authenticated user.
func (h *PurchaseHandler) GetPurchase(c *gin.Context) {
	userID := middleware.GetUserID(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid purchase id"})
		return
	}

	purchase, err := h.purchaseService.GetPurchase(c.Request.Context(), id, userID)
	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, apperr.ErrPurchaseNotFound):
			status = http.StatusNotFound
		case errors.Is(err, apperr.ErrPurchaseForbidden):
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, purchase)
}

// CancelPurchase handles PATCH /purchases/:id/cancel: it cancels a purchase owned
// by the authenticated user and restores the product's stock.
func (h *PurchaseHandler) CancelPurchase(c *gin.Context) {
	userID := middleware.GetUserID(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid purchase id"})
		return
	}

	if err := h.purchaseService.CancelPurchase(c.Request.Context(), id, userID); err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, apperr.ErrPurchaseNotFound):
			status = http.StatusNotFound
		case errors.Is(err, apperr.ErrPurchaseForbidden):
			status = http.StatusForbidden
		case errors.Is(err, apperr.ErrPurchaseNotCancellable):
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "purchase cancelled successfully"})
}
