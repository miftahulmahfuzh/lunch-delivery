// internal/models/models.go
package models

import (
	"time"

	"github.com/lib/pq"
)

type MenuItem struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Price     int       `json:"price" db:"price"` // in cents
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

type DailyMenu struct {
	ID          int           `json:"id" db:"id"`
	Date        time.Time     `json:"date" db:"date"`
	MenuItemIDs pq.Int64Array `json:"menu_item_ids" db:"menu_item_ids"`
	CreatedAt   time.Time     `json:"created_at" db:"created_at"`
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
	CreatedAt   time.Time     `json:"created_at" db:"created_at"`
}

const (
	StatusOpen           = "OPEN"
	StatusClosedOrders   = "CLOSED_FOR_ORDERS"
	StatusDelivered      = "DELIVERED"
	StatusPaymentPending = "PAYMENT_PENDING"
	StatusCompleted      = "COMPLETED"
)
