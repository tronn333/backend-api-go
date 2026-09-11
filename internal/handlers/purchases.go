package handlers

import (
	"backend-api-go/internal/models"
	"backend-api-go/internal/services"
	"backend-api-go/pkg/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PurchaseHandler struct {
	purchaseService services.PurchaseService
}

func NewPurchaseHandler(purchaseService services.PurchaseService) *PurchaseHandler {
	return &PurchaseHandler{purchaseService: purchaseService}
}

// CreatePurchase godoc
// POST /purchases
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
		msg := err.Error()
		if msg == "product not found" {
			status = http.StatusNotFound
		} else if msg == "insufficient stock" || containsPrefix(msg, "insufficient stock:") {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusCreated, purchase)
}

// GetPurchaseHistory godoc
// GET /purchases
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

// GetPurchase godoc
// GET /purchases/:id
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
		if err.Error() == "purchase not found" {
			status = http.StatusNotFound
		} else if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, purchase)
}

// CancelPurchase godoc
// PATCH /purchases/:id/cancel
func (h *PurchaseHandler) CancelPurchase(c *gin.Context) {
	userID := middleware.GetUserID(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid purchase id"})
		return
	}

	if err := h.purchaseService.CancelPurchase(c.Request.Context(), id, userID); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "purchase not found" {
			status = http.StatusNotFound
		} else if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "purchase cancelled successfully"})
}

func containsPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
