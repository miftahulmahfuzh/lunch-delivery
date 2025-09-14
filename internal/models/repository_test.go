package models

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testDate returns a consistent test date to avoid importing testutils
func testDate() time.Time {
	return time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
}

func setupMockDB(t *testing.T) (*Repository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewRepository(sqlxDB)

	cleanup := func() {
		_ = db.Close()
	}

	return repo, mock, cleanup
}

func TestNewRepository(t *testing.T) {
	t.Run("creates repository with valid database", func(t *testing.T) {
		db, _, err := sqlmock.New()
		require.NoError(t, err)
		defer func() { _ = db.Close() }()

		sqlxDB := sqlx.NewDb(db, "postgres")
		repo := NewRepository(sqlxDB)

		assert.NotNil(t, repo)
		assert.Equal(t, sqlxDB, repo.db)
	})

	t.Run("handles nil database", func(t *testing.T) {
		repo := NewRepository(nil)
		assert.NotNil(t, repo)
		assert.Nil(t, repo.db)
	})
}

func TestRepository_MenuItems(t *testing.T) {
	t.Run("CreateMenuItem success", func(t *testing.T) {
		repo, mock, cleanup := setupMockDB(t)
		defer cleanup()

		mock.ExpectQuery(`INSERT INTO menu_items`).
			WithArgs("Test Item", 25000).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price", "active", "created_at"}).
				AddRow(1, "Test Item", 25000, true, time.Now()))

		item, err := repo.CreateMenuItem("Test Item", 25000)

		assert.NoError(t, err)
		assert.NotNil(t, item)
		assert.Equal(t, "Test Item", item.Name)
		assert.Equal(t, 25000, item.Price)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("CreateMenuItem database error", func(t *testing.T) {
		repo, mock, cleanup := setupMockDB(t)
		defer cleanup()

		mock.ExpectQuery(`INSERT INTO menu_items`).
			WithArgs("Test Item", 25000).
			WillReturnError(assert.AnError)

		item, err := repo.CreateMenuItem("Test Item", 25000)

		assert.Error(t, err)
		assert.NotNil(t, item) // Implementation returns empty struct, not nil
		assert.Empty(t, item.Name)
		assert.Equal(t, 0, item.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetAllMenuItems success", func(t *testing.T) {
		repo, mock, cleanup := setupMockDB(t)
		defer cleanup()

		rows := sqlmock.NewRows([]string{"id", "name", "price", "active", "created_at"}).
			AddRow(1, "Item 1", 20000, true, time.Now()).
			AddRow(2, "Item 2", 30000, true, time.Now())

		mock.ExpectQuery(`SELECT \* FROM menu_items WHERE active = true ORDER BY name`).
			WillReturnRows(rows)

		items, err := repo.GetAllMenuItems()

		assert.NoError(t, err)
		assert.Len(t, items, 2)
		assert.Equal(t, "Item 1", items[0].Name)
		assert.Equal(t, "Item 2", items[1].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetAllMenuItems empty result", func(t *testing.T) {
		repo, mock, cleanup := setupMockDB(t)
		defer cleanup()

		mock.ExpectQuery(`SELECT \* FROM menu_items WHERE active = true ORDER BY name`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price", "active", "created_at"}))

		items, err := repo.GetAllMenuItems()

		assert.NoError(t, err)
		assert.Len(t, items, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("UpdateMenuItem success", func(t *testing.T) {
		repo, mock, cleanup := setupMockDB(t)
		defer cleanup()

		mock.ExpectExec(`UPDATE menu_items SET name = \$1, price = \$2 WHERE id = \$3`).
			WithArgs("Updated Item", 35000, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateMenuItem(1, "Updated Item", 35000)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DeleteMenuItem success", func(t *testing.T) {
		repo, mock, cleanup := setupMockDB(t)
		defer cleanup()

		mock.ExpectExec(`UPDATE menu_items SET active = false WHERE id = \$1`).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DeleteMenuItem(1)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetMenuItemsByIDs success", func(t *testing.T) {
		repo, mock, cleanup := setupMockDB(t)
		defer cleanup()

		rows := sqlmock.NewRows([]string{"id", "name", "price", "active", "created_at"}).
			AddRow(1, "Item 1", 20000, true, time.Now()).
			AddRow(3, "Item 3", 40000, true, time.Now())

		mock.ExpectQuery(`SELECT \* FROM menu_items WHERE id IN`).
			WithArgs(int64(1), int64(3)).
			WillReturnRows(rows)

		items, err := repo.GetMenuItemsByIDs([]int64{1, 3})

		assert.NoError(t, err)
		assert.Len(t, items, 2)
		assert.Equal(t, 1, items[0].ID)
		assert.Equal(t, 3, items[1].ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRepository_Companies(t *testing.T) {
	t.Run("CreateCompany success", func(t *testing.T) {
		repo, mock, cleanup := setupMockDB(t)
		defer cleanup()

		mock.ExpectQuery(`INSERT INTO companies`).
			WithArgs("Test Company", "Test Address", "test@company.com").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "address", "contact", "active", "created_at"}).
				AddRow(1, "Test Company", "Test Address", "test@company.com", true, time.Now()))

		company, err := repo.CreateCompany("Test Company", "Test Address", "test@company.com")

		assert.NoError(t, err)
		assert.NotNil(t, company)
		assert.Equal(t, "Test Company", company.Name)
		assert.Equal(t, "Test Address", company.Address)
		assert.Equal(t, "test@company.com", company.Contact)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetAllCompanies success", func(t *testing.T) {
		repo, mock, cleanup := setupMockDB(t)
		defer cleanup()

		rows := sqlmock.NewRows([]string{"id", "name", "address", "contact", "active", "created_at"}).
			AddRow(1, "Company 1", "Address 1", "contact1@test.com", true, time.Now()).
			AddRow(2, "Company 2", "Address 2", "contact2@test.com", true, time.Now())

		mock.ExpectQuery(`SELECT \* FROM companies WHERE active = true ORDER BY name`).
			WillReturnRows(rows)

		companies, err := repo.GetAllCompanies()

		assert.NoError(t, err)
		assert.Len(t, companies, 2)
		assert.Equal(t, "Company 1", companies[0].Name)
		assert.Equal(t, "Company 2", companies[1].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetCompanyByID success", func(t *testing.T) {
		repo, mock, cleanup := setupMockDB(t)
		defer cleanup()

		mock.ExpectQuery(`SELECT \* FROM companies WHERE id = \$1 AND active = true`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "address", "contact", "active", "created_at"}).
				AddRow(1, "Test Company", "Test Address", "test@company.com", true, time.Now()))

		company, err := repo.GetCompanyByID(1)

		assert.NoError(t, err)
		assert.NotNil(t, company)
		assert.Equal(t, 1, company.ID)
		assert.Equal(t, "Test Company", company.Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetCompanyByID not found", func(t *testing.T) {
		repo, mock, cleanup := setupMockDB(t)
		defer cleanup()

		mock.ExpectQuery(`SELECT \* FROM companies WHERE id = \$1 AND active = true`).
			WithArgs(999).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "address", "contact", "active", "created_at"}))

		company, err := repo.GetCompanyByID(999)

		assert.NoError(t, err)
		assert.Nil(t, company)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRepository_Employees(t *testing.T) {
	t.Run("CreateEmployee success", func(t *testing.T) {
		repo, mock, cleanup := setupMockDB(t)
		defer cleanup()

		mock.ExpectQuery(`INSERT INTO employees`).
			WithArgs(1, "John Doe", "john@test.com", "+628123456789", "hashed_password").
			WillReturnRows(sqlmock.NewRows([]string{"id", "company_id", "name", "email", "wa_contact", "password_hash", "active", "created_at"}).
				AddRow(1, 1, "John Doe", "john@test.com", "+628123456789", "hashed_password", true, time.Now()))

		employee, err := repo.CreateEmployee(1, "John Doe", "john@test.com", "+628123456789", "hashed_password")

		assert.NoError(t, err)
		assert.NotNil(t, employee)
		assert.Equal(t, "John Doe", employee.Name)
		assert.Equal(t, "john@test.com", employee.Email)
		assert.Equal(t, 1, employee.CompanyID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetEmployeeByEmail success", func(t *testing.T) {
		repo, mock, cleanup := setupMockDB(t)
		defer cleanup()

		mock.ExpectQuery(`SELECT \* FROM employees WHERE email = \$1 AND active = true`).
			WithArgs("john@test.com").
			WillReturnRows(sqlmock.NewRows([]string{"id", "company_id", "name", "email", "wa_contact", "password_hash", "active", "created_at"}).
				AddRow(1, 1, "John Doe", "john@test.com", "+628123456789", "hashed_password", true, time.Now()))

		employee, err := repo.GetEmployeeByEmail("john@test.com")

		assert.NoError(t, err)
		assert.NotNil(t, employee)
		assert.Equal(t, "john@test.com", employee.Email)
		assert.Equal(t, "John Doe", employee.Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetEmployeeByEmail not found", func(t *testing.T) {
		repo, mock, cleanup := setupMockDB(t)
		defer cleanup()

		mock.ExpectQuery(`SELECT \* FROM employees WHERE email = \$1 AND active = true`).
			WithArgs("nonexistent@test.com").
			WillReturnRows(sqlmock.NewRows([]string{"id", "company_id", "name", "email", "wa_contact", "password_hash", "active", "created_at"}))

		employee, err := repo.GetEmployeeByEmail("nonexistent@test.com")

		assert.NoError(t, err)
		assert.Nil(t, employee)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("UpdateEmployeePassword success", func(t *testing.T) {
		repo, mock, cleanup := setupMockDB(t)
		defer cleanup()

		mock.ExpectExec(`UPDATE employees SET password_hash = \$1 WHERE id = \$2`).
			WithArgs("new_hashed_password", 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateEmployeePassword(1, "new_hashed_password")

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRepository_OrderSessions(t *testing.T) {
	t.Run("CreateOrderSession success", func(t *testing.T) {
		repo, mock, cleanup := setupMockDB(t)
		defer cleanup()

		testDate := testDate()
		mock.ExpectQuery(`INSERT INTO order_sessions`).
			WithArgs(1, testDate.Format("2006-01-02"), StatusOpen).
			WillReturnRows(sqlmock.NewRows([]string{"id", "company_id", "date", "status", "created_at"}).
				AddRow(1, 1, testDate, StatusOpen, time.Now()))

		session, err := repo.CreateOrderSession(1, testDate, StatusOpen)

		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, 1, session.CompanyID)
		assert.Equal(t, StatusOpen, session.Status)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetOrderSession success", func(t *testing.T) {
		repo, mock, cleanup := setupMockDB(t)
		defer cleanup()

		testDate := testDate()
		mock.ExpectQuery(`SELECT \* FROM order_sessions WHERE company_id = \$1 AND date = \$2`).
			WithArgs(1, testDate.Format("2006-01-02")).
			WillReturnRows(sqlmock.NewRows([]string{"id", "company_id", "date", "status", "created_at", "closed_at"}).
				AddRow(1, 1, testDate, StatusOpen, time.Now(), nil))

		session, err := repo.GetOrderSession(1, testDate)

		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, 1, session.ID)
		assert.Equal(t, 1, session.CompanyID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("CloseOrderSession success", func(t *testing.T) {
		repo, mock, cleanup := setupMockDB(t)
		defer cleanup()

		mock.ExpectExec(`UPDATE order_sessions SET status = \$1, closed_at = CURRENT_TIMESTAMP WHERE id = \$2`).
			WithArgs(StatusClosedOrders, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.CloseOrderSession(1)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRepository_IndividualOrders(t *testing.T) {
	t.Run("CreateIndividualOrder success", func(t *testing.T) {
		repo, mock, cleanup := setupMockDB(t)
		defer cleanup()

		menuItemIDs := []int64{1, 2, 3}
		mock.ExpectQuery(`INSERT INTO individual_orders`).
			WithArgs(1, 1, pq.Array(menuItemIDs), 75000, OrderStatusPending).
			WillReturnRows(sqlmock.NewRows([]string{"id", "session_id", "employee_id", "menu_item_ids", "total_price", "paid", "status", "created_at"}).
				AddRow(1, 1, 1, pq.Array(menuItemIDs), 75000, false, OrderStatusPending, time.Now()))

		order, err := repo.CreateIndividualOrder(1, 1, menuItemIDs, 75000)

		assert.NoError(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, 1, order.SessionID)
		assert.Equal(t, 1, order.EmployeeID)
		assert.Equal(t, 75000, order.TotalPrice)
		assert.False(t, order.Paid)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetOrdersBySession success", func(t *testing.T) {
		repo, mock, cleanup := setupMockDB(t)
		defer cleanup()

		rows := sqlmock.NewRows([]string{"id", "session_id", "employee_id", "menu_item_ids", "total_price", "paid", "status", "created_at"}).
			AddRow(1, 1, 1, pq.Array([]int64{1, 2}), 50000, false, OrderStatusPending, time.Now()).
			AddRow(2, 1, 2, pq.Array([]int64{2, 3}), 60000, true, OrderStatusPending, time.Now())

		mock.ExpectQuery(`SELECT \* FROM individual_orders WHERE session_id = \$1`).
			WithArgs(1).
			WillReturnRows(rows)

		orders, err := repo.GetOrdersBySession(1)

		assert.NoError(t, err)
		assert.Len(t, orders, 2)
		assert.Equal(t, 1, orders[0].EmployeeID)
		assert.Equal(t, 2, orders[1].EmployeeID)
		assert.False(t, orders[0].Paid)
		assert.True(t, orders[1].Paid)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("MarkOrderPaid success", func(t *testing.T) {
		repo, mock, cleanup := setupMockDB(t)
		defer cleanup()

		mock.ExpectExec(`UPDATE individual_orders SET paid = true WHERE id = \$1`).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.MarkOrderPaid(1)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("UpdateOrderStatus success", func(t *testing.T) {
		repo, mock, cleanup := setupMockDB(t)
		defer cleanup()

		mock.ExpectExec(`UPDATE individual_orders SET status = \$1 WHERE id = \$2`).
			WithArgs(OrderStatusReadyDelivery, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateOrderStatus(1, OrderStatusReadyDelivery)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRepository_DailyMenu(t *testing.T) {
	t.Run("CreateDailyMenu success", func(t *testing.T) {
		repo, mock, cleanup := setupMockDB(t)
		defer cleanup()

		testDate := testDate()
		menuItemIDs := []int64{1, 2, 3, 4, 5}

		mock.ExpectQuery(`INSERT INTO daily_menus`).
			WithArgs(testDate.Format("2006-01-02"), pq.Array(menuItemIDs), true).
			WillReturnRows(sqlmock.NewRows([]string{"id", "date", "menu_item_ids", "nutritionist_reset", "created_at"}).
				AddRow(1, testDate, pq.Array(menuItemIDs), true, time.Now()))

		menu, err := repo.CreateDailyMenu(testDate, menuItemIDs)

		assert.NoError(t, err)
		assert.NotNil(t, menu)
		assert.Equal(t, testDate.Format("2006-01-02"), menu.Date.Format("2006-01-02"))
		assert.True(t, menu.NutritionistReset)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetDailyMenuByDate success", func(t *testing.T) {
		repo, mock, cleanup := setupMockDB(t)
		defer cleanup()

		testDate := testDate()
		menuItemIDs := []int64{1, 2, 3}

		mock.ExpectQuery(`SELECT \* FROM daily_menus WHERE date = \$1`).
			WithArgs(testDate.Format("2006-01-02")).
			WillReturnRows(sqlmock.NewRows([]string{"id", "date", "menu_item_ids", "nutritionist_reset", "created_at"}).
				AddRow(1, testDate, pq.Array(menuItemIDs), false, time.Now()))

		menu, err := repo.GetDailyMenuByDate(testDate)

		assert.NoError(t, err)
		assert.NotNil(t, menu)
		assert.Equal(t, 1, menu.ID)
		assert.False(t, menu.NutritionistReset)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetDailyMenuByDate not found", func(t *testing.T) {
		repo, mock, cleanup := setupMockDB(t)
		defer cleanup()

		testDate := testDate()
		mock.ExpectQuery(`SELECT \* FROM daily_menus WHERE date = \$1`).
			WithArgs(testDate.Format("2006-01-02")).
			WillReturnRows(sqlmock.NewRows([]string{"id", "date", "menu_item_ids", "nutritionist_reset", "created_at"}))

		menu, err := repo.GetDailyMenuByDate(testDate)

		assert.NoError(t, err)
		assert.Nil(t, menu)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRepository_UserNotifications(t *testing.T) {
	t.Run("CreateUserNotification success", func(t *testing.T) {
		repo, mock, cleanup := setupMockDB(t)
		defer cleanup()

		redirectURL := "/test/path"
		mock.ExpectExec(`INSERT INTO user_notifications`).
			WithArgs(1, NotificationStockEmpty, "Test Message", &redirectURL).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.CreateUserNotification(1, NotificationStockEmpty, "Test Message", &redirectURL)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetUserNotifications success", func(t *testing.T) {
		repo, mock, cleanup := setupMockDB(t)
		defer cleanup()

		rows := sqlmock.NewRows([]string{"id", "employee_id", "notification_type", "title", "message", "redirect_url", "is_read", "created_at"}).
			AddRow(1, 1, NotificationStockEmpty, "Title 1", "Message 1", "/path1", false, time.Now()).
			AddRow(2, 1, NotificationPaid, "Title 2", "Message 2", nil, true, time.Now())

		mock.ExpectQuery(`SELECT \* FROM user_notifications WHERE employee_id = \$1 ORDER BY created_at DESC LIMIT 10`).
			WithArgs(1).
			WillReturnRows(rows)

		notifications, err := repo.GetUserNotifications(1, 10)

		assert.NoError(t, err)
		assert.Len(t, notifications, 2)
		assert.Equal(t, NotificationStockEmpty, notifications[0].NotificationType)
		assert.Equal(t, NotificationPaid, notifications[1].NotificationType)
		assert.False(t, notifications[0].IsRead)
		assert.True(t, notifications[1].IsRead)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("MarkNotificationRead success", func(t *testing.T) {
		repo, mock, cleanup := setupMockDB(t)
		defer cleanup()

		mock.ExpectExec(`UPDATE user_notifications SET is_read = true WHERE id = \$1`).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.MarkNotificationRead(1)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DeleteAllUserNotifications success", func(t *testing.T) {
		repo, mock, cleanup := setupMockDB(t)
		defer cleanup()

		mock.ExpectExec(`DELETE FROM user_notifications WHERE employee_id = \$1`).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 5))

		err := repo.DeleteAllUserNotifications(1)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}