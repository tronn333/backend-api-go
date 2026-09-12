package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"backend-api-go/internal/apperr"
	"backend-api-go/internal/models"
	"backend-api-go/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func newProductCtx(t *testing.T, method, path, body string, userID int, role string) (*httptest.ResponseRecorder, *gin.Context) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		c.Request.Header.Set("Content-Type", "application/json")
	}
	c.Set(middleware.UserIDKey, userID)
	c.Set(middleware.RoleKey, role)
	return w, c
}

func TestListProducts_Success(t *testing.T) {
	mock := &mockProductService{}
	mock.listProductsFn = func(ctx context.Context, page, pageSize int, category string) (*models.ProductListResponse, error) {
		return &models.ProductListResponse{Products: []models.Product{}, Total: 0, Page: page, PageSize: pageSize}, nil
	}
	h := NewProductHandler(mock)

	w, c := newProductCtx(t, http.MethodGet, "/products?page=1&page_size=20", "", 0, "")
	h.ListProducts(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetProduct_Success(t *testing.T) {
	mock := &mockProductService{}
	mock.getProductFn = func(ctx context.Context, id int) (*models.Product, error) {
		return &models.Product{ID: id, Name: "Widget"}, nil
	}
	h := NewProductHandler(mock)

	w, c := newProductCtx(t, http.MethodGet, "/products/5", "", 0, "")
	c.Params = gin.Params{{Key: "id", Value: "5"}}
	h.GetProduct(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetProduct_NotFound(t *testing.T) {
	mock := &mockProductService{}
	mock.getProductFn = func(ctx context.Context, id int) (*models.Product, error) {
		return nil, apperr.ErrProductNotFound
	}
	h := NewProductHandler(mock)

	w, c := newProductCtx(t, http.MethodGet, "/products/5", "", 0, "")
	c.Params = gin.Params{{Key: "id", Value: "5"}}
	h.GetProduct(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestGetProduct_InvalidID(t *testing.T) {
	h := NewProductHandler(&mockProductService{})

	w, c := newProductCtx(t, http.MethodGet, "/products/abc", "", 0, "")
	c.Params = gin.Params{{Key: "id", Value: "abc"}}
	h.GetProduct(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreateProduct_Success(t *testing.T) {
	mock := &mockProductService{}
	mock.createProductFn = func(ctx context.Context, sellerID int, req *models.CreateProductRequest) (*models.Product, error) {
		return &models.Product{ID: 1, SellerID: sellerID, Name: req.Name}, nil
	}
	h := NewProductHandler(mock)

	w, c := newProductCtx(t, http.MethodPost, "/products", `{"name":"Widget","price":10,"stock":3}`, 7, "")
	h.CreateProduct(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}

func TestCreateProduct_BadRequest(t *testing.T) {
	h := NewProductHandler(&mockProductService{})

	w, c := newProductCtx(t, http.MethodPost, "/products", `{"name":"Widget"}`, 7, "")
	h.CreateProduct(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUpdateProduct_Success(t *testing.T) {
	mock := &mockProductService{}
	mock.updateProductFn = func(ctx context.Context, id, sellerID int, isAdmin bool, req *models.UpdateProductRequest) (*models.Product, error) {
		return &models.Product{ID: id, SellerID: sellerID, Name: req.Name}, nil
	}
	h := NewProductHandler(mock)

	w, c := newProductCtx(t, http.MethodPatch, "/products/5", `{"name":"New"}`, 7, "")
	c.Params = gin.Params{{Key: "id", Value: "5"}}
	h.UpdateProduct(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestUpdateProduct_Forbidden(t *testing.T) {
	mock := &mockProductService{}
	mock.updateProductFn = func(ctx context.Context, id, sellerID int, isAdmin bool, req *models.UpdateProductRequest) (*models.Product, error) {
		return nil, apperr.ErrProductForbidden
	}
	h := NewProductHandler(mock)

	w, c := newProductCtx(t, http.MethodPatch, "/products/5", `{"name":"New"}`, 7, "")
	c.Params = gin.Params{{Key: "id", Value: "5"}}
	h.UpdateProduct(c)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestUpdateProduct_PassesAdminFlag(t *testing.T) {
	mock := &mockProductService{}
	var gotAdmin bool
	mock.updateProductFn = func(ctx context.Context, id, sellerID int, isAdmin bool, req *models.UpdateProductRequest) (*models.Product, error) {
		gotAdmin = isAdmin
		return &models.Product{ID: id}, nil
	}
	h := NewProductHandler(mock)

	w, c := newProductCtx(t, http.MethodPatch, "/products/5", `{"name":"New"}`, 7, "admin")
	c.Params = gin.Params{{Key: "id", Value: "5"}}
	h.UpdateProduct(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !gotAdmin {
		t.Fatal("expected isAdmin=true to be passed to the service")
	}
}

func TestDeleteProduct_Success(t *testing.T) {
	mock := &mockProductService{}
	mock.deleteProductFn = func(ctx context.Context, id, sellerID int, isAdmin bool) error { return nil }
	h := NewProductHandler(mock)

	w, c := newProductCtx(t, http.MethodDelete, "/products/5", "", 7, "")
	c.Params = gin.Params{{Key: "id", Value: "5"}}
	h.DeleteProduct(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestDeleteProduct_Forbidden(t *testing.T) {
	mock := &mockProductService{}
	mock.deleteProductFn = func(ctx context.Context, id, sellerID int, isAdmin bool) error {
		return apperr.ErrProductForbidden
	}
	h := NewProductHandler(mock)

	w, c := newProductCtx(t, http.MethodDelete, "/products/5", "", 7, "")
	c.Params = gin.Params{{Key: "id", Value: "5"}}
	h.DeleteProduct(c)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}
