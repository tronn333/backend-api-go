package main

import (
    "net/http"
    "my-ecommerce-backend/internal/handlers"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/register", handlers.RegisterHandler)
    mux.HandleFunc("/login", handlers.LoginHandler)

    http.ListenAndServe(":8080", mux)
}