// internal/models/repository.go
package models

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// Menu Items
func (r *Repository) CreateMenuItem(name string, price int) (*MenuItem, error) {
	var item MenuItem
	err := r.db.Get(&item,
		`INSERT INTO menu_items (name, price) VALUES ($1, $2) RETURNING *`,
		name, price)
	return &item, err
}

func (r *Repository) GetAllMenuItems() ([]MenuItem, error) {
	var items []MenuItem
	err := r.db.Select(&items, `SELECT * FROM menu_items WHERE active = true ORDER BY name`)
	return items, err
}

func (r *Repository) UpdateMenuItem(id int, name string, price int) error {
	_, err := r.db.Exec(`UPDATE menu_items SET name = $1, price = $2 WHERE id = $3`,
		name, price, id)
	return err
}

func (r *Repository) DeleteMenuItem(id int) error {
	_, err := r.db.Exec(`UPDATE menu_items SET active = false WHERE id = $1`, id)
	return err
}

// Companies
func (r *Repository) CreateCompany(name, address, contact string) (*Company, error) {
	var company Company
	err := r.db.Get(&company,
		`INSERT INTO companies (name, address, contact) VALUES ($1, $2, $3) RETURNING *`,
		name, address, contact)
	return &company, err
}

func (r *Repository) GetAllCompanies() ([]Company, error) {
	var companies []Company
	err := r.db.Select(&companies, `SELECT * FROM companies WHERE active = true ORDER BY name`)
	return companies, err
}

func (r *Repository) GetCompanyByID(id int) (*Company, error) {
	var company Company
	err := r.db.Get(&company, `SELECT * FROM companies WHERE id = $1 AND active = true`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &company, err
}

// Employees
func (r *Repository) CreateEmployee(companyID int, name, email, waContact, passwordHash string) (*Employee, error) {
	var employee Employee
	err := r.db.Get(&employee,
		`INSERT INTO employees (company_id, name, email, wa_contact, password_hash) VALUES ($1, $2, $3, $4, $5) RETURNING *`,
		companyID, name, email, waContact, passwordHash)
	return &employee, err
}

func (r *Repository) GetEmployeeByEmail(email string) (*Employee, error) {
	var employee Employee
	err := r.db.Get(&employee, `SELECT * FROM employees WHERE email = $1 AND active = true`, email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &employee, err
}

func (r *Repository) GetEmployeesByCompany(companyID int) ([]Employee, error) {
	var employees []Employee
	err := r.db.Select(&employees,
		`SELECT * FROM employees WHERE company_id = $1 AND active = true ORDER BY name`,
		companyID)
	return employees, err
}

func (r *Repository) UpdateEmployeePassword(employeeID int, passwordHash string) error {
	_, err := r.db.Exec(`UPDATE employees SET password_hash = $1 WHERE id = $2`, passwordHash, employeeID)
	return err
}

// Password Reset Tokens
func (r *Repository) CreatePasswordResetToken(employeeID int, token string, expiresAt time.Time) error {
	_, err := r.db.Exec(
		`INSERT INTO password_reset_tokens (employee_id, token, expires_at) VALUES ($1, $2, $3)`,
		employeeID, token, expiresAt)
	return err
}

func (r *Repository) GetPasswordResetToken(token string) (*PasswordResetToken, error) {
	var resetToken PasswordResetToken
	err := r.db.Get(&resetToken,
		`SELECT * FROM password_reset_tokens WHERE token = $1 AND used = false AND expires_at > NOW()`,
		token)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &resetToken, err
}

func (r *Repository) MarkPasswordResetTokenAsUsed(token string) error {
	_, err := r.db.Exec(`UPDATE password_reset_tokens SET used = true WHERE token = $1`, token)
	return err
}

func (r *Repository) DeletePasswordResetToken(token string) error {
	_, err := r.db.Exec(`DELETE FROM password_reset_tokens WHERE token = $1`, token)
	return err
}

func (r *Repository) CleanupExpiredPasswordResetTokens() error {
	_, err := r.db.Exec(`DELETE FROM password_reset_tokens WHERE expires_at < NOW() OR used = true`)
	return err
}

// Daily Menu
func (r *Repository) CreateDailyMenu(date time.Time, menuItemIDs []int64) (*DailyMenu, error) {
	var menu DailyMenu
	err := r.db.Get(&menu,
		`INSERT INTO daily_menus (date, menu_item_ids, nutritionist_reset) VALUES ($1, $2, $3)
         ON CONFLICT (date) DO UPDATE SET 
             menu_item_ids = $2,
             nutritionist_reset = $3
         RETURNING *`,
		date.Format("2006-01-02"), pq.Array(menuItemIDs), true) // Always set reset flag to true when menu is updated
	return &menu, err
}

func (r *Repository) GetDailyMenuByDate(date time.Time) (*DailyMenu, error) {
	var menu DailyMenu
	err := r.db.Get(&menu, `SELECT * FROM daily_menus WHERE date = $1`, date.Format("2006-01-02"))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &menu, err
}

// Order Sessions
func (r *Repository) CreateOrderSession(companyID int, date time.Time, status string) (*OrderSession, error) {
	var session OrderSession
	err := r.db.Get(&session,
		`INSERT INTO order_sessions (company_id, date, status) VALUES ($1, $2, $3) RETURNING *`,
		companyID, date.Format("2006-01-02"), status)
	return &session, err
}

func (r *Repository) GetOrderSessionsByDate(date time.Time) ([]OrderSession, error) {
	var sessions []OrderSession
	err := r.db.Select(&sessions, `SELECT * FROM order_sessions WHERE date = $1`, date.Format("2006-01-02"))
	return sessions, err
}

func (r *Repository) GetOrderSession(companyID int, date time.Time) (*OrderSession, error) {
	var session OrderSession
	err := r.db.Get(&session,
		`SELECT * FROM order_sessions WHERE company_id = $1 AND date = $2`,
		companyID, date.Format("2006-01-02"))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &session, err
}

func (r *Repository) UpdateOrderSessionStatus(id int, status string) error {
	_, err := r.db.Exec(`UPDATE order_sessions SET status = $1 WHERE id = $2`, status, id)
	return err
}

func (r *Repository) CloseOrderSession(id int) error {
	_, err := r.db.Exec(`UPDATE order_sessions SET status = $1, closed_at = CURRENT_TIMESTAMP WHERE id = $2`,
		StatusClosedOrders, id)
	return err
}

// Individual Orders
func (r *Repository) CreateIndividualOrder(sessionID, employeeID int, menuItemIDs []int64, totalPrice int) (*IndividualOrder, error) {
	var order IndividualOrder
	err := r.db.Get(&order,
		`INSERT INTO individual_orders (session_id, employee_id, menu_item_ids, total_price, status)
         VALUES ($1, $2, $3, $4, $5)
         ON CONFLICT (session_id, employee_id)
         DO UPDATE SET menu_item_ids = $3, total_price = $4, status = $5
         RETURNING *`,
		sessionID, employeeID, pq.Array(menuItemIDs), totalPrice, OrderStatusPending)
	return &order, err
}

func (r *Repository) GetOrdersBySession(sessionID int) ([]IndividualOrder, error) {
	var orders []IndividualOrder
	err := r.db.Select(&orders, `SELECT * FROM individual_orders WHERE session_id = $1`, sessionID)
	return orders, err
}

func (r *Repository) MarkOrderPaid(id int) error {
	_, err := r.db.Exec(`UPDATE individual_orders SET paid = true WHERE id = $1`, id)
	return err
}

func (r *Repository) MarkOrderUnpaid(id int) error {
	_, err := r.db.Exec(`UPDATE individual_orders SET paid = false WHERE id = $1`, id)
	return err
}

func (r *Repository) UpdateOrderStatus(id int, status string) error {
	_, err := r.db.Exec(`UPDATE individual_orders SET status = $1 WHERE id = $2`, status, id)
	return err
}

func (r *Repository) GetMenuItemsByIDs(ids []int64) ([]MenuItem, error) {
	var items []MenuItem
	query, args, err := sqlx.In(`SELECT * FROM menu_items WHERE id IN (?)`, ids)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)
	err = r.db.Select(&items, query, args...)
	return items, err
}

func (r *Repository) GetOrderSessionsByDateWithCompany(date time.Time) ([]OrderSessionWithCompany, error) {
	var sessions []OrderSessionWithCompany
	err := r.db.Select(&sessions, `
        SELECT os.*, c.name as company_name
        FROM order_sessions os
        JOIN companies c ON os.company_id = c.id
        WHERE os.date = $1
        ORDER BY c.name`,
		date.Format("2006-01-02"))
	return sessions, err
}

func (r *Repository) ReopenOrderSession(id int) error {
	_, err := r.db.Exec(`UPDATE order_sessions SET status = $1, closed_at = NULL WHERE id = $2`,
		StatusOpen, id)
	return err
}

// For individual orders

func (r *Repository) GetOrdersBySessionWithDetails(sessionID int) ([]IndividualOrderWithDetails, error) {
	var orders []IndividualOrderWithDetails
	err := r.db.Select(&orders, `
        SELECT io.*, e.name as employee_name
        FROM individual_orders io
        JOIN employees e ON io.employee_id = e.id
        WHERE io.session_id = $1
        ORDER BY e.name`,
		sessionID)
	if err != nil {
		return nil, err
	}

	// Get menu item names for each order
	for i, order := range orders {
		var itemIDs []int64
		for _, id := range order.MenuItemIDs {
			itemIDs = append(itemIDs, id)
		}

		if len(itemIDs) > 0 {
			items, err := r.GetMenuItemsByIDs(itemIDs)
			if err != nil {
				continue
			}

			var names []string
			for _, item := range items {
				names = append(names, item.Name)
			}
			orders[i].MenuItemNames = strings.Join(names, ", ")
		}
	}

	return orders, nil
}

func (r *Repository) GetOrderSessionWithCompany(id int) (*OrderSessionWithCompany, error) {
	var session OrderSessionWithCompany
	err := r.db.Get(&session, `
        SELECT os.*, c.name as company_name
        FROM order_sessions os
        JOIN companies c ON os.company_id = c.id
        WHERE os.id = $1`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &session, err
}

// Companies Page
func (r *Repository) UpdateCompany(id int, name, address, contact string) error {
	_, err := r.db.Exec(`UPDATE companies SET name = $1, address = $2, contact = $3 WHERE id = $4`,
		name, address, contact, id)
	return err
}

func (r *Repository) DeleteCompany(id int) error {
	_, err := r.db.Exec(`UPDATE companies SET active = false WHERE id = $1`, id)
	return err
}

// Add to internal/models/repository.go
func (r *Repository) UpdateEmployee(id int, name, email, waContact string) error {
	_, err := r.db.Exec(`UPDATE employees SET name = $1, email = $2, wa_contact = $3 WHERE id = $4`,
		name, email, waContact, id)
	return err
}

func (r *Repository) DeleteEmployee(id int) error {
	_, err := r.db.Exec(`UPDATE employees SET active = false WHERE id = $1`, id)
	return err
}

// Add to internal/models/repository.go
func (r *Repository) GetEmployeeByID(id int) (*Employee, error) {
	var employee Employee
	err := r.db.Get(&employee, `SELECT * FROM employees WHERE id = $1 AND active = true`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &employee, err
}

func (r *Repository) GetRecentOrdersByEmployee(employeeID int, startDate, endDate time.Time) ([]RecentOrder, error) {
	type OrderRow struct {
		Date        time.Time     `db:"date"`
		TotalPrice  int           `db:"total_price"`
		Paid        bool          `db:"paid"`
		MenuItemIDs pq.Int64Array `db:"menu_item_ids"`
	}

	var orderRows []OrderRow
	err := r.db.Select(&orderRows, `
        SELECT os.date, io.total_price, io.paid, io.menu_item_ids
        FROM individual_orders io
        JOIN order_sessions os ON io.session_id = os.id
        WHERE io.employee_id = $1 AND os.date >= $2 AND os.date <= $3
        ORDER BY os.date DESC`, employeeID, startDate, endDate)

	if err != nil {
		return nil, err
	}

	var orders []RecentOrder

	// Get menu item names for each order
	for _, row := range orderRows {
		var menuItemNames []string

		if len(row.MenuItemIDs) > 0 {
			// Convert pq.Int64Array to []int64
			var itemIDs []int64
			for _, id := range row.MenuItemIDs {
				itemIDs = append(itemIDs, id)
			}

			items, err := r.GetMenuItemsByIDs(itemIDs)
			if err == nil {
				for _, item := range items {
					menuItemNames = append(menuItemNames, item.Name)
				}
			}
		}

		order := RecentOrder{
			Date:          row.Date,
			TotalPrice:    row.TotalPrice,
			Paid:          row.Paid,
			MenuItemNames: strings.Join(menuItemNames, ", "),
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// Nutritionist Selections
func (r *Repository) GetNutritionistSelectionByDate(date time.Time) (*NutritionistSelection, error) {
	var selection NutritionistSelection
	err := r.db.Get(&selection, `SELECT * FROM nutritionist_selections WHERE date = $1`, date.Format("2006-01-02"))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &selection, err
}

func (r *Repository) CreateNutritionistSelection(date time.Time, menuItemIDs []int64, selectedIndices []int32, reasoning string, nutritionalSummary string) (*NutritionistSelection, error) {
	var selection NutritionistSelection
	err := r.db.Get(&selection,
		`INSERT INTO nutritionist_selections (date, menu_item_ids, selected_indices, reasoning, nutritional_summary)
         VALUES ($1, $2, $3, $4, $5) RETURNING *`,
		date.Format("2006-01-02"), pq.Array(menuItemIDs), pq.Array(selectedIndices), reasoning, nutritionalSummary)
	return &selection, err
}

func (r *Repository) DeleteNutritionistSelection(date time.Time) error {
	_, err := r.db.Exec(`DELETE FROM nutritionist_selections WHERE date = $1`, date.Format("2006-01-02"))
	return err
}

// Nutritionist User Selection Tracking
func (r *Repository) CreateNutritionistUserSelection(date time.Time, employeeID int, menuItemIDs []int64) (*NutritionistUserSelection, error) {
	var selection NutritionistUserSelection
	err := r.db.Get(&selection,
		`INSERT INTO nutritionist_user_selections (employee_id, date)
         VALUES ($1, $2)
         ON CONFLICT (employee_id, date)
         DO NOTHING RETURNING *`,
		employeeID, date.Format("2006-01-02"))
	return &selection, err
}

func (r *Repository) GetNutritionistUserSelectionByDate(employeeID int, date time.Time) (*NutritionistUserSelection, error) {
	var selection NutritionistUserSelection
	err := r.db.Get(&selection,
		`SELECT * FROM nutritionist_user_selections WHERE employee_id = $1 AND date = $2`,
		employeeID, date.Format("2006-01-02"))
	return &selection, err
}

func (r *Repository) GetNutritionistUsersByDateAndUnpaid(date time.Time) ([]NutritionistUserSelection, error) {
	var selections []NutritionistUserSelection
	err := r.db.Select(&selections, `
		SELECT nus.* FROM nutritionist_user_selections nus
		JOIN individual_orders io ON nus.order_id = io.id
		WHERE nus.date = $1 AND io.paid = false`,
		date.Format("2006-01-02"))
	return selections, err
}

// Daily Menu Reset Flag Methods
func (r *Repository) SetDailyMenuResetFlag(date time.Time, reset bool) error {
	_, err := r.db.Exec(`UPDATE daily_menus SET nutritionist_reset = $1 WHERE date = $2`,
		reset, date.Format("2006-01-02"))
	return err
}

func (r *Repository) GetDailyMenuResetFlag(date time.Time) (bool, error) {
	var resetFlag bool
	err := r.db.Get(&resetFlag, `SELECT nutritionist_reset FROM daily_menus WHERE date = $1`,
		date.Format("2006-01-02"))
	return resetFlag, err
}

// Stock Empty and Notification methods
func (r *Repository) CreateStockEmptyItem(menuItemID int) error {
	_, err := r.db.Exec(`INSERT INTO stock_empty_items (menu_item_id, date) VALUES ($1, CURRENT_DATE)
		ON CONFLICT (menu_item_id, date) DO NOTHING`, menuItemID)
	return err
}

func (r *Repository) DeleteStockEmptyItem(menuItemID int) error {
	_, err := r.db.Exec(`DELETE FROM stock_empty_items WHERE menu_item_id = $1 AND date = CURRENT_DATE`, menuItemID)
	return err
}

func (r *Repository) GetStockEmptyItems() ([]StockEmptyItem, error) {
	var items []StockEmptyItem
	err := r.db.Select(&items, `SELECT * FROM stock_empty_items WHERE date = CURRENT_DATE`)
	return items, err
}

func (r *Repository) CreateUserStockEmptyNotification(employeeID, menuItemID int) error {
	_, err := r.db.Exec(`INSERT INTO user_stock_empty_notifications (individual_order_id, menu_item_id, date)
		VALUES (0, $1, CURRENT_DATE) ON CONFLICT DO NOTHING`, menuItemID)
	return err
}

func (r *Repository) GetUsersNeedingNotification(date time.Time) ([]int, error) {
	var userIDs []int
	err := r.db.Select(&userIDs, `SELECT DISTINCT employee_id FROM user_stock_empty_notifications WHERE date = $1`, date.Format("2006-01-02"))
	return userIDs, err
}

func (r *Repository) GetOrderItemsByOrderID(orderID int) ([]MenuItem, error) {
	var items []MenuItem
	query := `
		SELECT mi.* FROM menu_items mi
		JOIN unnest((SELECT menu_item_ids FROM individual_orders WHERE id = $1)::int[]) item_id ON mi.id = item_id
		ORDER BY mi.name
	`
	err := r.db.Select(&items, query, orderID)
	return items, err
}

func (r *Repository) GetOrderByID(id int) (*IndividualOrder, error) {
	var order IndividualOrder
	err := r.db.Get(&order, `SELECT * FROM individual_orders WHERE id = $1`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &order, err
}

func (r *Repository) GetOrderSessionByID(id int) (*OrderSession, error) {
	var session OrderSession
	err := r.db.Get(&session, `SELECT * FROM order_sessions WHERE id = $1`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &session, err
}

func (r *Repository) GetEmployeeWithCompany(id int) (*EmployeeWithCompany, error) {
	var employee EmployeeWithCompany
	err := r.db.Get(&employee, `
		SELECT e.*, c.name as company_name
		FROM employees e
		JOIN companies c ON e.company_id = c.id
		WHERE e.id = $1 AND e.active = true`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &employee, err
}

func (r *Repository) MarkItemsStockEmpty(itemIDs []int, date time.Time, orderID int) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("Error rolling back transaction: %v", err)
		}
	}()

	// Get the order to find employee ID
	order, err := r.GetOrderByID(orderID)
	if err != nil {
		return err
	}

	// Mark items as stock empty for this specific user/order only
	for _, itemID := range itemIDs {
		// Create user-specific stock empty notification
		_, err = tx.Exec(`
			INSERT INTO user_stock_empty_notifications (individual_order_id, menu_item_id, date) 
			VALUES ($1, $2, $3) 
			ON CONFLICT (individual_order_id, menu_item_id) DO NOTHING`,
			orderID, itemID, date.Format("2006-01-02"))
		if err != nil {
			return err
		}

		// Get menu item name for notification
		var itemName string
		err = tx.Get(&itemName, `SELECT name FROM menu_items WHERE id = $1`, itemID)
		if err != nil {
			return err
		}

		// Get the employee's company ID for the redirect URL
		var companyID int
		err = tx.Get(&companyID, `SELECT company_id FROM employees WHERE id = $1`, order.EmployeeID)
		if err != nil {
			return err
		}

		// Create notification for the user
		redirectURL := fmt.Sprintf("/order/%d/%s", companyID, date.Format("2006-01-02"))
		_, err = tx.Exec(`
			INSERT INTO user_notifications (employee_id, notification_type, title, message, redirect_url) 
			VALUES ($1, $2, $3, $4, $5)`,
			order.EmployeeID,
			NotificationStockEmpty,
			"Item Out of Stock",
			"Unfortunately, "+itemName+" is out of stock. Please update your order with alternative items.",
			&redirectURL)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *Repository) GetUserNotifications(employeeID int, limit int) ([]UserNotification, error) {
	var notifications []UserNotification
	query := `SELECT * FROM user_notifications WHERE employee_id = $1 ORDER BY created_at DESC`
	if limit > 0 {
		query += fmt.Sprintf(` LIMIT %d`, limit)
	}
	err := r.db.Select(&notifications, query, employeeID)
	return notifications, err
}

func (r *Repository) MarkNotificationRead(id int) error {
	_, err := r.db.Exec(`UPDATE user_notifications SET is_read = true WHERE id = $1`, id)
	return err
}

func (r *Repository) DeleteUserNotification(id int) error {
	_, err := r.db.Exec(`DELETE FROM user_notifications WHERE id = $1`, id)
	return err
}

func (r *Repository) CreateUserNotification(employeeID int, notificationType, message string, redirectURL *string) error {
	_, err := r.db.Exec(`
		INSERT INTO user_notifications (employee_id, notification_type, message, redirect_url)
		VALUES ($1, $2, $3, $4)`,
		employeeID, notificationType, message, redirectURL)
	return err
}

func (r *Repository) GetStockEmptyItemsForOrder(orderID int) ([]int, error) {
	var itemIDs []int
	err := r.db.Select(&itemIDs, `
		SELECT DISTINCT menu_item_id 
		FROM user_stock_empty_notifications 
		WHERE individual_order_id = $1`, orderID)
	return itemIDs, err
}

// Get stock empty items for a specific user on a specific date
func (r *Repository) GetStockEmptyItemsForUser(employeeID int, date time.Time) ([]int, error) {
	var itemIDs []int
	err := r.db.Select(&itemIDs, `
		SELECT DISTINCT usn.menu_item_id
		FROM user_stock_empty_notifications usn
		JOIN individual_orders io ON usn.individual_order_id = io.id
		JOIN order_sessions os ON io.session_id = os.id
		WHERE io.employee_id = $1 AND os.date = $2`,
		employeeID, date.Format("2006-01-02"))
	return itemIDs, err
}

func (r *Repository) UnmarkItemStockEmpty(itemID int, date time.Time, orderID int) error {
	// Simply remove the user-specific stock empty notification
	// No need for global stock management since we're using user-specific logic only
	_, err := r.db.Exec(`
		DELETE FROM user_stock_empty_notifications 
		WHERE individual_order_id = $1 AND menu_item_id = $2`,
		orderID, itemID)

	return err
}

func (r *Repository) DeleteAllUserNotifications(employeeID int) error {
	_, err := r.db.Exec(`DELETE FROM user_notifications WHERE employee_id = $1`, employeeID)
	return err
}

func (r *Repository) DeleteUserNotificationsByType(employeeID int, notificationType string) error {
	_, err := r.db.Exec(`DELETE FROM user_notifications WHERE employee_id = $1 AND notification_type = $2`,
		employeeID, notificationType)
	return err
}

func (r *Repository) DeleteUserNotificationsByTypes(employeeID int, notificationTypes []string) error {
	if len(notificationTypes) == 0 {
		return nil
	}

	// Build placeholders for the IN clause
	placeholders := make([]string, len(notificationTypes))
	args := make([]interface{}, len(notificationTypes)+1)
	args[0] = employeeID

	for i, notificationType := range notificationTypes {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		args[i+1] = notificationType
	}

	query := fmt.Sprintf(`DELETE FROM user_notifications WHERE employee_id = $1 AND notification_type IN (%s)`,
		strings.Join(placeholders, ","))

	_, err := r.db.Exec(query, args...)
	return err
}
