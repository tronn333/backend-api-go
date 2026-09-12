package services

import (
	"context"
	"errors"
	"testing"

	"backend-api-go/internal/apperr"
	"backend-api-go/internal/models"
)

func TestListProducts_ClampsPagination(t *testing.T) {
	repo := newFakeProductRepo()
	var gotPage, gotPageSize int
	repo.getAllFn = func(_ context.Context, page, pageSize int, _ string) ([]models.Product, int, error) {
		gotPage, gotPageSize = page, pageSize
		return []models.Product{}, 0, nil
	}
	svc := NewProductService(repo)

	resp, err := svc.ListProducts(context.Background(), 0, 500, "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if gotPage != 1 {
		t.Fatalf("expected page clamped to 1, got %d", gotPage)
	}
	if gotPageSize != 20 {
		t.Fatalf("expected pageSize clamped to 20, got %d", gotPageSize)
	}
	if resp.Page != 1 || resp.PageSize != 20 {
		t.Fatalf("unexpected response pagination: %+v", resp)
	}
}

func TestListProducts_EmptySliceNotNil(t *testing.T) {
	repo := newFakeProductRepo()
	repo.getAllFn = func(_ context.Context, _, _ int, _ string) ([]models.Product, int, error) {
		return nil, 0, nil
	}
	svc := NewProductService(repo)

	resp, err := svc.ListProducts(context.Background(), 1, 20, "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Products == nil {
		t.Fatal("expected a non-nil empty slice")
	}
}

func TestGetProduct_Found(t *testing.T) {
	repo := newFakeProductRepo()
	repo.products[1] = &models.Product{ID: 1, Name: "Widget", Price: 9.99, Stock: 5}
	svc := NewProductService(repo)

	p, err := svc.GetProduct(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p.Name != "Widget" {
		t.Fatalf("unexpected product: %+v", p)
	}
}

func TestGetProduct_NotFound(t *testing.T) {
	svc := NewProductService(newFakeProductRepo())

	_, err := svc.GetProduct(context.Background(), 999)
	if !errors.Is(err, apperr.ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
}

func TestCreateProduct_SetsSeller(t *testing.T) {
	repo := newFakeProductRepo()
	svc := NewProductService(repo)

	p, err := svc.CreateProduct(context.Background(), 42, &models.CreateProductRequest{
		Name:  "Widget",
		Price: 10,
		Stock: 3,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p.SellerID != 42 {
		t.Fatalf("expected seller id 42, got %d", p.SellerID)
	}
	if p.ID == 0 {
		t.Fatal("expected product id to be assigned")
	}
}

func TestUpdateProduct_Owner(t *testing.T) {
	repo := newFakeProductRepo()
	repo.products[1] = &models.Product{ID: 1, SellerID: 7, Name: "Old", Price: 5, Stock: 2}
	svc := NewProductService(repo)

	p, err := svc.UpdateProduct(context.Background(), 1, 7, false, &models.UpdateProductRequest{Name: "New"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p.Name != "New" {
		t.Fatalf("expected name updated, got %q", p.Name)
	}
}

func TestUpdateProduct_Forbidden(t *testing.T) {
	repo := newFakeProductRepo()
	repo.products[1] = &models.Product{ID: 1, SellerID: 7}
	svc := NewProductService(repo)

	_, err := svc.UpdateProduct(context.Background(), 1, 99, false, &models.UpdateProductRequest{Name: "New"})
	if !errors.Is(err, apperr.ErrProductForbidden) {
		t.Fatalf("expected ErrProductForbidden, got %v", err)
	}
}

func TestUpdateProduct_AdminBypass(t *testing.T) {
	repo := newFakeProductRepo()
	repo.products[1] = &models.Product{ID: 1, SellerID: 7, Name: "Old"}
	svc := NewProductService(repo)

	p, err := svc.UpdateProduct(context.Background(), 1, 99, true, &models.UpdateProductRequest{Name: "New"})
	if err != nil {
		t.Fatalf("expected admin bypass, got %v", err)
	}
	if p.Name != "New" {
		t.Fatalf("expected name updated, got %q", p.Name)
	}
}

func TestUpdateProduct_NotFound(t *testing.T) {
	svc := NewProductService(newFakeProductRepo())

	_, err := svc.UpdateProduct(context.Background(), 999, 7, false, &models.UpdateProductRequest{Name: "New"})
	if !errors.Is(err, apperr.ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
}

func TestDeleteProduct_Owner(t *testing.T) {
	repo := newFakeProductRepo()
	repo.products[1] = &models.Product{ID: 1, SellerID: 7}
	svc := NewProductService(repo)

	if err := svc.DeleteProduct(context.Background(), 1, 7, false); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(repo.products) != 0 {
		t.Fatal("expected product to be deleted")
	}
}

func TestDeleteProduct_Forbidden(t *testing.T) {
	repo := newFakeProductRepo()
	repo.products[1] = &models.Product{ID: 1, SellerID: 7}
	svc := NewProductService(repo)

	err := svc.DeleteProduct(context.Background(), 1, 99, false)
	if !errors.Is(err, apperr.ErrProductForbidden) {
		t.Fatalf("expected ErrProductForbidden, got %v", err)
	}
}

func TestDeleteProduct_AdminBypass(t *testing.T) {
	repo := newFakeProductRepo()
	repo.products[1] = &models.Product{ID: 1, SellerID: 7}
	svc := NewProductService(repo)

	if err := svc.DeleteProduct(context.Background(), 1, 99, true); err != nil {
		t.Fatalf("expected admin bypass, got %v", err)
	}
	if len(repo.products) != 0 {
		t.Fatal("expected product to be deleted")
	}
}
