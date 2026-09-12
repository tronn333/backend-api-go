// Package apperr defines sentinel errors for domain-level failures.
// Services and repositories return these, and HTTP handlers map them to
// status codes using errors.Is.
package apperr

import "errors"

var (
	// Auth
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrInvalidCredentials     = errors.New("invalid email or password")

	// Users
	ErrUserNotFound = errors.New("user not found")

	// Products
	ErrProductNotFound  = errors.New("product not found")
	ErrProductForbidden = errors.New("forbidden: you do not own this product")

	// Purchases
	ErrInsufficientStock      = errors.New("insufficient stock")
	ErrPurchaseNotFound       = errors.New("purchase not found")
	ErrPurchaseForbidden      = errors.New("forbidden: you do not own this purchase")
	ErrPurchaseNotCancellable = errors.New("purchase cannot be cancelled in its current status")
)
