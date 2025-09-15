# Handlers Package Unit Tests Guide

## Overview

This guide explains the unit tests for the `internal/handlers` package, which contains HTTP request handlers for the lunch delivery system's web interface, including authentication, ordering, admin functionality, and API endpoints.

## Test Structure

### Files Tested
- `auth.go` - Authentication handlers (login, signup, password reset)
- `admin.go` - Administrative functionality (menu, companies, orders, sessions)
- `orders.go` - Customer ordering and notification management
- `employees.go` - Employee management operations
- `handlers.go` - Route setup and handler initialization

### Test Files
- `auth_test.go` - Authentication handler tests (IMPLEMENTED)
- Missing test files for other handlers (TO BE IMPLEMENTED)

## Test Categories

### 1. Authentication Tests (`auth_test.go`) - ✅ IMPLEMENTED

#### A. Login Handler Tests (`TestHandler_Login`)
**Purpose**: Test user authentication functionality.

**Scenarios Tested**:
- ✅ Successful login with valid credentials
- ✅ Login fails with invalid email (user not found)
- ✅ Login fails with invalid password
- ✅ Login fails with missing credentials
- ✅ Login handles database errors gracefully

**Test Patterns**:
- Mock repository setup with `mocks.RepositoryMock`
- Password hashing validation using `bcrypt`
- Cookie verification (user_id, company_id)
- Redirect validation to `/my-orders`
- Template panic handling with `SafeHandlerCall`

#### B. Logout Handler Tests (`TestHandler_Logout`)
**Purpose**: Test user logout functionality.

**Scenarios Tested**:
- ✅ Logout clears cookies and redirects to login page

**Test Patterns**:
- Cookie clearing verification (MaxAge -1)
- Redirect validation to `/login`

#### C. Signup Form Handler Tests (`TestHandler_SignupForm`)
**Purpose**: Test signup form display functionality.

**Scenarios Tested**:
- ✅ Displays signup form with companies list
- ✅ Handles database error when fetching companies

**Test Patterns**:
- Company data mocking with `testutils.MockCompany`
- Database error simulation

#### D. Signup Handler Tests (`TestHandler_Signup`)
**Purpose**: Test user registration functionality.

**Scenarios Tested**:
- ✅ Successful signup creates new employee
- ✅ Signup fails when passwords don't match
- ✅ Signup fails when email already exists
- ✅ Signup fails with short password (validation)

**Test Patterns**:
- Form data validation
- Email uniqueness checking
- Password confirmation validation
- Employee creation verification

#### E. Login Form Handler Tests (`TestHandler_LoginForm`)
**Purpose**: Test login form display.

**Scenarios Tested**:
- ✅ Displays login form successfully

### 2. Password Reset Tests - ❌ MISSING TESTS

#### Missing Test Coverage (`auth.go` functions):
- `forgotPasswordForm()` - Display forgot password form
- `forgotPassword()` - Process forgot password requests
- `resetPasswordForm()` - Display reset password form
- `resetPassword()` - Process password reset

**Should Test**:
- Form display functionality
- Email validation for password reset
- Token generation and validation
- Password reset token expiration
- Database error handling

### 3. Admin Handler Tests - ❌ MISSING TESTS

#### Missing Test Coverage (`admin.go` functions):
- `adminDashboard()` - Admin dashboard display
- `menuList()` - Menu item listing
- `createMenuItem()` - Menu item creation
- `updateMenuItem()` - Menu item updates
- `deleteMenuItem()` - Menu item deletion
- `companiesList()` - Company listing
- `createCompany()` - Company creation
- `companyEmployees()` - Employee listing by company
- `dailyMenuForm()` - Daily menu form
- `createDailyMenu()` - Daily menu creation
- `orderSessionsList()` - Order session management
- `createOrderSession()` - Session creation
- `closeOrderSession()` - Session closure
- `reopenOrderSession()` - Session reopening
- `viewSessionOrders()` - Session order viewing
- `markOrderPaid()` - Payment status updates
- `markOrderUnpaid()` - Payment status updates
- `updateOrderStatus()` - Order status management
- `updateCompany()` - Company updates
- `deleteCompany()` - Company deletion
- `getOrderItems()` - Order item retrieval
- `markItemsStockEmpty()` - Stock management
- `unmarkItemsStockEmpty()` - Stock management
- `getEmployeeDetails()` - Employee detail retrieval
- `getEmptyStockItemsForOrder()` - Stock item queries

### 4. Order Handler Tests - ❌ MISSING TESTS

#### Missing Test Coverage (`orders.go` functions):
- `orderRedirect()` - Order page redirection logic
- `orderForm()` - Order form display
- `submitOrder()` - Order submission processing
- `myOrders()` - User order history
- `nutritionistSelect()` - AI nutritionist selection
- `markNotificationRead()` - Notification management
- `deleteNotification()` - Notification deletion
- `clearAllNotifications()` - Bulk notification clearing
- `clearStockEmptyNotifications()` - Stock notification clearing
- `clearMenuRelatedNotifications()` - Menu notification clearing

### 5. Employee Handler Tests - ❌ MISSING TESTS

#### Missing Test Coverage (`employees.go` functions):
- `createEmployee()` - Employee creation (admin)
- `updateEmployee()` - Employee updates
- `deleteEmployee()` - Employee deletion

### 6. Route Setup Tests - ❌ MISSING TESTS

#### Missing Test Coverage (`handlers.go` functions):
- `NewHandler()` - Handler initialization
- `SetupRoutes()` - Route configuration
- Root redirect functionality
- Favicon serving
- Middleware integration testing

## Running the Tests

### Run All Handler Tests
```bash
go test ./internal/handlers
```

### Run Specific Test Categories
```bash
# Authentication tests only
go test ./internal/handlers -run TestHandler_Login
go test ./internal/handlers -run TestHandler_Signup

# Run with verbose output
go test ./internal/handlers -v
```

### Run Tests with Coverage
```bash
go test ./internal/handlers -cover
```

## Test Utilities Used

### From `testutils` Package
- `SetupGinTest()` - Gin test mode initialization
- `CreateFormRequest()` - HTTP form request creation
- `CreateTestGinContext()` - Gin context creation
- `SafeHandlerCall()` - Template panic handling
- `MockEmployee()` - Employee test data
- `MockCompany()` - Company test data

### From `mocks` Package
- `RepositoryMock` - Repository interface mocking
- Mock expectations with `On()` method
- Assertion verification with `AssertExpectations()`

### External Libraries
- `testify/assert` - Assertions
- `testify/mock` - Mock objects
- `gin-gonic/gin` - HTTP testing
- `golang.org/x/crypto/bcrypt` - Password hashing

## Example Test Execution

```bash
$ go test ./internal/handlers -v
=== RUN   TestHandler_Login
=== RUN   TestHandler_Login/successful_login_with_valid_credentials
=== RUN   TestHandler_Login/login_fails_with_invalid_email
=== RUN   TestHandler_Login/login_fails_with_invalid_password
=== RUN   TestHandler_Login/login_fails_with_missing_credentials
=== RUN   TestHandler_Login/login_handles_database_error
--- PASS: TestHandler_Login (0.05s)
=== RUN   TestHandler_Logout
--- PASS: TestHandler_Logout (0.01s)
=== RUN   TestHandler_SignupForm
--- PASS: TestHandler_SignupForm (0.02s)
=== RUN   TestHandler_Signup
--- PASS: TestHandler_Signup (0.03s)
=== RUN   TestHandler_LoginForm
--- PASS: TestHandler_LoginForm (0.01s)
PASS
ok      github.com/miftahulmahfuzh/lunch-delivery/internal/handlers    0.125s
```

## Test Coverage Analysis

### Currently Implemented (auth_test.go)
- ✅ Authentication flow (login/logout)
- ✅ User registration (signup)
- ✅ Form display handlers
- ✅ Input validation
- ✅ Database error handling
- ✅ Cookie management
- ✅ Redirect behavior

### Missing Test Coverage
- ❌ Password reset functionality (4 handlers)
- ❌ Admin functionality (20+ handlers)
- ❌ Order management (9 handlers)
- ❌ Employee management (3 handlers)
- ❌ Route setup and initialization
- ❌ Middleware integration
- ❌ JSON API endpoints
- ❌ File serving (favicon)

**Estimated Test Coverage**: ~15% (5 out of 33+ handlers tested)

## Recommended Test Implementation Priority

### High Priority (Core Functionality)
1. **Order Handlers**: `submitOrder()`, `orderForm()`, `myOrders()`
2. **Admin Core**: `createMenuItem()`, `updateMenuItem()`, `deleteMenuItem()`
3. **Session Management**: `createOrderSession()`, `closeOrderSession()`
4. **Password Reset**: `forgotPassword()`, `resetPassword()`

### Medium Priority (Management Features)
1. **Company Management**: `createCompany()`, `updateCompany()`, `deleteCompany()`
2. **Employee Management**: All employee handlers
3. **Daily Menu**: `createDailyMenu()`, `dailyMenuForm()`
4. **Order Status**: `markOrderPaid()`, `updateOrderStatus()`

### Lower Priority (Supporting Features)
1. **Notifications**: All notification handlers
2. **Stock Management**: `markItemsStockEmpty()`, `unmarkItemsStockEmpty()`
3. **Nutritionist**: `nutritionistSelect()`
4. **Admin Dashboard**: `adminDashboard()`, `viewSessionOrders()`

## Common Test Patterns for Missing Tests

### 1. Admin Handler Testing Pattern
```go
func TestHandler_CreateMenuItem(t *testing.T) {
    testutils.SetupGinTest()

    t.Run("successful menu item creation", func(t *testing.T) {
        mockRepo := &mocks.RepositoryMock{}
        handler := NewHandler(mockRepo, nil)

        newItem := testutils.MockMenuItem(1, "Test Item", 15000)
        mockRepo.On("CreateMenuItem", "Test Item", 15000, "test description").
            Return(&newItem, nil)

        formData := map[string]string{
            "name":        "Test Item",
            "price":       "15000",
            "description": "test description",
        }
        req := testutils.CreateFormRequest("POST", "/admin/menu", formData)

        ctx, recorder := testutils.CreateTestGinContext("POST", "/admin/menu", nil)
        ctx.Request = req

        testutils.SafeHandlerCall(func() {
            handler.createMenuItem(ctx)
        }, recorder, http.StatusOK)

        assert.Equal(t, http.StatusOK, recorder.Code)
        mockRepo.AssertExpectations(t)
    })
}
```

### 2. Order Handler Testing Pattern
```go
func TestHandler_SubmitOrder(t *testing.T) {
    testutils.SetupGinTest()

    t.Run("successful order submission", func(t *testing.T) {
        mockRepo := &mocks.RepositoryMock{}
        handler := NewHandler(mockRepo, nil)

        // Setup session and menu mocks
        session := testutils.MockOrderSession(1, 1, testutils.TestDate())
        mockRepo.On("GetOrderSession", 1, testutils.TestDate()).
            Return(&session, nil)

        // Setup order creation mock
        newOrder := testutils.MockIndividualOrder(1, 1, 1)
        mockRepo.On("CreateIndividualOrder", 1, 1, mock.AnythingOfType("[]int")).
            Return(&newOrder, nil)

        formData := map[string]string{
            "menu_items": "1,2,3",
        }
        req := testutils.CreateFormRequest("POST", "/order", formData)

        ctx, recorder := testutils.CreateTestGinContext("POST", "/order", nil)
        ctx.Request = req
        // Set cookies for authentication
        ctx.Request.AddCookie(&http.Cookie{Name: "user_id", Value: "1"})
        ctx.Request.AddCookie(&http.Cookie{Name: "company_id", Value: "1"})

        testutils.SafeHandlerCall(func() {
            handler.submitOrder(ctx)
        }, recorder, http.StatusOK)

        assert.Equal(t, http.StatusOK, recorder.Code)
        mockRepo.AssertExpectations(t)
    })
}
```

## Maintenance Guidelines

### When to Update Tests
- When adding new handler functions
- When modifying handler request/response formats
- When changing authentication logic
- When updating business rules
- When adding new form fields or validation

### Best Practices
- Always use `testutils.SetupGinTest()` for Gin initialization
- Use `SafeHandlerCall()` for handlers that render templates
- Mock repository dependencies properly
- Test both success and error scenarios
- Verify HTTP status codes and response headers
- Test cookie handling for authentication
- Include input validation tests
- Test database error handling

### Common Issues
- **Template Panics**: Use `SafeHandlerCall()` for template rendering
- **Cookie Testing**: Verify both setting and clearing of cookies
- **Form Data**: Use `CreateFormRequest()` for proper form encoding
- **Authentication**: Mock user sessions with cookies
- **Database Errors**: Test error propagation from repository layer
- **Redirects**: Check both status codes and Location headers

## Integration Points

These tests support:
- Middleware testing (authentication, authorization)
- Service layer testing (through handler integration)
- API endpoint validation
- Template rendering verification
- Business logic validation
- User experience testing