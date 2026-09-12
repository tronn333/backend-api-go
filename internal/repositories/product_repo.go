package repositories

import (
	"backend-api-go/internal/apperr"
	"backend-api-go/internal/models"
	"context"
	"database/sql"
)

type ProductRepo interface {
	GetAll(ctx context.Context, page, pageSize int, category string) ([]models.Product, int, error)
	GetByID(ctx context.Context, id int) (*models.Product, error)
	Create(ctx context.Context, product *models.Product) error
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id int) error
	DecrementStock(ctx context.Context, id, quantity int) error
}

type productRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) ProductRepo {
	return &productRepo{db: db}
}

func (r *productRepo) GetAll(ctx context.Context, page, pageSize int, category string) ([]models.Product, int, error) {
	offset := (page - 1) * pageSize

	var (
		rows  *sql.Rows
		count int
		err   error
	)

	if category != "" {
		err = r.db.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM products WHERE category = $1`, category,
		).Scan(&count)
		if err != nil {
			return nil, 0, err
		}
		rows, err = r.db.QueryContext(ctx,
			`SELECT id, name, description, price, stock, category, seller_id, created_at, updated_at
			 FROM products WHERE category = $1
			 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
			category, pageSize, offset,
		)
	} else {
		err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM products`).Scan(&count)
		if err != nil {
			return nil, 0, err
		}
		rows, err = r.db.QueryContext(ctx,
			`SELECT id, name, description, price, stock, category, seller_id, created_at, updated_at
			 FROM products
			 ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
			pageSize, offset,
		)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		var desc, cat sql.NullString
		if err := rows.Scan(
			&p.ID, &p.Name, &desc, &p.Price, &p.Stock,
			&cat, &p.SellerID, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		p.Description = desc.String
		p.Category = cat.String
		products = append(products, p)
	}
	return products, count, rows.Err()
}

func (r *productRepo) GetByID(ctx context.Context, id int) (*models.Product, error) {
	p := &models.Product{}
	var desc, cat sql.NullString
	query := `SELECT id, name, description, price, stock, category, seller_id, created_at, updated_at
	          FROM products WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Name, &desc, &p.Price, &p.Stock,
		&cat, &p.SellerID, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, apperr.ErrProductNotFound
	}
	if err != nil {
		return nil, err
	}
	p.Description = desc.String
	p.Category = cat.String
	return p, nil
}

func (r *productRepo) Create(ctx context.Context, product *models.Product) error {
	query := `INSERT INTO products (name, description, price, stock, category, seller_id)
	          VALUES ($1, $2, $3, $4, $5, $6)
	          RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query,
		product.Name, product.Description, product.Price,
		product.Stock, product.Category, product.SellerID,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
}

func (r *productRepo) Update(ctx context.Context, product *models.Product) error {
	query := `UPDATE products
	          SET name = $1, description = $2, price = $3, stock = $4, category = $5, updated_at = NOW()
	          WHERE id = $6
	          RETURNING updated_at`
	return r.db.QueryRowContext(ctx, query,
		product.Name, product.Description, product.Price,
		product.Stock, product.Category, product.ID,
	).Scan(&product.UpdatedAt)
}

func (r *productRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM products WHERE id = $1`, id)
	return err
}

func (r *productRepo) DecrementStock(ctx context.Context, id, quantity int) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE products SET stock = stock - $1, updated_at = NOW()
		 WHERE id = $2 AND stock >= $1`,
		quantity, id,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return apperr.ErrInsufficientStock
	}
	return nil
}
