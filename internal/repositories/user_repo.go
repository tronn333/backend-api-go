package repositories

import (
	"backend-api-go/internal/models"
	"context"
	"database/sql"
	"fmt"
)

type UserRepo interface {
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int) error
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) GetByID(ctx context.Context, id int) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password, role, created_at, updated_at
	          FROM users WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password, role, created_at, updated_at
	          FROM users WHERE email = $1`
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepo) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (username, email, password, role)
	          VALUES ($1, $2, $3, $4)
	          RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query,
		user.Username, user.Email, user.Password, user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *userRepo) Update(ctx context.Context, user *models.User) error {
	query := `UPDATE users
	          SET username = $1, email = $2, updated_at = NOW()
	          WHERE id = $3
	          RETURNING updated_at`
	return r.db.QueryRowContext(ctx, query,
		user.Username, user.Email, user.ID,
	).Scan(&user.UpdatedAt)
}

func (r *userRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	return err
}
