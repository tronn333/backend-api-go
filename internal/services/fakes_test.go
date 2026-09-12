package services

import (
	"context"

	"backend-api-go/internal/apperr"
	"backend-api-go/internal/models"
)

// ── Auth repository fake ────────────────────────────────────────

type fakeAuthRepo struct {
	users map[string]*models.User
}

func newFakeAuthRepo() *fakeAuthRepo {
	return &fakeAuthRepo{users: map[string]*models.User{}}
}

func (f *fakeAuthRepo) Create(_ context.Context, u *models.User) error {
	if u.ID == 0 {
		u.ID = len(f.users) + 1
	}
	f.users[u.Email] = u
	return nil
}

func (f *fakeAuthRepo) GetByEmail(_ context.Context, email string) (*models.User, error) {
	u, ok := f.users[email]
	if !ok {
		return nil, nil
	}
	return u, nil
}

// ── User repository fake ────────────────────────────────────────

type fakeUserRepo struct {
	users     map[int]*models.User
	getErr    error
	updateErr error
	deleteErr error
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{users: map[int]*models.User{}}
}

func (f *fakeUserRepo) GetByID(_ context.Context, id int) (*models.User, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	u, ok := f.users[id]
	if !ok {
		return nil, apperr.ErrUserNotFound
	}
	return u, nil
}

func (f *fakeUserRepo) GetByEmail(_ context.Context, email string) (*models.User, error) {
	for _, u := range f.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, apperr.ErrUserNotFound
}

func (f *fakeUserRepo) Create(_ context.Context, u *models.User) error {
	if u.ID == 0 {
		u.ID = len(f.users) + 1
	}
	f.users[u.ID] = u
	return nil
}

func (f *fakeUserRepo) Update(_ context.Context, u *models.User) error {
	if f.updateErr != nil {
		return f.updateErr
	}
	f.users[u.ID] = u
	return nil
}

func (f *fakeUserRepo) Delete(_ context.Context, id int) error {
	if f.deleteErr != nil {
		return f.deleteErr
	}
	delete(f.users, id)
	return nil
}

// ── Product repository fake ─────────────────────────────────────

type fakeProductRepo struct {
	products  map[int]*models.Product
	getAllFn  func(ctx context.Context, page, pageSize int, category string) ([]models.Product, int, error)
	getErr    error
	createErr error
	updateErr error
	deleteErr error
	decErr    error
}

func newFakeProductRepo() *fakeProductRepo {
	return &fakeProductRepo{products: map[int]*models.Product{}}
}

func (f *fakeProductRepo) GetAll(ctx context.Context, page, pageSize int, category string) ([]models.Product, int, error) {
	if f.getAllFn != nil {
		return f.getAllFn(ctx, page, pageSize, category)
	}
	out := make([]models.Product, 0, len(f.products))
	for _, p := range f.products {
		out = append(out, *p)
	}
	return out, len(out), nil
}

func (f *fakeProductRepo) GetByID(_ context.Context, id int) (*models.Product, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	p, ok := f.products[id]
	if !ok {
		return nil, apperr.ErrProductNotFound
	}
	return p, nil
}

func (f *fakeProductRepo) Create(_ context.Context, p *models.Product) error {
	if f.createErr != nil {
		return f.createErr
	}
	if p.ID == 0 {
		p.ID = len(f.products) + 1
	}
	f.products[p.ID] = p
	return nil
}

func (f *fakeProductRepo) Update(_ context.Context, p *models.Product) error {
	if f.updateErr != nil {
		return f.updateErr
	}
	f.products[p.ID] = p
	return nil
}

func (f *fakeProductRepo) Delete(_ context.Context, id int) error {
	if f.deleteErr != nil {
		return f.deleteErr
	}
	delete(f.products, id)
	return nil
}

func (f *fakeProductRepo) DecrementStock(_ context.Context, id, quantity int) error {
	if f.decErr != nil {
		return f.decErr
	}
	p, ok := f.products[id]
	if !ok {
		return apperr.ErrProductNotFound
	}
	if p.Stock < quantity {
		return apperr.ErrInsufficientStock
	}
	p.Stock -= quantity
	return nil
}

// ── Purchase repository fake ────────────────────────────────────

type fakePurchaseRepo struct {
	purchases       map[int]*models.Purchase
	createErr       error
	getByIDErr      error
	updateStatusErr error
	historyFn       func(ctx context.Context, userID, page, pageSize int) ([]models.Purchase, int, error)
}

func newFakePurchaseRepo() *fakePurchaseRepo {
	return &fakePurchaseRepo{purchases: map[int]*models.Purchase{}}
}

func (f *fakePurchaseRepo) Create(_ context.Context, pu *models.Purchase) error {
	if f.createErr != nil {
		return f.createErr
	}
	if pu.ID == 0 {
		pu.ID = len(f.purchases) + 1
	}
	f.purchases[pu.ID] = pu
	return nil
}

func (f *fakePurchaseRepo) GetByUserID(ctx context.Context, userID, page, pageSize int) ([]models.Purchase, int, error) {
	if f.historyFn != nil {
		return f.historyFn(ctx, userID, page, pageSize)
	}
	var out []models.Purchase
	for _, pu := range f.purchases {
		if pu.UserID == userID {
			out = append(out, *pu)
		}
	}
	return out, len(out), nil
}

func (f *fakePurchaseRepo) GetByID(_ context.Context, id int) (*models.Purchase, error) {
	if f.getByIDErr != nil {
		return nil, f.getByIDErr
	}
	pu, ok := f.purchases[id]
	if !ok {
		return nil, apperr.ErrPurchaseNotFound
	}
	return pu, nil
}

func (f *fakePurchaseRepo) UpdateStatus(_ context.Context, id int, status models.PurchaseStatus) error {
	if f.updateStatusErr != nil {
		return f.updateStatusErr
	}
	pu, ok := f.purchases[id]
	if !ok {
		return apperr.ErrPurchaseNotFound
	}
	pu.Status = status
	return nil
}
