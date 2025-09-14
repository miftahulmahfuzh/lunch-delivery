package mocks

import (
	"time"

	"github.com/miftahulmahfuzh/lunch-delivery/internal/interfaces"
	"github.com/miftahulmahfuzh/lunch-delivery/internal/models"
	"github.com/stretchr/testify/mock"
)

// RepositoryMock is a mock implementation of the RepositoryInterface
type RepositoryMock struct {
	mock.Mock
}

// Compile-time check to ensure RepositoryMock implements RepositoryInterface
var _ interfaces.RepositoryInterface = (*RepositoryMock)(nil)

// Menu Items
func (m *RepositoryMock) CreateMenuItem(name string, price int) (*models.MenuItem, error) {
	args := m.Called(name, price)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MenuItem), args.Error(1)
}

func (m *RepositoryMock) GetAllMenuItems() ([]models.MenuItem, error) {
	args := m.Called()
	return args.Get(0).([]models.MenuItem), args.Error(1)
}

func (m *RepositoryMock) UpdateMenuItem(id int, name string, price int) error {
	args := m.Called(id, name, price)
	return args.Error(0)
}

func (m *RepositoryMock) DeleteMenuItem(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *RepositoryMock) GetMenuItemsByIDs(ids []int64) ([]models.MenuItem, error) {
	args := m.Called(ids)
	return args.Get(0).([]models.MenuItem), args.Error(1)
}

// Companies
func (m *RepositoryMock) CreateCompany(name, address, contact string) (*models.Company, error) {
	args := m.Called(name, address, contact)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Company), args.Error(1)
}

func (m *RepositoryMock) GetAllCompanies() ([]models.Company, error) {
	args := m.Called()
	return args.Get(0).([]models.Company), args.Error(1)
}

func (m *RepositoryMock) GetCompanyByID(id int) (*models.Company, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Company), args.Error(1)
}

func (m *RepositoryMock) UpdateCompany(id int, name, address, contact string) error {
	args := m.Called(id, name, address, contact)
	return args.Error(0)
}

func (m *RepositoryMock) DeleteCompany(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// Employees
func (m *RepositoryMock) CreateEmployee(companyID int, name, email, waContact, passwordHash string) (*models.Employee, error) {
	args := m.Called(companyID, name, email, waContact, passwordHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Employee), args.Error(1)
}

func (m *RepositoryMock) GetEmployeeByEmail(email string) (*models.Employee, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Employee), args.Error(1)
}

func (m *RepositoryMock) GetEmployeeByID(id int) (*models.Employee, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Employee), args.Error(1)
}

func (m *RepositoryMock) GetEmployeesByCompany(companyID int) ([]models.Employee, error) {
	args := m.Called(companyID)
	return args.Get(0).([]models.Employee), args.Error(1)
}

func (m *RepositoryMock) UpdateEmployee(id int, name, email, waContact string) error {
	args := m.Called(id, name, email, waContact)
	return args.Error(0)
}

func (m *RepositoryMock) DeleteEmployee(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *RepositoryMock) UpdateEmployeePassword(employeeID int, passwordHash string) error {
	args := m.Called(employeeID, passwordHash)
	return args.Error(0)
}

func (m *RepositoryMock) GetEmployeeWithCompany(id int) (*models.EmployeeWithCompany, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EmployeeWithCompany), args.Error(1)
}

// Password Reset Tokens
func (m *RepositoryMock) CreatePasswordResetToken(employeeID int, token string, expiresAt time.Time) (*models.PasswordResetToken, error) {
	args := m.Called(employeeID, token, expiresAt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PasswordResetToken), args.Error(1)
}

func (m *RepositoryMock) GetPasswordResetToken(token string) (*models.PasswordResetToken, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PasswordResetToken), args.Error(1)
}

func (m *RepositoryMock) MarkPasswordResetTokenAsUsed(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *RepositoryMock) CleanupExpiredPasswordResetTokens() error {
	args := m.Called()
	return args.Error(0)
}

// Daily Menu
func (m *RepositoryMock) CreateDailyMenu(date time.Time, menuItemIDs []int64) (*models.DailyMenu, error) {
	args := m.Called(date, menuItemIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DailyMenu), args.Error(1)
}

func (m *RepositoryMock) GetDailyMenuByDate(date time.Time) (*models.DailyMenu, error) {
	args := m.Called(date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DailyMenu), args.Error(1)
}

func (m *RepositoryMock) SetDailyMenuResetFlag(date time.Time, reset bool) error {
	args := m.Called(date, reset)
	return args.Error(0)
}

func (m *RepositoryMock) GetDailyMenuResetFlag(date time.Time) (bool, error) {
	args := m.Called(date)
	return args.Bool(0), args.Error(1)
}

// Order Sessions
func (m *RepositoryMock) CreateOrderSession(companyID int, date time.Time) (*models.OrderSession, error) {
	args := m.Called(companyID, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.OrderSession), args.Error(1)
}

func (m *RepositoryMock) GetOrderSession(companyID int, date time.Time) (*models.OrderSession, error) {
	args := m.Called(companyID, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.OrderSession), args.Error(1)
}

func (m *RepositoryMock) GetOrderSessionByID(id int) (*models.OrderSession, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.OrderSession), args.Error(1)
}

func (m *RepositoryMock) GetOrderSessionsByDate(date time.Time) ([]models.OrderSession, error) {
	args := m.Called(date)
	return args.Get(0).([]models.OrderSession), args.Error(1)
}

func (m *RepositoryMock) GetOrderSessionsByDateWithCompany(date time.Time) ([]models.OrderSessionWithCompany, error) {
	args := m.Called(date)
	return args.Get(0).([]models.OrderSessionWithCompany), args.Error(1)
}

func (m *RepositoryMock) GetOrderSessionWithCompany(id int) (*models.OrderSessionWithCompany, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.OrderSessionWithCompany), args.Error(1)
}

func (m *RepositoryMock) CloseOrderSession(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *RepositoryMock) ReopenOrderSession(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *RepositoryMock) UpdateOrderSessionStatus(id int, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}

// Individual Orders
func (m *RepositoryMock) CreateIndividualOrder(sessionID, employeeID int, menuItemIDs []int64, totalPrice int) (*models.IndividualOrder, error) {
	args := m.Called(sessionID, employeeID, menuItemIDs, totalPrice)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.IndividualOrder), args.Error(1)
}

func (m *RepositoryMock) GetOrdersBySession(sessionID int) ([]models.IndividualOrder, error) {
	args := m.Called(sessionID)
	return args.Get(0).([]models.IndividualOrder), args.Error(1)
}

func (m *RepositoryMock) GetOrdersBySessionWithDetails(sessionID int) ([]models.IndividualOrderWithDetails, error) {
	args := m.Called(sessionID)
	return args.Get(0).([]models.IndividualOrderWithDetails), args.Error(1)
}

func (m *RepositoryMock) GetOrderByID(id int) (*models.IndividualOrder, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.IndividualOrder), args.Error(1)
}

func (m *RepositoryMock) MarkOrderPaid(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *RepositoryMock) MarkOrderUnpaid(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *RepositoryMock) UpdateOrderStatus(id int, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *RepositoryMock) GetRecentOrdersByEmployee(employeeID int, startDate, endDate time.Time) ([]models.RecentOrder, error) {
	args := m.Called(employeeID, startDate, endDate)
	return args.Get(0).([]models.RecentOrder), args.Error(1)
}

func (m *RepositoryMock) GetOrderItemsByOrderID(orderID int) ([]models.MenuItem, error) {
	args := m.Called(orderID)
	return args.Get(0).([]models.MenuItem), args.Error(1)
}

// Nutritionist Selections
func (m *RepositoryMock) GetNutritionistSelectionByDate(date time.Time) (*models.NutritionistSelection, error) {
	args := m.Called(date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.NutritionistSelection), args.Error(1)
}

func (m *RepositoryMock) CreateNutritionistSelection(date time.Time, menuItemIDs []int64, selectedIndices []int32, reasoning, nutritionalSummary string) (*models.NutritionistSelection, error) {
	args := m.Called(date, menuItemIDs, selectedIndices, reasoning, nutritionalSummary)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.NutritionistSelection), args.Error(1)
}

func (m *RepositoryMock) DeleteNutritionistSelection(date time.Time) error {
	args := m.Called(date)
	return args.Error(0)
}

func (m *RepositoryMock) CreateNutritionistUserSelection(employeeID int, date time.Time, orderID *int) error {
	args := m.Called(employeeID, date, orderID)
	return args.Error(0)
}

func (m *RepositoryMock) GetNutritionistUsersByDateAndUnpaid(date time.Time) ([]models.NutritionistUserSelection, error) {
	args := m.Called(date)
	return args.Get(0).([]models.NutritionistUserSelection), args.Error(1)
}

// Stock and Notifications
func (m *RepositoryMock) MarkItemsStockEmpty(itemIDs []int, date time.Time, orderID int) error {
	args := m.Called(itemIDs, date, orderID)
	return args.Error(0)
}

func (m *RepositoryMock) UnmarkItemStockEmpty(itemID int, date time.Time, orderID int) error {
	args := m.Called(itemID, date, orderID)
	return args.Error(0)
}

func (m *RepositoryMock) GetStockEmptyItemsForOrder(orderID int) ([]int, error) {
	args := m.Called(orderID)
	return args.Get(0).([]int), args.Error(1)
}

func (m *RepositoryMock) GetStockEmptyItemsForUser(employeeID int, date time.Time) ([]int, error) {
	args := m.Called(employeeID, date)
	return args.Get(0).([]int), args.Error(1)
}

// User Notifications
func (m *RepositoryMock) CreateUserNotification(employeeID int, notificationType, title, message string, redirectURL *string) error {
	args := m.Called(employeeID, notificationType, title, message, redirectURL)
	return args.Error(0)
}

func (m *RepositoryMock) GetUserNotifications(employeeID int, limit int) ([]models.UserNotification, error) {
	args := m.Called(employeeID, limit)
	return args.Get(0).([]models.UserNotification), args.Error(1)
}

func (m *RepositoryMock) MarkNotificationRead(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *RepositoryMock) DeleteUserNotification(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *RepositoryMock) DeleteAllUserNotifications(employeeID int) error {
	args := m.Called(employeeID)
	return args.Error(0)
}

func (m *RepositoryMock) DeleteUserNotificationsByType(employeeID int, notificationType string) error {
	args := m.Called(employeeID, notificationType)
	return args.Error(0)
}

func (m *RepositoryMock) DeleteUserNotificationsByTypes(employeeID int, notificationTypes []string) error {
	args := m.Called(employeeID, notificationTypes)
	return args.Error(0)
}