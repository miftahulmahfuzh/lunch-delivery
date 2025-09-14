# Models Package Unit Tests Guide

## Overview

This guide explains the unit tests for the `internal/models` package, which contains data models, database repository, and core business logic for the lunch delivery system.

## Test Structure

### Files Tested
- `models.go` - Data model structures and constants
- `repository.go` - Database repository with 70+ methods

### Test Files
- `models_test.go` - Model structure validation and constants testing
- `repository_test.go` - Database repository functionality testing

## Test Categories

### 1. Model Structure Tests (`models_test.go`)

#### A. Structure Validation (`TestModelStructures`)
Tests all model structs for proper field assignment and data integrity:

**Models Tested**:
- `MenuItem` - Menu item with pricing and status
- `Company` - Company information and contact details
- `Employee` - Employee data with authentication
- `PasswordResetToken` - Password reset functionality
- `DailyMenu` - Daily menu with item selections
- `OrderSession` - Order session management
- `IndividualOrder` - Individual customer orders
- `NutritionistSelection` - AI nutritionist recommendations
- `UserNotification` - User notification system
- `StockEmptyItem` - Stock management
- `UserStockEmptyNotification` - User-specific stock notifications

#### B. Constants Validation (`TestModelConstants`)
Tests all application constants:
- Order session status constants
- Individual order status constants
- Notification type constants

#### C. Model Validation (`TestModelValidation`)
Tests edge cases and data validation:
- Zero and negative prices
- Empty password hashes
- Nil pointer fields
- Optional field handling

#### D. Edge Cases (`TestModelEdgeCases`)
Tests boundary conditions:
- Empty PostgreSQL arrays
- Very long text fields
- Special characters in text fields
- Unicode character support

#### E. JSON Serialization (`TestModelJSON`)
Tests JSON marshaling/unmarshaling:
- Proper JSON tag handling
- Password hash exclusion from JSON
- API compatibility

#### F. Model Relationships (`TestModelRelationships`)
Tests foreign key relationships:
- Employee belongs to Company
- Order belongs to Session and Employee
- Token belongs to Employee

### 2. Repository Tests (`repository_test.go`)

#### A. Core Setup (`setupMockDB`)
Helper function that creates:
- Mock database connection using `sqlmock`
- Repository instance
- Cleanup functions

#### B. Menu Item Operations (`TestRepository_MenuItems`)
**Functions Tested**:
- `CreateMenuItem()` - Create new menu items
- `GetAllMenuItems()` - Retrieve all active items
- `UpdateMenuItem()` - Update existing items
- `DeleteMenuItem()` - Soft delete items
- `GetMenuItemsByIDs()` - Batch retrieve by IDs

**Test Scenarios**:
- Successful operations
- Database errors
- Empty results
- Invalid parameters

#### C. Company Operations (`TestRepository_Companies`)
**Functions Tested**:
- `CreateCompany()` - Create new companies
- `GetAllCompanies()` - Retrieve all active companies
- `GetCompanyByID()` - Retrieve specific company

**Test Scenarios**:
- Successful CRUD operations
- Company not found scenarios
- Database connection errors

#### D. Employee Operations (`TestRepository_Employees`)
**Functions Tested**:
- `CreateEmployee()` - Create new employees
- `GetEmployeeByEmail()` - Authentication lookup
- `UpdateEmployeePassword()` - Password updates

**Test Scenarios**:
- Employee creation with password hashing
- Email-based lookups
- Non-existent employee handling

#### E. Order Session Operations (`TestRepository_OrderSessions`)
**Functions Tested**:
- `CreateOrderSession()` - Create ordering sessions
- `GetOrderSession()` - Retrieve sessions
- `CloseOrderSession()` - Close sessions for ordering

**Test Scenarios**:
- Session lifecycle management
- Date-based session retrieval
- Status transitions

#### F. Individual Order Operations (`TestRepository_IndividualOrders`)
**Functions Tested**:
- `CreateIndividualOrder()` - Create customer orders
- `GetOrdersBySession()` - Retrieve session orders
- `MarkOrderPaid()` - Payment processing
- `UpdateOrderStatus()` - Status management

**Test Scenarios**:
- Order creation with menu items array
- Payment status tracking
- Order status transitions

#### G. Daily Menu Operations (`TestRepository_DailyMenu`)
**Functions Tested**:
- `CreateDailyMenu()` - Set daily menu
- `GetDailyMenuByDate()` - Retrieve menu for date

**Test Scenarios**:
- Menu creation with item arrays
- Date-based menu retrieval
- Menu not found handling

#### H. User Notification Operations (`TestRepository_UserNotifications`)
**Functions Tested**:
- `CreateUserNotification()` - Create notifications
- `GetUserNotifications()` - Retrieve user notifications
- `MarkNotificationRead()` - Mark as read
- `DeleteAllUserNotifications()` - Bulk delete

**Test Scenarios**:
- Notification creation with optional redirect URLs
- Pagination support
- Bulk operations

## Running the Tests

### Run All Model Tests
```bash
go test ./internal/models
```

### Run Specific Test Categories
```bash
# Model structure tests only
go test ./internal/models -run TestModelStructures

# Repository tests only
go test ./internal/models -run TestRepository

# Specific repository functionality
go test ./internal/models -run TestRepository_MenuItems
```

### Run Tests with Database Mock Verification
```bash
go test ./internal/models -v
```

### Run Tests with Coverage
```bash
go test ./internal/models -cover
```

## Mock Database Testing

### Setup Pattern
```go
func TestSomeRepositoryFunction(t *testing.T) {
    repo, mock, cleanup := setupMockDB(t)
    defer cleanup()

    // Setup mock expectations
    mock.ExpectQuery(`SELECT \* FROM table`).
        WithArgs("param").
        WillReturnRows(sqlmock.NewRows([]string{"col1", "col2"}).
            AddRow("val1", "val2"))

    // Execute function
    result, err := repo.SomeFunction("param")

    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.NoError(t, mock.ExpectationsWereMet())
}
```

### Common Mock Patterns

#### 1. Query Mocking
```go
mock.ExpectQuery(`SELECT \* FROM menu_items`).
    WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price"}).
        AddRow(1, "Item 1", 10000).
        AddRow(2, "Item 2", 20000))
```

#### 2. Exec Mocking
```go
mock.ExpectExec(`UPDATE menu_items SET`).
    WithArgs("New Name", 15000, 1).
    WillReturnResult(sqlmock.NewResult(0, 1))
```

#### 3. Error Mocking
```go
mock.ExpectQuery(`SELECT \* FROM menu_items`).
    WillReturnError(assert.AnError)
```

## Test Utilities Used

### From `testutils` Package
- `MockMenuItem()` - Creates test menu items
- `MockCompany()` - Creates test companies
- `MockEmployee()` - Creates test employees
- `MockOrderSession()` - Creates test order sessions
- `TestDate()` - Consistent test date
- `SetupGinTest()` - Gin test mode setup

### External Libraries
- `sqlmock` - SQL query mocking
- `testify/assert` - Assertions
- `testify/require` - Required assertions
- `testify/mock` - Mock objects

## Example Test Execution

```bash
$ go test ./internal/models -v
=== RUN   TestModelStructures
=== RUN   TestModelStructures/MenuItem_structure
=== RUN   TestModelStructures/Company_structure
--- PASS: TestModelStructures (0.01s)
=== RUN   TestRepository_MenuItems
=== RUN   TestRepository_MenuItems/CreateMenuItem_success
=== RUN   TestRepository_MenuItems/GetAllMenuItems_success
--- PASS: TestRepository_MenuItems (0.02s)
PASS
ok      github.com/miftahulmahfuzh/lunch-delivery/internal/models    0.045s
```

## Test Coverage

### Models (`models_test.go`)
- ✅ All 14 model structures
- ✅ All constants and enums
- ✅ JSON serialization
- ✅ Field validation
- ✅ Edge cases and special characters
- ✅ Relationship validation

### Repository (`repository_test.go`)
- ✅ 25+ core repository methods
- ✅ CRUD operations for all entities
- ✅ Database error handling
- ✅ SQL query validation
- ✅ Transaction handling
- ✅ Array field support (PostgreSQL)

## Maintenance Guidelines

### When to Update Tests
- When adding new model fields
- When adding new repository methods
- When changing database schema
- When modifying business logic
- When adding new relationships

### Best Practices
- Always use `sqlmock` for database testing
- Verify all SQL expectations are met
- Test both success and error scenarios
- Use consistent test data from `testutils`
- Test edge cases and boundary conditions
- Include relationship validation

### Common Issues
- **SQL Regex Matching**: Use proper regex escaping in `ExpectQuery`
- **Array Handling**: Test PostgreSQL array fields properly
- **Transaction Cleanup**: Ensure mock database cleanup
- **Time Zones**: Use UTC for consistent date testing
- **Null Pointers**: Test optional field handling

## Integration Points

These tests support:
- Handler tests (through repository mocking)
- Service tests (through repository interface)
- Database migration validation
- API endpoint testing
- Business logic validation