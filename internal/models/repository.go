// internal/models/repository.go
package models

import (
	"database/sql"
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
		`INSERT INTO daily_menus (date, menu_item_ids) VALUES ($1, $2)
         ON CONFLICT (date) DO UPDATE SET menu_item_ids = $2
         RETURNING *`,
		date.Format("2006-01-02"), pq.Array(menuItemIDs))
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
