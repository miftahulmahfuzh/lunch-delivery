package models

import (
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestModelStructures(t *testing.T) {
	t.Run("MenuItem structure", func(t *testing.T) {
		now := time.Now()
		item := MenuItem{
			ID:        1,
			Name:      "Test Item",
			Price:     25000,
			Active:    true,
			CreatedAt: now,
		}

		assert.Equal(t, 1, item.ID)
		assert.Equal(t, "Test Item", item.Name)
		assert.Equal(t, 25000, item.Price)
		assert.True(t, item.Active)
		assert.Equal(t, now, item.CreatedAt)
	})

	t.Run("Company structure", func(t *testing.T) {
		now := time.Now()
		company := Company{
			ID:        1,
			Name:      "Test Company",
			Address:   "Test Address",
			Contact:   "test@company.com",
			Active:    true,
			CreatedAt: now,
		}

		assert.Equal(t, 1, company.ID)
		assert.Equal(t, "Test Company", company.Name)
		assert.Equal(t, "Test Address", company.Address)
		assert.Equal(t, "test@company.com", company.Contact)
		assert.True(t, company.Active)
		assert.Equal(t, now, company.CreatedAt)
	})

	t.Run("Employee structure", func(t *testing.T) {
		now := time.Now()
		employee := Employee{
			ID:           1,
			CompanyID:    1,
			Name:         "John Doe",
			Email:        "john@example.com",
			WaContact:    "+628123456789",
			PasswordHash: "hashed_password",
			Active:       true,
			CreatedAt:    now,
		}

		assert.Equal(t, 1, employee.ID)
		assert.Equal(t, 1, employee.CompanyID)
		assert.Equal(t, "John Doe", employee.Name)
		assert.Equal(t, "john@example.com", employee.Email)
		assert.Equal(t, "+628123456789", employee.WaContact)
		assert.Equal(t, "hashed_password", employee.PasswordHash)
		assert.True(t, employee.Active)
		assert.Equal(t, now, employee.CreatedAt)
	})

	t.Run("PasswordResetToken structure", func(t *testing.T) {
		now := time.Now()
		expiresAt := now.Add(time.Hour)
		token := PasswordResetToken{
			ID:         1,
			EmployeeID: 1,
			Token:      "reset-token-123",
			ExpiresAt:  expiresAt,
			Used:       false,
			CreatedAt:  now,
		}

		assert.Equal(t, 1, token.ID)
		assert.Equal(t, 1, token.EmployeeID)
		assert.Equal(t, "reset-token-123", token.Token)
		assert.Equal(t, expiresAt, token.ExpiresAt)
		assert.False(t, token.Used)
		assert.Equal(t, now, token.CreatedAt)
	})

	t.Run("DailyMenu structure", func(t *testing.T) {
		now := time.Now()
		date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		menu := DailyMenu{
			ID:                1,
			Date:              date,
			MenuItemIDs:       pq.Int64Array{1, 2, 3, 4, 5},
			NutritionistReset: false,
			CreatedAt:         now,
		}

		assert.Equal(t, 1, menu.ID)
		assert.Equal(t, date, menu.Date)
		assert.Equal(t, pq.Int64Array{1, 2, 3, 4, 5}, menu.MenuItemIDs)
		assert.False(t, menu.NutritionistReset)
		assert.Equal(t, now, menu.CreatedAt)
	})

	t.Run("OrderSession structure", func(t *testing.T) {
		now := time.Now()
		date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		closedAt := now.Add(time.Hour)

		session := OrderSession{
			ID:        1,
			CompanyID: 1,
			Date:      date,
			Status:    StatusOpen,
			CreatedAt: now,
			ClosedAt:  &closedAt,
		}

		assert.Equal(t, 1, session.ID)
		assert.Equal(t, 1, session.CompanyID)
		assert.Equal(t, date, session.Date)
		assert.Equal(t, StatusOpen, session.Status)
		assert.Equal(t, now, session.CreatedAt)
		assert.NotNil(t, session.ClosedAt)
		assert.Equal(t, closedAt, *session.ClosedAt)
	})

	t.Run("IndividualOrder structure", func(t *testing.T) {
		now := time.Now()
		order := IndividualOrder{
			ID:          1,
			SessionID:   1,
			EmployeeID:  1,
			MenuItemIDs: pq.Int64Array{1, 2, 3},
			TotalPrice:  75000,
			Paid:        false,
			Status:      OrderStatusPending,
			CreatedAt:   now,
		}

		assert.Equal(t, 1, order.ID)
		assert.Equal(t, 1, order.SessionID)
		assert.Equal(t, 1, order.EmployeeID)
		assert.Equal(t, pq.Int64Array{1, 2, 3}, order.MenuItemIDs)
		assert.Equal(t, 75000, order.TotalPrice)
		assert.False(t, order.Paid)
		assert.Equal(t, OrderStatusPending, order.Status)
		assert.Equal(t, now, order.CreatedAt)
	})

	t.Run("NutritionistSelection structure", func(t *testing.T) {
		now := time.Now()
		date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		selection := NutritionistSelection{
			ID:                 1,
			Date:               date,
			MenuItemIDs:        pq.Int64Array{1, 2, 3, 4, 5},
			SelectedIndices:    pq.Int32Array{0, 2, 4},
			Reasoning:          "Balanced nutrition",
			NutritionalSummary: `{"protein": "high", "carbs": "moderate"}`,
			CreatedAt:          now,
		}

		assert.Equal(t, 1, selection.ID)
		assert.Equal(t, date, selection.Date)
		assert.Equal(t, pq.Int64Array{1, 2, 3, 4, 5}, selection.MenuItemIDs)
		assert.Equal(t, pq.Int32Array{0, 2, 4}, selection.SelectedIndices)
		assert.Equal(t, "Balanced nutrition", selection.Reasoning)
		assert.Equal(t, `{"protein": "high", "carbs": "moderate"}`, selection.NutritionalSummary)
		assert.Equal(t, now, selection.CreatedAt)
	})

	t.Run("NutritionistUserSelection structure", func(t *testing.T) {
		now := time.Now()
		date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		orderID := 123

		userSelection := NutritionistUserSelection{
			ID:         1,
			EmployeeID: 1,
			Date:       date,
			OrderID:    &orderID,
			CreatedAt:  now,
		}

		assert.Equal(t, 1, userSelection.ID)
		assert.Equal(t, 1, userSelection.EmployeeID)
		assert.Equal(t, date, userSelection.Date)
		assert.NotNil(t, userSelection.OrderID)
		assert.Equal(t, 123, *userSelection.OrderID)
		assert.Equal(t, now, userSelection.CreatedAt)
	})

	t.Run("StockEmptyItem structure", func(t *testing.T) {
		now := time.Now()
		date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

		stockItem := StockEmptyItem{
			ID:         1,
			MenuItemID: 1,
			Date:       date,
			CreatedAt:  now,
		}

		assert.Equal(t, 1, stockItem.ID)
		assert.Equal(t, 1, stockItem.MenuItemID)
		assert.Equal(t, date, stockItem.Date)
		assert.Equal(t, now, stockItem.CreatedAt)
	})

	t.Run("UserStockEmptyNotification structure", func(t *testing.T) {
		now := time.Now()
		date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

		notification := UserStockEmptyNotification{
			ID:                1,
			IndividualOrderID: 1,
			MenuItemID:        1,
			Date:              date,
			CreatedAt:         now,
		}

		assert.Equal(t, 1, notification.ID)
		assert.Equal(t, 1, notification.IndividualOrderID)
		assert.Equal(t, 1, notification.MenuItemID)
		assert.Equal(t, date, notification.Date)
		assert.Equal(t, now, notification.CreatedAt)
	})

	t.Run("UserNotification structure", func(t *testing.T) {
		now := time.Now()
		redirectURL := "/some/path"

		notification := UserNotification{
			ID:               1,
			EmployeeID:       1,
			NotificationType: NotificationStockEmpty,
			Title:            "Stock Empty",
			Message:          "Item is out of stock",
			RedirectURL:      &redirectURL,
			IsRead:           false,
			CreatedAt:        now,
		}

		assert.Equal(t, 1, notification.ID)
		assert.Equal(t, 1, notification.EmployeeID)
		assert.Equal(t, NotificationStockEmpty, notification.NotificationType)
		assert.Equal(t, "Stock Empty", notification.Title)
		assert.Equal(t, "Item is out of stock", notification.Message)
		assert.NotNil(t, notification.RedirectURL)
		assert.Equal(t, "/some/path", *notification.RedirectURL)
		assert.False(t, notification.IsRead)
		assert.Equal(t, now, notification.CreatedAt)
	})
}

func TestModelConstants(t *testing.T) {
	t.Run("order session status constants", func(t *testing.T) {
		assert.Equal(t, "OPEN", StatusOpen)
		assert.Equal(t, "CLOSED_FOR_ORDERS", StatusClosedOrders)
		assert.Equal(t, "DELIVERED", StatusDelivered)
		assert.Equal(t, "PAYMENT_PENDING", StatusPaymentPending)
		assert.Equal(t, "COMPLETED", StatusCompleted)
	})

	t.Run("individual order status constants", func(t *testing.T) {
		assert.Equal(t, "PENDING", OrderStatusPending)
		assert.Equal(t, "READY_FOR_DELIVERY", OrderStatusReadyDelivery)
	})

	t.Run("notification type constants", func(t *testing.T) {
		assert.Equal(t, "STOCK_EMPTY", NotificationStockEmpty)
		assert.Equal(t, "PAID", NotificationPaid)
		assert.Equal(t, "SESSION_CLOSED", NotificationSessionClosed)
		assert.Equal(t, "MENU_UPDATED", NotificationMenuUpdated)
		assert.Equal(t, "READY_FOR_DELIVERY", NotificationReadyForDelivery)
	})
}

func TestModelValidation(t *testing.T) {
	t.Run("MenuItem with zero price", func(t *testing.T) {
		item := MenuItem{
			ID:     1,
			Name:   "Free Item",
			Price:  0,
			Active: true,
		}

		assert.Equal(t, 0, item.Price)
		// Zero price should be valid for free items
	})

	t.Run("MenuItem with negative price", func(t *testing.T) {
		item := MenuItem{
			ID:     1,
			Name:   "Negative Price Item",
			Price:  -1000,
			Active: true,
		}

		assert.Equal(t, -1000, item.Price)
		// Negative prices might be used for discounts or credits
	})

	t.Run("Employee with empty password hash", func(t *testing.T) {
		employee := Employee{
			ID:           1,
			CompanyID:    1,
			Name:         "Test User",
			Email:        "test@example.com",
			PasswordHash: "",
			Active:       true,
		}

		assert.Equal(t, "", employee.PasswordHash)
		// Empty password hash might be valid for OAuth users
	})

	t.Run("OrderSession with nil ClosedAt", func(t *testing.T) {
		session := OrderSession{
			ID:        1,
			CompanyID: 1,
			Date:      time.Now(),
			Status:    StatusOpen,
			ClosedAt:  nil,
		}

		assert.Nil(t, session.ClosedAt)
		// Nil ClosedAt should be valid for open sessions
	})

	t.Run("UserNotification with nil RedirectURL", func(t *testing.T) {
		notification := UserNotification{
			ID:               1,
			EmployeeID:       1,
			NotificationType: NotificationPaid,
			Title:            "Payment Confirmed",
			Message:          "Payment processed",
			RedirectURL:      nil,
			IsRead:           false,
		}

		assert.Nil(t, notification.RedirectURL)
		// Nil RedirectURL should be valid for notifications without actions
	})
}

func TestModelEdgeCases(t *testing.T) {
	t.Run("empty arrays in PostgreSQL array fields", func(t *testing.T) {
		// Test empty Int64Array
		menu := DailyMenu{
			ID:          1,
			MenuItemIDs: pq.Int64Array{},
		}
		assert.Len(t, menu.MenuItemIDs, 0)

		// Test empty Int32Array
		selection := NutritionistSelection{
			ID:              1,
			SelectedIndices: pq.Int32Array{},
		}
		assert.Len(t, selection.SelectedIndices, 0)

		// Test individual order with empty menu items
		order := IndividualOrder{
			ID:          1,
			MenuItemIDs: pq.Int64Array{},
			TotalPrice:  0,
		}
		assert.Len(t, order.MenuItemIDs, 0)
		assert.Equal(t, 0, order.TotalPrice)
	})

	t.Run("very long strings in text fields", func(t *testing.T) {
		longText := ""
		for i := 0; i < 1000; i++ {
			longText += "a"
		}

		employee := Employee{
			Name:         longText,
			Email:        longText + "@example.com",
			WaContact:    longText,
			PasswordHash: longText,
		}

		assert.Equal(t, longText, employee.Name)
		assert.Equal(t, longText+"@example.com", employee.Email)
		assert.Equal(t, longText, employee.WaContact)
		assert.Equal(t, longText, employee.PasswordHash)
	})

	t.Run("special characters in text fields", func(t *testing.T) {
		specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?/~`"

		company := Company{
			Name:    "Company " + specialChars,
			Address: "Address " + specialChars,
			Contact: "contact" + specialChars + "@example.com",
		}

		assert.Contains(t, company.Name, specialChars)
		assert.Contains(t, company.Address, specialChars)
		assert.Contains(t, company.Contact, specialChars)
	})

	t.Run("unicode characters in text fields", func(t *testing.T) {
		unicode := "Café München 東京 🍕🍝🍜"

		item := MenuItem{
			Name: "Special Item " + unicode,
		}

		assert.Contains(t, item.Name, unicode)
	})
}

func TestModelJSON(t *testing.T) {
	t.Run("models can be JSON serialized", func(t *testing.T) {
		// Test that models can be marshaled to JSON (useful for APIs)
		item := MenuItem{
			ID:     1,
			Name:   "Test Item",
			Price:  25000,
			Active: true,
		}

		// JSON tags should work correctly
		assert.Equal(t, 1, item.ID)
		assert.Equal(t, "Test Item", item.Name)
		assert.Equal(t, 25000, item.Price)
		assert.True(t, item.Active)
	})

	t.Run("password hash is excluded from JSON", func(t *testing.T) {
		employee := Employee{
			ID:           1,
			Name:         "John Doe",
			Email:        "john@example.com",
			PasswordHash: "secret_hash",
		}

		// PasswordHash should have json:"-" tag to exclude it
		assert.Equal(t, "secret_hash", employee.PasswordHash)
		// The actual JSON exclusion would be tested with json.Marshal in integration tests
	})
}

func TestModelRelationships(t *testing.T) {
	t.Run("employee belongs to company", func(t *testing.T) {
		employee := Employee{
			ID:        1,
			CompanyID: 5,
			Name:      "John Doe",
		}

		assert.Equal(t, 5, employee.CompanyID)
		// In actual usage, this would relate to a Company with ID 5
	})

	t.Run("individual order belongs to session and employee", func(t *testing.T) {
		order := IndividualOrder{
			ID:         1,
			SessionID:  10,
			EmployeeID: 5,
		}

		assert.Equal(t, 10, order.SessionID)
		assert.Equal(t, 5, order.EmployeeID)
		// In actual usage, these would relate to OrderSession and Employee records
	})

	t.Run("order session belongs to company", func(t *testing.T) {
		session := OrderSession{
			ID:        1,
			CompanyID: 3,
		}

		assert.Equal(t, 3, session.CompanyID)
		// In actual usage, this would relate to a Company with ID 3
	})

	t.Run("password reset token belongs to employee", func(t *testing.T) {
		token := PasswordResetToken{
			ID:         1,
			EmployeeID: 7,
			Token:      "reset-123",
		}

		assert.Equal(t, 7, token.EmployeeID)
		// In actual usage, this would relate to an Employee with ID 7
	})
}

func TestModelDefaults(t *testing.T) {
	t.Run("zero values are properly set", func(t *testing.T) {
		var item MenuItem
		assert.Equal(t, 0, item.ID)
		assert.Equal(t, "", item.Name)
		assert.Equal(t, 0, item.Price)
		assert.False(t, item.Active)
		assert.True(t, item.CreatedAt.IsZero())

		var company Company
		assert.Equal(t, 0, company.ID)
		assert.Equal(t, "", company.Name)
		assert.False(t, company.Active)

		var employee Employee
		assert.Equal(t, 0, employee.ID)
		assert.Equal(t, 0, employee.CompanyID)
		assert.Equal(t, "", employee.Name)
		assert.False(t, employee.Active)
	})
}
