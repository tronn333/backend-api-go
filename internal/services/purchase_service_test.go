package services

import (
	"context"
	"errors"
	"testing"

	"backend-api-go/internal/apperr"
	"backend-api-go/internal/models"
)

func TestCreatePurchase_Success(t *testing.T) {
	prodRepo := newFakeProductRepo()
	prodRepo.products[1] = &models.Product{ID: 1, Price: 10, Stock: 5}
	purRepo := newFakePurchaseRepo()
	svc := NewPurchaseService(purRepo, prodRepo)

	pu, err := svc.CreatePurchase(context.Background(), 3, &models.CreatePurchaseRequest{ProductID: 1, Quantity: 2})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if pu.TotalPrice != 20 {
		t.Fatalf("expected total price 20, got %v", pu.TotalPrice)
	}
	if pu.Status != models.StatusCompleted {
		t.Fatalf("expected status completed, got %v", pu.Status)
	}
	if prodRepo.products[1].Stock != 3 {
		t.Fatalf("expected stock decremented to 3, got %d", prodRepo.products[1].Stock)
	}
}

func TestCreatePurchase_ProductNotFound(t *testing.T) {
	svc := NewPurchaseService(newFakePurchaseRepo(), newFakeProductRepo())

	_, err := svc.CreatePurchase(context.Background(), 3, &models.CreatePurchaseRequest{ProductID: 999, Quantity: 1})
	if !errors.Is(err, apperr.ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
}

func TestCreatePurchase_InsufficientStock(t *testing.T) {
	prodRepo := newFakeProductRepo()
	prodRepo.products[1] = &models.Product{ID: 1, Price: 10, Stock: 2}
	svc := NewPurchaseService(newFakePurchaseRepo(), prodRepo)

	_, err := svc.CreatePurchase(context.Background(), 3, &models.CreatePurchaseRequest{ProductID: 1, Quantity: 5})
	if !errors.Is(err, apperr.ErrInsufficientStock) {
		t.Fatalf("expected ErrInsufficientStock, got %v", err)
	}
	if prodRepo.products[1].Stock != 2 {
		t.Fatal("stock must not change when insufficient")
	}
}

func TestCreatePurchase_RollsBackStockOnCreateFailure(t *testing.T) {
	prodRepo := newFakeProductRepo()
	prodRepo.products[1] = &models.Product{ID: 1, Price: 10, Stock: 5}
	purRepo := newFakePurchaseRepo()
	purRepo.createErr = errors.New("db down")
	svc := NewPurchaseService(purRepo, prodRepo)

	_, err := svc.CreatePurchase(context.Background(), 3, &models.CreatePurchaseRequest{ProductID: 1, Quantity: 2})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if prodRepo.products[1].Stock != 5 {
		t.Fatalf("expected stock rolled back to 5, got %d", prodRepo.products[1].Stock)
	}
}

func TestGetPurchase_Owner(t *testing.T) {
	purRepo := newFakePurchaseRepo()
	purRepo.purchases[1] = &models.Purchase{ID: 1, UserID: 3, ProductID: 1}
	svc := NewPurchaseService(purRepo, newFakeProductRepo())

	pu, err := svc.GetPurchase(context.Background(), 1, 3)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if pu.UserID != 3 {
		t.Fatalf("unexpected purchase: %+v", pu)
	}
}

func TestGetPurchase_Forbidden(t *testing.T) {
	purRepo := newFakePurchaseRepo()
	purRepo.purchases[1] = &models.Purchase{ID: 1, UserID: 3}
	svc := NewPurchaseService(purRepo, newFakeProductRepo())

	_, err := svc.GetPurchase(context.Background(), 1, 99)
	if !errors.Is(err, apperr.ErrPurchaseForbidden) {
		t.Fatalf("expected ErrPurchaseForbidden, got %v", err)
	}
}

func TestGetPurchase_NotFound(t *testing.T) {
	svc := NewPurchaseService(newFakePurchaseRepo(), newFakeProductRepo())

	_, err := svc.GetPurchase(context.Background(), 999, 3)
	if !errors.Is(err, apperr.ErrPurchaseNotFound) {
		t.Fatalf("expected ErrPurchaseNotFound, got %v", err)
	}
}

func TestCancelPurchase_Success_RestoresStock(t *testing.T) {
	prodRepo := newFakeProductRepo()
	prodRepo.products[1] = &models.Product{ID: 1, Stock: 10}
	purRepo := newFakePurchaseRepo()
	purRepo.purchases[1] = &models.Purchase{ID: 1, UserID: 3, ProductID: 1, Quantity: 2, Status: models.StatusCompleted}
	svc := NewPurchaseService(purRepo, prodRepo)

	if err := svc.CancelPurchase(context.Background(), 1, 3); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if purRepo.purchases[1].Status != models.StatusCancelled {
		t.Fatalf("expected status cancelled, got %v", purRepo.purchases[1].Status)
	}
	if prodRepo.products[1].Stock != 12 {
		t.Fatalf("expected stock restored to 12, got %d", prodRepo.products[1].Stock)
	}
}

func TestCancelPurchase_Forbidden(t *testing.T) {
	purRepo := newFakePurchaseRepo()
	purRepo.purchases[1] = &models.Purchase{ID: 1, UserID: 3, Status: models.StatusCompleted}
	svc := NewPurchaseService(purRepo, newFakeProductRepo())

	err := svc.CancelPurchase(context.Background(), 1, 99)
	if !errors.Is(err, apperr.ErrPurchaseForbidden) {
		t.Fatalf("expected ErrPurchaseForbidden, got %v", err)
	}
}

func TestCancelPurchase_NotFound(t *testing.T) {
	svc := NewPurchaseService(newFakePurchaseRepo(), newFakeProductRepo())

	err := svc.CancelPurchase(context.Background(), 999, 3)
	if !errors.Is(err, apperr.ErrPurchaseNotFound) {
		t.Fatalf("expected ErrPurchaseNotFound, got %v", err)
	}
}

func TestCancelPurchase_InvalidStatus(t *testing.T) {
	purRepo := newFakePurchaseRepo()
	purRepo.purchases[1] = &models.Purchase{ID: 1, UserID: 3, Status: models.StatusCancelled}
	svc := NewPurchaseService(purRepo, newFakeProductRepo())

	err := svc.CancelPurchase(context.Background(), 1, 3)
	if !errors.Is(err, apperr.ErrPurchaseNotCancellable) {
		t.Fatalf("expected ErrPurchaseNotCancellable, got %v", err)
	}
}
