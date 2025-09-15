package interfaces

import (
	"time"

	"github.com/miftahulmahfuzh/lunch-delivery/internal/models"
)

// RepositoryInterface defines all repository methods
type RepositoryInterface interface {
	// Menu Items
	CreateMenuItem(name string, price int) (*models.MenuItem, error)
	GetAllMenuItems() ([]models.MenuItem, error)
	UpdateMenuItem(id int, name string, price int) error
	DeleteMenuItem(id int) error
	GetMenuItemsByIDs(ids []int64) ([]models.MenuItem, error)

	// Companies
	CreateCompany(name, address, contact string) (*models.Company, error)
	GetAllCompanies() ([]models.Company, error)
	GetCompanyByID(id int) (*models.Company, error)
	UpdateCompany(id int, name, address, contact string) error
	DeleteCompany(id int) error

	// Employees
	CreateEmployee(companyID int, name, email, waContact, passwordHash string) (*models.Employee, error)
	GetEmployeeByEmail(email string) (*models.Employee, error)
	GetEmployeeByID(id int) (*models.Employee, error)
	GetEmployeesByCompany(companyID int) ([]models.Employee, error)
	GetEmployeeWithCompany(id int) (*models.EmployeeWithCompany, error)
	GetRecentOrdersByEmployee(employeeID int, startDate, endDate time.Time) ([]models.RecentOrder, error)
	UpdateEmployee(id int, name, email, waContact string) error
	UpdateEmployeePassword(id int, passwordHash string) error
	DeleteEmployee(id int) error

	// Order Sessions
	CreateOrderSession(companyID int, date time.Time, status string) (*models.OrderSession, error)
	GetOrderSession(companyID int, date time.Time) (*models.OrderSession, error)
	GetOrderSessionsByDate(date time.Time) ([]models.OrderSession, error)
	GetOrderSessionsByDateWithCompany(date time.Time) ([]models.OrderSessionWithCompany, error)
	GetOrderSessionWithCompany(id int) (*models.OrderSessionWithCompany, error)
	GetOrderSessionByID(id int) (*models.OrderSession, error)
	CloseOrderSession(id int) error
	ReopenOrderSession(id int) error

	// Individual Orders
	CreateIndividualOrder(sessionID, employeeID int, menuItemIDs []int64, totalPrice int) (*models.IndividualOrder, error)
	GetOrdersBySession(sessionID int) ([]models.IndividualOrder, error)
	GetOrdersBySessionWithDetails(sessionID int) ([]models.IndividualOrderWithDetails, error)
	GetOrderByID(orderID int) (*models.IndividualOrder, error)
	GetOrderItemsByOrderID(orderID int) ([]models.MenuItem, error)
	GetStockEmptyItemsForOrder(orderID int) ([]int, error)
	MarkOrderPaid(orderID int) error
	MarkOrderUnpaid(orderID int) error
	UpdateOrderStatus(orderID int, status string) error
	MarkItemsStockEmpty(itemIDs []int, date time.Time, orderID int) error
	UnmarkItemStockEmpty(itemID int, date time.Time, orderID int) error

	// Daily Menu
	CreateDailyMenu(date time.Time, menuItemIDs []int64) (*models.DailyMenu, error)
	GetDailyMenuByDate(date time.Time) (*models.DailyMenu, error)

	// Nutritionist Selection
	CreateNutritionistSelection(date time.Time, menuItemIDs []int64, selectedIndices []int32, reasoning string, nutritionalSummary string) (*models.NutritionistSelection, error)
	GetNutritionistSelectionByDate(date time.Time) (*models.NutritionistSelection, error)
	DeleteNutritionistSelection(date time.Time) error

	// Nutritionist User Selection
	CreateNutritionistUserSelection(date time.Time, employeeID int, menuItemIDs []int64) (*models.NutritionistUserSelection, error)
	GetNutritionistUserSelectionByDate(employeeID int, date time.Time) (*models.NutritionistUserSelection, error)
	GetNutritionistUsersByDateAndUnpaid(date time.Time) ([]models.NutritionistUserSelection, error)

	// Stock Empty Items
	CreateStockEmptyItem(menuItemID int) error
	DeleteStockEmptyItem(menuItemID int) error
	GetStockEmptyItems() ([]models.StockEmptyItem, error)
	GetStockEmptyItemsForUser(employeeID int, date time.Time) ([]int, error)

	// User Stock Empty Notifications
	CreateUserStockEmptyNotification(employeeID, menuItemID int) error
	GetUsersNeedingNotification(date time.Time) ([]int, error)

	// User Notifications
	CreateUserNotification(employeeID int, notificationType, message string, redirectURL *string) error
	GetUserNotifications(employeeID int, limit int) ([]models.UserNotification, error)
	MarkNotificationRead(notificationID int) error
	DeleteUserNotification(notificationID int) error
	DeleteUserNotificationsByType(employeeID int, notificationType string) error
	DeleteUserNotificationsByTypes(employeeID int, notificationTypes []string) error
	DeleteAllUserNotifications(employeeID int) error

	// Daily Menu Reset Flag
	GetDailyMenuResetFlag(date time.Time) (bool, error)
	SetDailyMenuResetFlag(date time.Time, value bool) error

	// Password Reset Tokens
	CreatePasswordResetToken(employeeID int, token string, expiresAt time.Time) error
	GetPasswordResetToken(token string) (*models.PasswordResetToken, error)
	MarkPasswordResetTokenAsUsed(token string) error
	DeletePasswordResetToken(token string) error
}
