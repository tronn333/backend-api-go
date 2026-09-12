package services

import (
	"context"
	"errors"
	"testing"

	"backend-api-go/internal/apperr"
	"backend-api-go/internal/models"
)

func TestGetProfile_Success(t *testing.T) {
	repo := newFakeUserRepo()
	repo.users[1] = &models.User{ID: 1, Username: "alice", Email: "alice@example.com"}
	svc := NewUserService(repo)

	u, err := svc.GetProfile(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if u.ID != 1 || u.Username != "alice" {
		t.Fatalf("unexpected user: %+v", u)
	}
}

func TestGetProfile_NotFound(t *testing.T) {
	svc := NewUserService(newFakeUserRepo())

	_, err := svc.GetProfile(context.Background(), 999)
	if !errors.Is(err, apperr.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUpdateProfile_Success(t *testing.T) {
	repo := newFakeUserRepo()
	repo.users[1] = &models.User{ID: 1, Username: "alice", Email: "alice@example.com"}
	svc := NewUserService(repo)

	u, err := svc.UpdateProfile(context.Background(), 1, &models.UpdateProfileRequest{
		Username: "alice2",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if u.Username != "alice2" {
		t.Fatalf("expected username updated, got %q", u.Username)
	}
	if u.Email != "alice@example.com" {
		t.Fatalf("expected email unchanged, got %q", u.Email)
	}
}

func TestUpdateProfile_NotFound(t *testing.T) {
	svc := NewUserService(newFakeUserRepo())

	_, err := svc.UpdateProfile(context.Background(), 999, &models.UpdateProfileRequest{Username: "x"})
	if !errors.Is(err, apperr.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestDeleteAccount_Success(t *testing.T) {
	repo := newFakeUserRepo()
	repo.users[1] = &models.User{ID: 1}
	svc := NewUserService(repo)

	if err := svc.DeleteAccount(context.Background(), 1); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(repo.users) != 0 {
		t.Fatal("expected user to be deleted")
	}
}

func TestDeleteAccount_Error(t *testing.T) {
	repo := newFakeUserRepo()
	repo.users[1] = &models.User{ID: 1}
	repo.deleteErr = errors.New("db down")
	svc := NewUserService(repo)

	if err := svc.DeleteAccount(context.Background(), 1); err == nil {
		t.Fatal("expected an error, got nil")
	}
}
