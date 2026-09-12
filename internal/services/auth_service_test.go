package services

import (
	"context"
	"errors"
	"testing"

	"backend-api-go/internal/apperr"
	"backend-api-go/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func TestRegister_Success(t *testing.T) {
	repo := newFakeAuthRepo()
	svc := NewAuthService(repo)

	resp, err := svc.Register(context.Background(), &models.RegisterRequest{
		Username: "alice",
		Email:    "alice@example.com",
		Password: "secret123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Token == "" {
		t.Fatal("expected a non-empty token")
	}
	if resp.User.Email != "alice@example.com" {
		t.Fatalf("unexpected email: %s", resp.User.Email)
	}
	if resp.User.Role != "user" {
		t.Fatalf("expected default role 'user', got %q", resp.User.Role)
	}
	if resp.User.Password == "secret123" {
		t.Fatal("expected password to be hashed, got plaintext")
	}
	if len(repo.users) != 1 {
		t.Fatalf("expected 1 stored user, got %d", len(repo.users))
	}
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	repo := newFakeAuthRepo()
	repo.users["alice@example.com"] = &models.User{ID: 1, Email: "alice@example.com"}
	svc := NewAuthService(repo)

	_, err := svc.Register(context.Background(), &models.RegisterRequest{
		Username: "alice",
		Email:    "alice@example.com",
		Password: "secret123",
	})
	if !errors.Is(err, apperr.ErrEmailAlreadyRegistered) {
		t.Fatalf("expected ErrEmailAlreadyRegistered, got %v", err)
	}
}

func TestLogin_Success(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	repo := newFakeAuthRepo()
	repo.users["alice@example.com"] = &models.User{
		ID:       1,
		Username: "alice",
		Email:    "alice@example.com",
		Password: string(hash),
		Role:     "user",
	}
	svc := NewAuthService(repo)

	resp, err := svc.Login(context.Background(), &models.LoginRequest{
		Email:    "alice@example.com",
		Password: "secret123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Token == "" {
		t.Fatal("expected a non-empty token")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.MinCost)
	repo := newFakeAuthRepo()
	repo.users["alice@example.com"] = &models.User{
		ID:       1,
		Email:    "alice@example.com",
		Password: string(hash),
	}
	svc := NewAuthService(repo)

	_, err := svc.Login(context.Background(), &models.LoginRequest{
		Email:    "alice@example.com",
		Password: "wrong-password",
	})
	if !errors.Is(err, apperr.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_UnknownEmail(t *testing.T) {
	svc := NewAuthService(newFakeAuthRepo())

	_, err := svc.Login(context.Background(), &models.LoginRequest{
		Email:    "nobody@example.com",
		Password: "secret123",
	})
	if !errors.Is(err, apperr.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}
