package repositories

import (
	"backend-api-go/internal/models"
	"context"
	"database/sql"
)

// AuthRepo is now an alias that delegates to UserRepo for auth-specific queries.
// Registration and login both operate on the users table.
type AuthRepo interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type authRepo struct {
	db *sql.DB
}

// NewAuthRepo constructs an AuthRepo backed by the given database handle.
func NewAuthRepo(db *sql.DB) AuthRepo {
	return &authRepo{db: db}
}

// Create inserts a new user and writes back the generated id and timestamps.
func (r *authRepo) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (username, email, password, role)
	          VALUES ($1, $2, $3, $4)
	          RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query,
		user.Username, user.Email, user.Password, user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

// GetByEmail fetches a user by email, returning nil without an error when no
// matching user exists.
func (r *authRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password, role, created_at, updated_at
	          FROM users WHERE email = $1`
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}
