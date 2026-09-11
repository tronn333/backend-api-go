package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	UserIDKey = "user_id"
	EmailKey  = "email"
	RoleKey   = "role"
)

// AuthRequired validates the Bearer JWT token in the Authorization header.
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header format must be: Bearer <token>"})
			return
		}

		tokenStr := parts[1]
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "changeme"
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		// user_id comes from JWT as float64 in MapClaims
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token payload"})
			return
		}

		c.Set(UserIDKey, int(userIDFloat))
		c.Set(EmailKey, claims["email"])
		c.Set(RoleKey, claims["role"])
		c.Next()
	}
}

// AdminRequired ensures the authenticated user has the "admin" role.
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get(RoleKey)
		if role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			return
		}
		c.Next()
	}
}

// GetUserID extracts the user ID from the Gin context (set by AuthRequired).
func GetUserID(c *gin.Context) int {
	id, _ := c.Get(UserIDKey)
	if v, ok := id.(int); ok {
		return v
	}
	return 0
}
