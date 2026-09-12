package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"backend-api-go/internal/apperr"
	"backend-api-go/internal/models"
	"backend-api-go/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func newUserContext(t *testing.T, method, path, body string, userID int) (*httptest.ResponseRecorder, *gin.Context) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		c.Request.Header.Set("Content-Type", "application/json")
	}
	c.Set(middleware.UserIDKey, userID)
	return w, c
}

func TestGetProfile_Success(t *testing.T) {
	mock := &mockUserService{}
	mock.getProfileFn = func(ctx context.Context, userID int) (*models.User, error) {
		return &models.User{ID: userID, Username: "alice"}, nil
	}
	h := NewUserHandler(mock)

	w, c := newUserContext(t, http.MethodGet, "/users/me", "", 1)
	h.GetProfile(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetProfile_NotFound(t *testing.T) {
	mock := &mockUserService{}
	mock.getProfileFn = func(ctx context.Context, userID int) (*models.User, error) {
		return nil, apperr.ErrUserNotFound
	}
	h := NewUserHandler(mock)

	w, c := newUserContext(t, http.MethodGet, "/users/me", "", 1)
	h.GetProfile(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUpdateProfile_Success(t *testing.T) {
	mock := &mockUserService{}
	mock.updateProfileFn = func(ctx context.Context, userID int, req *models.UpdateProfileRequest) (*models.User, error) {
		return &models.User{ID: userID, Username: req.Username}, nil
	}
	h := NewUserHandler(mock)

	w, c := newUserContext(t, http.MethodPatch, "/users/me", `{"username":"alice2"}`, 1)
	h.UpdateProfile(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestUpdateProfile_BadRequest(t *testing.T) {
	h := NewUserHandler(&mockUserService{})

	w, c := newUserContext(t, http.MethodPatch, "/users/me", `{"username":"ab"}`, 1)
	h.UpdateProfile(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUpdateProfile_NotFound(t *testing.T) {
	mock := &mockUserService{}
	mock.updateProfileFn = func(ctx context.Context, userID int, req *models.UpdateProfileRequest) (*models.User, error) {
		return nil, apperr.ErrUserNotFound
	}
	h := NewUserHandler(mock)

	w, c := newUserContext(t, http.MethodPatch, "/users/me", `{"username":"alice2"}`, 1)
	h.UpdateProfile(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestDeleteAccount_Success(t *testing.T) {
	mock := &mockUserService{}
	mock.deleteAccountFn = func(ctx context.Context, userID int) error { return nil }
	h := NewUserHandler(mock)

	w, c := newUserContext(t, http.MethodDelete, "/users/me", "", 1)
	h.DeleteAccount(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
