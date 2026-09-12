package repositories

import (
	"backend-api-go/internal/apperr"
	"backend-api-go/internal/models"
	"context"
	"database/sql"
)

type PurchaseRepo interface {
	Create(ctx context.Context, purchase *models.Purchase) error
	GetByUserID(ctx context.Context, userID, page, pageSize int) ([]models.Purchase, int, error)
	GetByID(ctx context.Context, id int) (*models.Purchase, error)
	UpdateStatus(ctx context.Context, id int, status models.PurchaseStatus) error
}

type purchaseRepo struct {
	db *sql.DB
}

// NewPurchaseRepo constructs a PurchaseRepo backed by the given database handle.
func NewPurchaseRepo(db *sql.DB) PurchaseRepo {
	return &purchaseRepo{db: db}
}

// Create inserts a new purchase and writes back its generated id and timestamps.
func (r *purchaseRepo) Create(ctx context.Context, purchase *models.Purchase) error {
	query := `INSERT INTO purchases (user_id, product_id, quantity, total_price, status)
	          VALUES ($1, $2, $3, $4, $5)
	          RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query,
		purchase.UserID, purchase.ProductID, purchase.Quantity,
		purchase.TotalPrice, purchase.Status,
	).Scan(&purchase.ID, &purchase.CreatedAt, &purchase.UpdatedAt)
}

// GetByUserID returns a page of a user's purchases (joined with their product)
// and the total matching count.
func (r *purchaseRepo) GetByUserID(ctx context.Context, userID, page, pageSize int) ([]models.Purchase, int, error) {
	offset := (page - 1) * pageSize

	var count int
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM purchases WHERE user_id = $1`, userID,
	).Scan(&count); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT pu.id, pu.user_id, pu.product_id, pu.quantity, pu.total_price, pu.status,
		        pu.created_at, pu.updated_at,
		        p.id, p.name, p.description, p.price, p.stock, p.category, p.seller_id, p.created_at, p.updated_at
		 FROM purchases pu
		 JOIN products p ON p.id = pu.product_id
		 WHERE pu.user_id = $1
		 ORDER BY pu.created_at DESC
		 LIMIT $2 OFFSET $3`,
		userID, pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var purchases []models.Purchase
	for rows.Next() {
		var pu models.Purchase
		var prod models.Product
		var desc, cat sql.NullString
		if err := rows.Scan(
			&pu.ID, &pu.UserID, &pu.ProductID, &pu.Quantity, &pu.TotalPrice, &pu.Status,
			&pu.CreatedAt, &pu.UpdatedAt,
			&prod.ID, &prod.Name, &desc, &prod.Price, &prod.Stock,
			&cat, &prod.SellerID, &prod.CreatedAt, &prod.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		prod.Description = desc.String
		prod.Category = cat.String
		pu.Product = &prod
		purchases = append(purchases, pu)
	}
	return purchases, count, rows.Err()
}

// GetByID fetches a single purchase (joined with its product) by id, returning
// ErrPurchaseNotFound when it doesn't exist.
func (r *purchaseRepo) GetByID(ctx context.Context, id int) (*models.Purchase, error) {
	pu := &models.Purchase{}
	prod := &models.Product{}
	var desc, cat sql.NullString

	query := `SELECT pu.id, pu.user_id, pu.product_id, pu.quantity, pu.total_price, pu.status,
	                 pu.created_at, pu.updated_at,
	                 p.id, p.name, p.description, p.price, p.stock, p.category, p.seller_id, p.created_at, p.updated_at
	          FROM purchases pu
	          JOIN products p ON p.id = pu.product_id
	          WHERE pu.id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&pu.ID, &pu.UserID, &pu.ProductID, &pu.Quantity, &pu.TotalPrice, &pu.Status,
		&pu.CreatedAt, &pu.UpdatedAt,
		&prod.ID, &prod.Name, &desc, &prod.Price, &prod.Stock,
		&cat, &prod.SellerID, &prod.CreatedAt, &prod.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, apperr.ErrPurchaseNotFound
	}
	if err != nil {
		return nil, err
	}
	prod.Description = desc.String
	prod.Category = cat.String
	pu.Product = prod
	return pu, nil
}

// UpdateStatus sets a purchase's status, returning ErrPurchaseNotFound when no
// row was updated.
func (r *purchaseRepo) UpdateStatus(ctx context.Context, id int, status models.PurchaseStatus) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE purchases SET status = $1, updated_at = NOW() WHERE id = $2`,
		status, id,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return apperr.ErrPurchaseNotFound
	}
	return nil
}
