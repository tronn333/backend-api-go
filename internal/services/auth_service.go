package services

import (
    "errors"
)

type User struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func RegisterUser(user User) error {
    // Implement user registration logic here
    if user.Username == "" || user.Password == "" {
        return errors.New("username and password are required")
    }
    // Save user logic...
    return nil
}