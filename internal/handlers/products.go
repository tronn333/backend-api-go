package handlers

import (
	"backend-api-go/internal/models"
	"backend-api-go/internal/services"
	"backend-api-go/pkg/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService services.ProductService
}

func NewProductHandler(productService services.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

// ListProducts godoc
// GET /products
func (h *ProductHandler) ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	category := c.Query("category")

	resp, err := h.productService.ListProducts(c.Request.Context(), page, pageSize, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetProduct godoc
// GET /products/:id
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	product, err := h.productService.GetProduct(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// CreateProduct godoc
// POST /products
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req models.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.productService.CreateProduct(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// UpdateProduct godoc
// PATCH /products/:id
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	userID := middleware.GetUserID(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	var req models.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.productService.UpdateProduct(c.Request.Context(), id, userID, &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "product not found" {
			status = http.StatusNotFound
		} else if err.Error() == "forbidden: you do not own this product" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct godoc
// DELETE /products/:id
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	userID := middleware.GetUserID(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	if err := h.productService.DeleteProduct(c.Request.Context(), id, userID); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "product not found" {
			status = http.StatusNotFound
		} else if err.Error() == "forbidden: you do not own this product" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
}
