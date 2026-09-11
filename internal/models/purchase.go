package models

import "time"

type PurchaseStatus string

const (
	StatusPending   PurchaseStatus = "pending"
	StatusCompleted PurchaseStatus = "completed"
	StatusCancelled PurchaseStatus = "cancelled"
)

type Purchase struct {
	ID         int            `json:"id"`
	UserID     int            `json:"user_id"`
	ProductID  int            `json:"product_id"`
	Quantity   int            `json:"quantity"`
	TotalPrice float64        `json:"total_price"`
	Status     PurchaseStatus `json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	// Populated via JOIN
	Product *Product `json:"product,omitempty"`
}

type CreatePurchaseRequest struct {
	ProductID int `json:"product_id" binding:"required,gt=0"`
	Quantity  int `json:"quantity" binding:"required,gt=0"`
}

type PurchaseHistoryResponse struct {
	Purchases []Purchase `json:"purchases"`
	Total     int        `json:"total"`
	Page      int        `json:"page"`
	PageSize  int        `json:"page_size"`
}
