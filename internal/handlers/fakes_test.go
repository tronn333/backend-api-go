package handlers

import (
	"context"

	"backend-api-go/internal/models"
)

// ── Auth service mock ───────────────────────────────────────────

type mockAuthService struct {
	registerFn func(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error)
	loginFn    func(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error)
}

func (m *mockAuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	if m.registerFn != nil {
		return m.registerFn(ctx, req)
	}
	return nil, nil
}

func (m *mockAuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	if m.loginFn != nil {
		return m.loginFn(ctx, req)
	}
	return nil, nil
}

// ── User service mock ───────────────────────────────────────────

type mockUserService struct {
	getProfileFn    func(ctx context.Context, userID int) (*models.User, error)
	updateProfileFn func(ctx context.Context, userID int, req *models.UpdateProfileRequest) (*models.User, error)
	deleteAccountFn func(ctx context.Context, userID int) error
}

func (m *mockUserService) GetProfile(ctx context.Context, userID int) (*models.User, error) {
	if m.getProfileFn != nil {
		return m.getProfileFn(ctx, userID)
	}
	return &models.User{}, nil
}

func (m *mockUserService) UpdateProfile(ctx context.Context, userID int, req *models.UpdateProfileRequest) (*models.User, error) {
	if m.updateProfileFn != nil {
		return m.updateProfileFn(ctx, userID, req)
	}
	return &models.User{}, nil
}

func (m *mockUserService) DeleteAccount(ctx context.Context, userID int) error {
	if m.deleteAccountFn != nil {
		return m.deleteAccountFn(ctx, userID)
	}
	return nil
}

// ── Product service mock ────────────────────────────────────────

type mockProductService struct {
	listProductsFn  func(ctx context.Context, page, pageSize int, category string) (*models.ProductListResponse, error)
	getProductFn    func(ctx context.Context, id int) (*models.Product, error)
	createProductFn func(ctx context.Context, sellerID int, req *models.CreateProductRequest) (*models.Product, error)
	updateProductFn func(ctx context.Context, id, sellerID int, isAdmin bool, req *models.UpdateProductRequest) (*models.Product, error)
	deleteProductFn func(ctx context.Context, id, sellerID int, isAdmin bool) error
}

func (m *mockProductService) ListProducts(ctx context.Context, page, pageSize int, category string) (*models.ProductListResponse, error) {
	if m.listProductsFn != nil {
		return m.listProductsFn(ctx, page, pageSize, category)
	}
	return &models.ProductListResponse{}, nil
}

func (m *mockProductService) GetProduct(ctx context.Context, id int) (*models.Product, error) {
	if m.getProductFn != nil {
		return m.getProductFn(ctx, id)
	}
	return &models.Product{}, nil
}

func (m *mockProductService) CreateProduct(ctx context.Context, sellerID int, req *models.CreateProductRequest) (*models.Product, error) {
	if m.createProductFn != nil {
		return m.createProductFn(ctx, sellerID, req)
	}
	return &models.Product{}, nil
}

func (m *mockProductService) UpdateProduct(ctx context.Context, id, sellerID int, isAdmin bool, req *models.UpdateProductRequest) (*models.Product, error) {
	if m.updateProductFn != nil {
		return m.updateProductFn(ctx, id, sellerID, isAdmin, req)
	}
	return &models.Product{}, nil
}

func (m *mockProductService) DeleteProduct(ctx context.Context, id, sellerID int, isAdmin bool) error {
	if m.deleteProductFn != nil {
		return m.deleteProductFn(ctx, id, sellerID, isAdmin)
	}
	return nil
}

// ── Purchase service mock ───────────────────────────────────────

type mockPurchaseService struct {
	createPurchaseFn     func(ctx context.Context, userID int, req *models.CreatePurchaseRequest) (*models.Purchase, error)
	getPurchaseHistoryFn func(ctx context.Context, userID, page, pageSize int) (*models.PurchaseHistoryResponse, error)
	getPurchaseFn        func(ctx context.Context, id, userID int) (*models.Purchase, error)
	cancelPurchaseFn     func(ctx context.Context, id, userID int) error
}

func (m *mockPurchaseService) CreatePurchase(ctx context.Context, userID int, req *models.CreatePurchaseRequest) (*models.Purchase, error) {
	if m.createPurchaseFn != nil {
		return m.createPurchaseFn(ctx, userID, req)
	}
	return &models.Purchase{}, nil
}

func (m *mockPurchaseService) GetPurchaseHistory(ctx context.Context, userID, page, pageSize int) (*models.PurchaseHistoryResponse, error) {
	if m.getPurchaseHistoryFn != nil {
		return m.getPurchaseHistoryFn(ctx, userID, page, pageSize)
	}
	return &models.PurchaseHistoryResponse{}, nil
}

func (m *mockPurchaseService) GetPurchase(ctx context.Context, id, userID int) (*models.Purchase, error) {
	if m.getPurchaseFn != nil {
		return m.getPurchaseFn(ctx, id, userID)
	}
	return &models.Purchase{}, nil
}

func (m *mockPurchaseService) CancelPurchase(ctx context.Context, id, userID int) error {
	if m.cancelPurchaseFn != nil {
		return m.cancelPurchaseFn(ctx, id, userID)
	}
	return nil
}
