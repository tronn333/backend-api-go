package services

import (
	"backend-api-go/internal/models"
	"backend-api-go/internal/repositories"
	"context"
	"fmt"
)

type PurchaseService interface {
	CreatePurchase(ctx context.Context, userID int, req *models.CreatePurchaseRequest) (*models.Purchase, error)
	GetPurchaseHistory(ctx context.Context, userID, page, pageSize int) (*models.PurchaseHistoryResponse, error)
	GetPurchase(ctx context.Context, id, userID int) (*models.Purchase, error)
	CancelPurchase(ctx context.Context, id, userID int) error
}

type purchaseService struct {
	purchaseRepo repositories.PurchaseRepo
	productRepo  repositories.ProductRepo
}

func NewPurchaseService(purchaseRepo repositories.PurchaseRepo, productRepo repositories.ProductRepo) PurchaseService {
	return &purchaseService{
		purchaseRepo: purchaseRepo,
		productRepo:  productRepo,
	}
}

func (s *purchaseService) CreatePurchase(ctx context.Context, userID int, req *models.CreatePurchaseRequest) (*models.Purchase, error) {
	// Fetch the product to validate and calculate total
	product, err := s.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found")
	}

	if product.Stock < req.Quantity {
		return nil, fmt.Errorf("insufficient stock: only %d items available", product.Stock)
	}

	// Decrement stock atomically (checks stock >= quantity)
	if err := s.productRepo.DecrementStock(ctx, req.ProductID, req.Quantity); err != nil {
		return nil, err
	}

	purchase := &models.Purchase{
		UserID:     userID,
		ProductID:  req.ProductID,
		Quantity:   req.Quantity,
		TotalPrice: product.Price * float64(req.Quantity),
		Status:     models.StatusCompleted,
	}

	if err := s.purchaseRepo.Create(ctx, purchase); err != nil {
		// Attempt to roll back stock decrement
		_ = s.productRepo.DecrementStock(ctx, req.ProductID, -req.Quantity)
		return nil, fmt.Errorf("failed to create purchase: %w", err)
	}

	purchase.Product = product
	return purchase, nil
}

func (s *purchaseService) GetPurchaseHistory(ctx context.Context, userID, page, pageSize int) (*models.PurchaseHistoryResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	purchases, total, err := s.purchaseRepo.GetByUserID(ctx, userID, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch purchase history: %w", err)
	}
	if purchases == nil {
		purchases = []models.Purchase{}
	}

	return &models.PurchaseHistoryResponse{
		Purchases: purchases,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
	}, nil
}

func (s *purchaseService) GetPurchase(ctx context.Context, id, userID int) (*models.Purchase, error) {
	purchase, err := s.purchaseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Users can only see their own purchases
	if purchase.UserID != userID {
		return nil, fmt.Errorf("forbidden")
	}

	return purchase, nil
}

func (s *purchaseService) CancelPurchase(ctx context.Context, id, userID int) error {
	purchase, err := s.purchaseRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if purchase.UserID != userID {
		return fmt.Errorf("forbidden")
	}

	if purchase.Status != models.StatusPending && purchase.Status != models.StatusCompleted {
		return fmt.Errorf("purchase cannot be cancelled in its current status")
	}

	if err := s.purchaseRepo.UpdateStatus(ctx, id, models.StatusCancelled); err != nil {
		return fmt.Errorf("failed to cancel purchase: %w", err)
	}

	// Restore stock on cancellation
	if err := s.productRepo.DecrementStock(ctx, purchase.ProductID, -purchase.Quantity); err != nil {
		// Log but don't fail — status is already cancelled
		fmt.Printf("warning: failed to restore stock for product %d: %v\n", purchase.ProductID, err)
	}

	return nil
}
