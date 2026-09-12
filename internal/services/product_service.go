package services

import (
	"backend-api-go/internal/apperr"
	"backend-api-go/internal/models"
	"backend-api-go/internal/repositories"
	"context"
	"fmt"
)

type ProductService interface {
	ListProducts(ctx context.Context, page, pageSize int, category string) (*models.ProductListResponse, error)
	GetProduct(ctx context.Context, id int) (*models.Product, error)
	CreateProduct(ctx context.Context, sellerID int, req *models.CreateProductRequest) (*models.Product, error)
	UpdateProduct(ctx context.Context, id, sellerID int, isAdmin bool, req *models.UpdateProductRequest) (*models.Product, error)
	DeleteProduct(ctx context.Context, id, sellerID int, isAdmin bool) error
}

type productService struct {
	productRepo repositories.ProductRepo
}

// NewProductService constructs a ProductService backed by the given repository.
func NewProductService(productRepo repositories.ProductRepo) ProductService {
	return &productService{productRepo: productRepo}
}

// ListProducts clamps pagination to sane bounds and returns a paginated product listing.
func (s *productService) ListProducts(ctx context.Context, page, pageSize int, category string) (*models.ProductListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	products, total, err := s.productRepo.GetAll(ctx, page, pageSize, category)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}
	if products == nil {
		products = []models.Product{}
	}

	return &models.ProductListResponse{
		Products: products,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// GetProduct returns a single product by id.
func (s *productService) GetProduct(ctx context.Context, id int) (*models.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

// CreateProduct builds a product owned by sellerID and persists it.
func (s *productService) CreateProduct(ctx context.Context, sellerID int, req *models.CreateProductRequest) (*models.Product, error) {
	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
		SellerID:    sellerID,
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}
	return product, nil
}

// UpdateProduct applies partial updates to a product, restricted to its owner or
// an admin.
func (s *productService) UpdateProduct(ctx context.Context, id, sellerID int, isAdmin bool, req *models.UpdateProductRequest) (*models.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Only the seller or an admin can update
	if product.SellerID != sellerID && !isAdmin {
		return nil, apperr.ErrProductForbidden
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.Stock >= 0 {
		product.Stock = req.Stock
	}
	if req.Category != "" {
		product.Category = req.Category
	}

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}
	return product, nil
}

// DeleteProduct removes a product, restricted to its owner or an admin.
func (s *productService) DeleteProduct(ctx context.Context, id, sellerID int, isAdmin bool) error {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Only the seller or an admin can delete
	if product.SellerID != sellerID && !isAdmin {
		return apperr.ErrProductForbidden
	}

	return s.productRepo.Delete(ctx, id)
}
