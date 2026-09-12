// Package jwtutil centralises JWT configuration shared by token generation
// (services) and validation (middleware).
package jwtutil

import "os"

// defaultSecret is only used for local development when JWT_SECRET is unset.
// In production you must set JWT_SECRET to a strong, unique value.
const defaultSecret = "changeme"

// Secret returns the JWT signing key. It reads JWT_SECRET from the environment
// and falls back to a development default, ensuring token generation and
// validation always agree on the same key.
func Secret() []byte {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		s = defaultSecret
	}
	return []byte(s)
}
