package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var database *sql.DB

// Connect initialises the PostgreSQL connection using environment variables.
func Connect() *sql.DB {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	if dbPort == "" {
		dbPort = "5432"
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	var err error
	database, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open database connection: %v", err)
	}

	// Verify connectivity
	if err = database.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	database.SetMaxOpenConns(25)
	database.SetMaxIdleConns(10)

	log.Println("Successfully connected to PostgreSQL database!")
	return database
}

// GetDB returns the active database handle.
func GetDB() *sql.DB {
	return database
}

// Migrate runs the initial DDL to create all tables if they don't exist.
func Migrate(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id         SERIAL PRIMARY KEY,
			username   VARCHAR(50)  NOT NULL UNIQUE,
			email      VARCHAR(255) NOT NULL UNIQUE,
			password   VARCHAR(255) NOT NULL,
			role       VARCHAR(20)  NOT NULL DEFAULT 'user',
			created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS products (
			id          SERIAL PRIMARY KEY,
			name        VARCHAR(255)   NOT NULL,
			description TEXT,
			price       NUMERIC(12, 2) NOT NULL CHECK (price > 0),
			stock       INT            NOT NULL DEFAULT 0 CHECK (stock >= 0),
			category    VARCHAR(100),
			seller_id   INT            NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
			updated_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS purchases (
			id          SERIAL PRIMARY KEY,
			user_id     INT            NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			product_id  INT            NOT NULL REFERENCES products(id) ON DELETE CASCADE,
			quantity    INT            NOT NULL CHECK (quantity > 0),
			total_price NUMERIC(12, 2) NOT NULL,
			status      VARCHAR(20)    NOT NULL DEFAULT 'pending',
			created_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
			updated_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW()
		)`,
		// Indexes for common queries
		`CREATE INDEX IF NOT EXISTS idx_products_seller_id  ON products(seller_id)`,
		`CREATE INDEX IF NOT EXISTS idx_products_category   ON products(category)`,
		`CREATE INDEX IF NOT EXISTS idx_purchases_user_id   ON purchases(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_purchases_product_id ON purchases(product_id)`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return fmt.Errorf("migration error: %w", err)
		}
	}

	log.Println("Database migrations applied successfully!")
	return nil
}
