// internal/models/repository.go
package models

import (
	"database/sql"
	"fmt"
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
func (r *Repository) CreateOrderSession(companyID int, date time.Time) (*OrderSession, error) {
	var session OrderSession
	err := r.db.Get(&session,
		`INSERT INTO order_sessions (company_id, date) VALUES ($1, $2) RETURNING *`,
		companyID, date.Format("2006-01-02"))
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
		`INSERT INTO individual_orders (session_id, employee_id, menu_item_ids, total_price)
         VALUES ($1, $2, $3, $4)
         ON CONFLICT (session_id, employee_id)
         DO UPDATE SET menu_item_ids = $3, total_price = $4
         RETURNING *`,
		sessionID, employeeID, pq.Array(menuItemIDs), totalPrice)
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

// Add this to internal/models/repository.go

type OrderSessionWithCompany struct {
	OrderSession
	CompanyName string `json:"company_name" db:"company_name"`
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

type IndividualOrderWithDetails struct {
	IndividualOrder
	EmployeeName  string `json:"employee_name" db:"employee_name"`
	MenuItemNames string `json:"menu_item_names"`
}

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

// Add to internal/models/repository.go
func (r *Repository) MarkOrderUnpaid(id int) error {
	_, err := r.db.Exec(`UPDATE individual_orders SET paid = false WHERE id = $1`, id)
	return err
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

type RecentOrder struct {
	Date          time.Time `db:"date"`
	TotalPrice    int       `db:"total_price"`
	Paid          bool      `db:"paid"`
	MenuItemNames string    `json:"menu_item_names"`
}

func (r *Repository) GetRecentOrdersByEmployee(employeeID int, startDate, endDate time.Time) ([]RecentOrder, error) {
	type OrderRow struct {
		Date        time.Time      `db:"date"`
		TotalPrice  int            `db:"total_price"`
		Paid        bool           `db:"paid"`
		MenuItemIDs pq.Int64Array  `db:"menu_item_ids"`
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

func (r *Repository) CreateNutritionistSelection(date time.Time, menuItemIDs []int64, selectedIndices []int32, reasoning, nutritionalSummary string) (*NutritionistSelection, error) {
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
func (r *Repository) CreateNutritionistUserSelection(employeeID int, date time.Time, orderID *int) error {
	_, err := r.db.Exec(
		`INSERT INTO nutritionist_user_selections (employee_id, date, order_id) 
         VALUES ($1, $2, $3)
         ON CONFLICT (employee_id, date) 
         DO UPDATE SET order_id = $3`,
		employeeID, date.Format("2006-01-02"), orderID)
	return err
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

type EmployeeWithCompany struct {
	Employee
	CompanyName string `json:"company_name" db:"company_name"`
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
	defer tx.Rollback()

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

func (r *Repository) CreateUserNotification(employeeID int, notificationType, title, message string, redirectURL *string) error {
	_, err := r.db.Exec(`
		INSERT INTO user_notifications (employee_id, notification_type, title, message, redirect_url) 
		VALUES ($1, $2, $3, $4, $5)`,
		employeeID, notificationType, title, message, redirectURL)
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