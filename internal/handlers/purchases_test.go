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

func newPurchaseCtx(t *testing.T, method, path, body string, userID int) (*httptest.ResponseRecorder, *gin.Context) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		c.Request.Header.Set("Content-Type", "application/json")
	}
	c.Set(middleware.UserIDKey, userID)
	return w, c
}

func TestCreatePurchase_Success(t *testing.T) {
	mock := &mockPurchaseService{}
	mock.createPurchaseFn = func(ctx context.Context, userID int, req *models.CreatePurchaseRequest) (*models.Purchase, error) {
		return &models.Purchase{ID: 1, UserID: userID, ProductID: req.ProductID, Quantity: req.Quantity}, nil
	}
	h := NewPurchaseHandler(mock)

	w, c := newPurchaseCtx(t, http.MethodPost, "/purchases", `{"product_id":1,"quantity":2}`, 3)
	h.CreatePurchase(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}

func TestCreatePurchase_ProductNotFound(t *testing.T) {
	mock := &mockPurchaseService{}
	mock.createPurchaseFn = func(ctx context.Context, userID int, req *models.CreatePurchaseRequest) (*models.Purchase, error) {
		return nil, apperr.ErrProductNotFound
	}
	h := NewPurchaseHandler(mock)

	w, c := newPurchaseCtx(t, http.MethodPost, "/purchases", `{"product_id":999,"quantity":2}`, 3)
	h.CreatePurchase(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestCreatePurchase_InsufficientStock(t *testing.T) {
	mock := &mockPurchaseService{}
	mock.createPurchaseFn = func(ctx context.Context, userID int, req *models.CreatePurchaseRequest) (*models.Purchase, error) {
		return nil, apperr.ErrInsufficientStock
	}
	h := NewPurchaseHandler(mock)

	w, c := newPurchaseCtx(t, http.MethodPost, "/purchases", `{"product_id":1,"quantity":999}`, 3)
	h.CreatePurchase(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreatePurchase_BadRequest(t *testing.T) {
	h := NewPurchaseHandler(&mockPurchaseService{})

	w, c := newPurchaseCtx(t, http.MethodPost, "/purchases", `{"product_id":1}`, 3)
	h.CreatePurchase(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestGetPurchase_Success(t *testing.T) {
	mock := &mockPurchaseService{}
	mock.getPurchaseFn = func(ctx context.Context, id, userID int) (*models.Purchase, error) {
		return &models.Purchase{ID: id, UserID: userID}, nil
	}
	h := NewPurchaseHandler(mock)

	w, c := newPurchaseCtx(t, http.MethodGet, "/purchases/5", "", 3)
	c.Params = gin.Params{{Key: "id", Value: "5"}}
	h.GetPurchase(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetPurchase_Forbidden(t *testing.T) {
	mock := &mockPurchaseService{}
	mock.getPurchaseFn = func(ctx context.Context, id, userID int) (*models.Purchase, error) {
		return nil, apperr.ErrPurchaseForbidden
	}
	h := NewPurchaseHandler(mock)

	w, c := newPurchaseCtx(t, http.MethodGet, "/purchases/5", "", 3)
	c.Params = gin.Params{{Key: "id", Value: "5"}}
	h.GetPurchase(c)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestCancelPurchase_Success(t *testing.T) {
	mock := &mockPurchaseService{}
	mock.cancelPurchaseFn = func(ctx context.Context, id, userID int) error { return nil }
	h := NewPurchaseHandler(mock)

	w, c := newPurchaseCtx(t, http.MethodPatch, "/purchases/5/cancel", "", 3)
	c.Params = gin.Params{{Key: "id", Value: "5"}}
	h.CancelPurchase(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestCancelPurchase_Forbidden(t *testing.T) {
	mock := &mockPurchaseService{}
	mock.cancelPurchaseFn = func(ctx context.Context, id, userID int) error {
		return apperr.ErrPurchaseForbidden
	}
	h := NewPurchaseHandler(mock)

	w, c := newPurchaseCtx(t, http.MethodPatch, "/purchases/5/cancel", "", 3)
	c.Params = gin.Params{{Key: "id", Value: "5"}}
	h.CancelPurchase(c)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestCancelPurchase_NotCancellable(t *testing.T) {
	mock := &mockPurchaseService{}
	mock.cancelPurchaseFn = func(ctx context.Context, id, userID int) error {
		return apperr.ErrPurchaseNotCancellable
	}
	h := NewPurchaseHandler(mock)

	w, c := newPurchaseCtx(t, http.MethodPatch, "/purchases/5/cancel", "", 3)
	c.Params = gin.Params{{Key: "id", Value: "5"}}
	h.CancelPurchase(c)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", w.Code)
	}
}

func TestCancelPurchase_NotFound(t *testing.T) {
	mock := &mockPurchaseService{}
	mock.cancelPurchaseFn = func(ctx context.Context, id, userID int) error {
		return apperr.ErrPurchaseNotFound
	}
	h := NewPurchaseHandler(mock)

	w, c := newPurchaseCtx(t, http.MethodPatch, "/purchases/5/cancel", "", 3)
	c.Params = gin.Params{{Key: "id", Value: "5"}}
	h.CancelPurchase(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
