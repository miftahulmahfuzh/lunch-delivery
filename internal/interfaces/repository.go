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

	// Employees
	CreateEmployee(companyID int, name, email, waContact, passwordHash string) (*models.Employee, error)
	GetEmployeeByEmail(email string) (*models.Employee, error)
	UpdateEmployeePassword(id int, passwordHash string) error

	// Order Sessions
	CreateOrderSession(companyID int, date time.Time, status string) (*models.OrderSession, error)
	GetOrderSession(companyID int, date time.Time) (*models.OrderSession, error)
	CloseOrderSession(id int) error

	// Individual Orders
	CreateIndividualOrder(sessionID, employeeID int, menuItemIDs []int64, totalPrice int) (*models.IndividualOrder, error)
	GetOrdersBySession(sessionID int) ([]models.IndividualOrder, error)
	MarkOrderPaid(orderID int) error
	UpdateOrderStatus(orderID int, status string) error

	// Daily Menu
	CreateDailyMenu(date time.Time, menuItemIDs []int64) (*models.DailyMenu, error)
	GetDailyMenuByDate(date time.Time) (*models.DailyMenu, error)

	// Nutritionist Selection
	CreateNutritionistSelection(date time.Time, menuItemIDs []int64) (*models.NutritionistSelection, error)
	GetNutritionistSelectionByDate(date time.Time) (*models.NutritionistSelection, error)
	DeleteNutritionistSelection(date time.Time) error

	// Nutritionist User Selection
	CreateNutritionistUserSelection(date time.Time, employeeID int, menuItemIDs []int64) (*models.NutritionistUserSelection, error)
	GetNutritionistUserSelectionByDate(employeeID int, date time.Time) (*models.NutritionistUserSelection, error)

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
	DeleteAllUserNotifications(employeeID int) error

	// Daily Menu Reset Flag
	GetDailyMenuResetFlag(date time.Time) (bool, error)
	SetDailyMenuResetFlag(date time.Time, value bool) error

	// Password Reset Tokens
	CreatePasswordResetToken(employeeID int, token string, expiresAt time.Time) error
	GetPasswordResetToken(token string) (*models.PasswordResetToken, error)
	DeletePasswordResetToken(token string) error
}