package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"backend-api-go/pkg/jwtutil"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func newTestToken(t *testing.T, userID float64, role string) string {
	t.Helper()
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   "alice@example.com",
		"role":    role,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := tok.SignedString(jwtutil.Secret())
	if err != nil {
		t.Fatalf("failed to sign test token: %v", err)
	}
	return s
}

func newAuthContext(token string) (*httptest.ResponseRecorder, *gin.Context) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	if token != "" {
		c.Request.Header.Set("Authorization", "Bearer "+token)
	}
	return w, c
}

func TestAuthRequired_NoHeader(t *testing.T) {
	w, c := newAuthContext("")

	AuthRequired()(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthRequired_BadFormat(t *testing.T) {
	w, c := newAuthContext("not-a-real-token")
	c.Request.Header.Set("Authorization", "Token abc123")

	AuthRequired()(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthRequired_InvalidToken(t *testing.T) {
	w, c := newAuthContext("this-is-not-a-jwt")

	AuthRequired()(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthRequired_ValidToken(t *testing.T) {
	w, c := newAuthContext(newTestToken(t, 42, "admin"))

	AuthRequired()(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if got := GetUserID(c); got != 42 {
		t.Fatalf("expected user id 42, got %d", got)
	}
	if !IsAdmin(c) {
		t.Fatal("expected IsAdmin true for admin role")
	}
	if email, _ := c.Get(EmailKey); email != "alice@example.com" {
		t.Fatalf("unexpected email: %v", email)
	}
}

func TestAuthRequired_ValidTokenUserRole(t *testing.T) {
	_, c := newAuthContext(newTestToken(t, 1, "user"))

	AuthRequired()(c)

	if IsAdmin(c) {
		t.Fatal("expected IsAdmin false for user role")
	}
	if got := GetUserID(c); got != 1 {
		t.Fatalf("expected user id 1, got %d", got)
	}
}

func TestGetUserID_NotSet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	if got := GetUserID(c); got != 0 {
		t.Fatalf("expected 0 when user id is not set, got %d", got)
	}
}
