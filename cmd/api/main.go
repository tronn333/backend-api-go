package main

import (
	"backend-api-go/internal/handlers"
	"backend-api-go/internal/repositories"
	"backend-api-go/internal/services"
	"backend-api-go/pkg/db"
	"backend-api-go/pkg/middleware"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file (ignore error in production where env vars may be set directly)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Connect to PostgreSQL
	database := db.Connect()

	// Run migrations
	if err := db.Migrate(database); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	// ── Repositories ──────────────────────────────────────────
	authRepo := repositories.NewAuthRepo(database)
	userRepo := repositories.NewUserRepo(database)
	productRepo := repositories.NewProductRepo(database)
	purchaseRepo := repositories.NewPurchaseRepo(database)

	// ── Services ──────────────────────────────────────────────
	authService := services.NewAuthService(authRepo)
	userService := services.NewUserService(userRepo)
	productService := services.NewProductService(productRepo)
	purchaseService := services.NewPurchaseService(purchaseRepo, productRepo)

	// ── Handlers ──────────────────────────────────────────────
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	productHandler := handlers.NewProductHandler(productService)
	purchaseHandler := handlers.NewPurchaseHandler(purchaseService)

	// ── Router ────────────────────────────────────────────────
	r := gin.Default()

	// Trust only loopback proxies (removes the "trusted all proxies" warning).
	// If you run behind a reverse-proxy (nginx, etc.) add its IP here instead.
	if err := r.SetTrustedProxies([]string{"127.0.0.1", "::1"}); err != nil {
		log.Fatalf("failed to set trusted proxies: %v", err)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Auth (public)
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Products – listing & detail are public; create/update/delete require auth
	r.GET("/products", productHandler.ListProducts)
	r.GET("/products/:id", productHandler.GetProduct)

	// Authenticated routes
	protected := r.Group("/")
	protected.Use(middleware.AuthRequired())
	{
		// User profile
		protected.GET("/users/me", userHandler.GetProfile)
		protected.PATCH("/users/me", userHandler.UpdateProfile)
		protected.DELETE("/users/me", userHandler.DeleteAccount)

		// Products (create / update / delete require login)
		protected.POST("/products", productHandler.CreateProduct)
		protected.PATCH("/products/:id", productHandler.UpdateProduct)
		protected.DELETE("/products/:id", productHandler.DeleteProduct)

		// Purchases / order history
		protected.POST("/purchases", purchaseHandler.CreatePurchase)
		protected.GET("/purchases", purchaseHandler.GetPurchaseHistory)
		protected.GET("/purchases/:id", purchaseHandler.GetPurchase)
		protected.PATCH("/purchases/:id/cancel", purchaseHandler.CancelPurchase)
	}

	// ── Start server ──────────────────────────────────────────
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
