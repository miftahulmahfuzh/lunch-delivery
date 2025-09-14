// internal/models/models.go
package models

import (
	"time"

	"github.com/lib/pq"
)

type MenuItem struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Price     int       `json:"price" db:"price"` // in rupiah
	Active    bool      `json:"active" db:"active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Company struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Address   string    `json:"address" db:"address"`
	Contact   string    `json:"contact" db:"contact"`
	Active    bool      `json:"active" db:"active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Employee struct {
	ID           int       `json:"id" db:"id"`
	CompanyID    int       `json:"company_id" db:"company_id"`
	Name         string    `json:"name" db:"name"`
	Email        string    `json:"email" db:"email"`
	WaContact    string    `json:"wa_contact" db:"wa_contact"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Active       bool      `json:"active" db:"active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type PasswordResetToken struct {
	ID         int       `json:"id" db:"id"`
	EmployeeID int       `json:"employee_id" db:"employee_id"`
	Token      string    `json:"token" db:"token"`
	ExpiresAt  time.Time `json:"expires_at" db:"expires_at"`
	Used       bool      `json:"used" db:"used"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type DailyMenu struct {
	ID                int           `json:"id" db:"id"`
	Date              time.Time     `json:"date" db:"date"`
	MenuItemIDs       pq.Int64Array `json:"menu_item_ids" db:"menu_item_ids"`
	NutritionistReset bool          `json:"nutritionist_reset" db:"nutritionist_reset"`
	CreatedAt         time.Time     `json:"created_at" db:"created_at"`
}

type OrderSession struct {
	ID        int        `json:"id" db:"id"`
	CompanyID int        `json:"company_id" db:"company_id"`
	Date      time.Time  `json:"date" db:"date"`
	Status    string     `json:"status" db:"status"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	ClosedAt  *time.Time `json:"closed_at" db:"closed_at"`
}

type IndividualOrder struct {
	ID          int           `json:"id" db:"id"`
	SessionID   int           `json:"session_id" db:"session_id"`
	EmployeeID  int           `json:"employee_id" db:"employee_id"`
	MenuItemIDs pq.Int64Array `json:"menu_item_ids" db:"menu_item_ids"`
	TotalPrice  int           `json:"total_price" db:"total_price"`
	Paid        bool          `json:"paid" db:"paid"`
	Status      string        `json:"status" db:"status"`
	CreatedAt   time.Time     `json:"created_at" db:"created_at"`
}

type NutritionistSelection struct {
	ID                 int           `json:"id" db:"id"`
	Date               time.Time     `json:"date" db:"date"`
	MenuItemIDs        pq.Int64Array `json:"menu_item_ids" db:"menu_item_ids"`
	SelectedIndices    pq.Int32Array `json:"selected_indices" db:"selected_indices"`
	Reasoning          string        `json:"reasoning" db:"reasoning"`
	NutritionalSummary string        `json:"nutritional_summary" db:"nutritional_summary"`
	CreatedAt          time.Time     `json:"created_at" db:"created_at"`
}

type NutritionistUserSelection struct {
	ID         int       `json:"id" db:"id"`
	EmployeeID int       `json:"employee_id" db:"employee_id"`
	Date       time.Time `json:"date" db:"date"`
	OrderID    *int      `json:"order_id" db:"order_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type StockEmptyItem struct {
	ID         int       `json:"id" db:"id"`
	MenuItemID int       `json:"menu_item_id" db:"menu_item_id"`
	Date       time.Time `json:"date" db:"date"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type UserStockEmptyNotification struct {
	ID                int       `json:"id" db:"id"`
	IndividualOrderID int       `json:"individual_order_id" db:"individual_order_id"`
	MenuItemID        int       `json:"menu_item_id" db:"menu_item_id"`
	Date              time.Time `json:"date" db:"date"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}

type UserNotification struct {
	ID               int       `json:"id" db:"id"`
	EmployeeID       int       `json:"employee_id" db:"employee_id"`
	NotificationType string    `json:"notification_type" db:"notification_type"`
	Title            string    `json:"title" db:"title"`
	Message          string    `json:"message" db:"message"`
	RedirectURL      *string   `json:"redirect_url" db:"redirect_url"`
	IsRead           bool      `json:"is_read" db:"is_read"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

const (
	StatusOpen           = "OPEN"
	StatusClosedOrders   = "CLOSED_FOR_ORDERS"
	StatusDelivered      = "DELIVERED"
	StatusPaymentPending = "PAYMENT_PENDING"
	StatusCompleted      = "COMPLETED"
)

// Individual Order Status Constants
const (
	OrderStatusPending        = "PENDING"
	OrderStatusReadyDelivery  = "READY_FOR_DELIVERY"
)

const (
	NotificationStockEmpty  = "STOCK_EMPTY"
	NotificationPaid        = "PAID"
	NotificationSessionClosed = "SESSION_CLOSED"
	NotificationMenuUpdated = "MENU_UPDATED"
)
