package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"backend-api-go/internal/apperr"
	"backend-api-go/internal/models"

	"github.com/gin-gonic/gin"
)

func newAuthContext(t *testing.T, method, path, body string) (*httptest.ResponseRecorder, *gin.Context) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, path, strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	return w, c
}

func TestRegister_Success(t *testing.T) {
	mock := &mockAuthService{}
	mock.registerFn = func(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
		return &models.AuthResponse{Token: "tok", User: models.User{ID: 1, Email: req.Email}}, nil
	}
	h := NewAuthHandler(mock)

	w, c := newAuthContext(t, http.MethodPost, "/auth/register", `{"username":"alice","email":"a@example.com","password":"secret123"}`)
	h.Register(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), `"token"`) {
		t.Fatalf("expected token in body, got %s", w.Body.String())
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	mock := &mockAuthService{}
	mock.registerFn = func(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
		return nil, apperr.ErrEmailAlreadyRegistered
	}
	h := NewAuthHandler(mock)

	w, c := newAuthContext(t, http.MethodPost, "/auth/register", `{"username":"alice","email":"a@example.com","password":"secret123"}`)
	h.Register(c)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", w.Code)
	}
}

func TestRegister_BadRequest(t *testing.T) {
	h := NewAuthHandler(&mockAuthService{})

	w, c := newAuthContext(t, http.MethodPost, "/auth/register", `{"username":"alice"}`)
	h.Register(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestLogin_Success(t *testing.T) {
	mock := &mockAuthService{}
	mock.loginFn = func(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
		return &models.AuthResponse{Token: "tok"}, nil
	}
	h := NewAuthHandler(mock)

	w, c := newAuthContext(t, http.MethodPost, "/auth/login", `{"email":"a@example.com","password":"secret123"}`)
	h.Login(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	mock := &mockAuthService{}
	mock.loginFn = func(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
		return nil, apperr.ErrInvalidCredentials
	}
	h := NewAuthHandler(mock)

	w, c := newAuthContext(t, http.MethodPost, "/auth/login", `{"email":"a@example.com","password":"wrong"}`)
	h.Login(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestLogin_InternalError(t *testing.T) {
	mock := &mockAuthService{}
	mock.loginFn = func(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
		return nil, errors.New("db down")
	}
	h := NewAuthHandler(mock)

	w, c := newAuthContext(t, http.MethodPost, "/auth/login", `{"email":"a@example.com","password":"secret123"}`)
	h.Login(c)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
