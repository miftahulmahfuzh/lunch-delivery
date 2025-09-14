package handlers

import (
	"net/http"
	"testing"

	"github.com/miftahulmahfuzh/lunch-delivery/internal/mocks"
	"github.com/miftahulmahfuzh/lunch-delivery/internal/models"
	"github.com/miftahulmahfuzh/lunch-delivery/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestHandler_Login(t *testing.T) {
	testutils.SetupGinTest()

	t.Run("successful login with valid credentials", func(t *testing.T) {
		mockRepo := &mocks.RepositoryMock{}
		handler := NewHandler(mockRepo, nil)

		// Mock employee with hashed password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
		employee := &models.Employee{
			ID:           1,
			CompanyID:    1,
			Name:         "John Doe",
			Email:        "john@test.com",
			PasswordHash: string(hashedPassword),
			Active:       true,
		}

		mockRepo.On("GetEmployeeByEmail", "john@test.com").Return(employee, nil)

		// Create form request
		formData := map[string]string{
			"email":    "john@test.com",
			"password": "testpass",
		}
		req := testutils.CreateFormRequest("POST", "/login", formData)

		ctx, recorder := testutils.CreateTestGinContext("POST", "/login", nil)
		ctx.Request = req

		handler.login(ctx)

		assert.Equal(t, http.StatusFound, recorder.Code)
		assert.Equal(t, "/my-orders", recorder.Header().Get("Location"))

		// Check cookies are set
		cookies := recorder.Result().Cookies()
		assert.Len(t, cookies, 2)

		var userIDCookie, companyIDCookie *http.Cookie
		for _, cookie := range cookies {
			switch cookie.Name {
			case "user_id":
				userIDCookie = cookie
			case "company_id":
				companyIDCookie = cookie
			}
		}

		assert.NotNil(t, userIDCookie)
		assert.NotNil(t, companyIDCookie)
		assert.Equal(t, "1", userIDCookie.Value)
		assert.Equal(t, "1", companyIDCookie.Value)

		mockRepo.AssertExpectations(t)
	})

	t.Run("login fails with invalid email", func(t *testing.T) {
		mockRepo := &mocks.RepositoryMock{}
		handler := NewHandler(mockRepo, nil)

		mockRepo.On("GetEmployeeByEmail", "nonexistent@test.com").Return(nil, nil)

		formData := map[string]string{
			"email":    "nonexistent@test.com",
			"password": "testpass",
		}
		req := testutils.CreateFormRequest("POST", "/login", formData)

		ctx, recorder := testutils.CreateTestGinContext("POST", "/login", nil)
		ctx.Request = req

		handler.login(ctx)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Invalid credentials")

		mockRepo.AssertExpectations(t)
	})

	t.Run("login fails with invalid password", func(t *testing.T) {
		mockRepo := &mocks.RepositoryMock{}
		handler := NewHandler(mockRepo, nil)

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpass"), bcrypt.DefaultCost)
		employee := &models.Employee{
			ID:           1,
			CompanyID:    1,
			Email:        "john@test.com",
			PasswordHash: string(hashedPassword),
		}

		mockRepo.On("GetEmployeeByEmail", "john@test.com").Return(employee, nil)

		formData := map[string]string{
			"email":    "john@test.com",
			"password": "wrongpass",
		}
		req := testutils.CreateFormRequest("POST", "/login", formData)

		ctx, recorder := testutils.CreateTestGinContext("POST", "/login", nil)
		ctx.Request = req

		handler.login(ctx)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Invalid credentials")

		mockRepo.AssertExpectations(t)
	})

	t.Run("login fails with missing credentials", func(t *testing.T) {
		mockRepo := &mocks.RepositoryMock{}
		handler := NewHandler(mockRepo, nil)

		formData := map[string]string{
			"email":    "",
			"password": "",
		}
		req := testutils.CreateFormRequest("POST", "/login", formData)

		ctx, recorder := testutils.CreateTestGinContext("POST", "/login", nil)
		ctx.Request = req

		handler.login(ctx)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Email and password required")

		// No repository calls should be made
		mockRepo.AssertExpectations(t)
	})

	t.Run("login handles database error", func(t *testing.T) {
		mockRepo := &mocks.RepositoryMock{}
		handler := NewHandler(mockRepo, nil)

		mockRepo.On("GetEmployeeByEmail", "john@test.com").Return(nil, assert.AnError)

		formData := map[string]string{
			"email":    "john@test.com",
			"password": "testpass",
		}
		req := testutils.CreateFormRequest("POST", "/login", formData)

		ctx, recorder := testutils.CreateTestGinContext("POST", "/login", nil)
		ctx.Request = req

		handler.login(ctx)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)

		mockRepo.AssertExpectations(t)
	})
}

func TestHandler_Logout(t *testing.T) {
	testutils.SetupGinTest()

	t.Run("logout clears cookies and redirects", func(t *testing.T) {
		handler := NewHandler(nil, nil)

		ctx, recorder := testutils.CreateTestGinContext("GET", "/logout", nil)

		handler.logout(ctx)

		assert.Equal(t, http.StatusFound, recorder.Code)
		assert.Equal(t, "/login", recorder.Header().Get("Location"))

		// Check that cookies are cleared (set with MaxAge -1)
		cookies := recorder.Result().Cookies()
		assert.Len(t, cookies, 2)

		for _, cookie := range cookies {
			if cookie.Name == "user_id" || cookie.Name == "company_id" {
				assert.Equal(t, "", cookie.Value)
				assert.Equal(t, -1, cookie.MaxAge)
			}
		}
	})
}

func TestHandler_SignupForm(t *testing.T) {
	testutils.SetupGinTest()

	t.Run("displays signup form with companies", func(t *testing.T) {
		mockRepo := &mocks.RepositoryMock{}
		handler := NewHandler(mockRepo, nil)

		companies := []models.Company{
			testutils.MockCompany(1, "Company A"),
			testutils.MockCompany(2, "Company B"),
		}

		mockRepo.On("GetAllCompanies").Return(companies, nil)

		ctx, recorder := testutils.CreateTestGinContext("GET", "/signup", nil)

		handler.signupForm(ctx)

		assert.Equal(t, http.StatusOK, recorder.Code)
		// In a real implementation, you'd check that the HTML contains the companies

		mockRepo.AssertExpectations(t)
	})

	t.Run("handles database error when fetching companies", func(t *testing.T) {
		mockRepo := &mocks.RepositoryMock{}
		handler := NewHandler(mockRepo, nil)

		mockRepo.On("GetAllCompanies").Return(nil, assert.AnError)

		ctx, recorder := testutils.CreateTestGinContext("GET", "/signup", nil)

		handler.signupForm(ctx)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)

		mockRepo.AssertExpectations(t)
	})
}

func TestHandler_Signup(t *testing.T) {
	testutils.SetupGinTest()

	t.Run("successful signup creates new employee", func(t *testing.T) {
		mockRepo := &mocks.RepositoryMock{}
		handler := NewHandler(mockRepo, nil)

		companies := []models.Company{testutils.MockCompany(1, "Test Company")}
		newEmployee := testutils.MockEmployee(1, 1, "John Doe", "john@test.com")

		mockRepo.On("GetAllCompanies").Return(companies, nil).Times(1)
		mockRepo.On("GetEmployeeByEmail", "john@test.com").Return(nil, nil)
		mockRepo.On("CreateEmployee", 1, "John Doe", "john@test.com", "+628123456789", mock.AnythingOfType("string")).Return(&newEmployee, nil)

		formData := map[string]string{
			"name":             "John Doe",
			"email":            "john@test.com",
			"wa_contact":       "+628123456789",
			"company_id":       "1",
			"password":         "testpass123",
			"confirm_password": "testpass123",
		}
		req := testutils.CreateFormRequest("POST", "/signup", formData)

		ctx, recorder := testutils.CreateTestGinContext("POST", "/signup", nil)
		ctx.Request = req

		handler.signup(ctx)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Account created successfully")

		mockRepo.AssertExpectations(t)
	})

	t.Run("signup fails when passwords don't match", func(t *testing.T) {
		mockRepo := &mocks.RepositoryMock{}
		handler := NewHandler(mockRepo, nil)

		companies := []models.Company{testutils.MockCompany(1, "Test Company")}
		mockRepo.On("GetAllCompanies").Return(companies, nil)

		formData := map[string]string{
			"name":             "John Doe",
			"email":            "john@test.com",
			"wa_contact":       "+628123456789",
			"company_id":       "1",
			"password":         "testpass123",
			"confirm_password": "differentpass",
		}
		req := testutils.CreateFormRequest("POST", "/signup", formData)

		ctx, recorder := testutils.CreateTestGinContext("POST", "/signup", nil)
		ctx.Request = req

		handler.signup(ctx)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Passwords don't match")

		mockRepo.AssertExpectations(t)
	})

	t.Run("signup fails when email already exists", func(t *testing.T) {
		mockRepo := &mocks.RepositoryMock{}
		handler := NewHandler(mockRepo, nil)

		companies := []models.Company{testutils.MockCompany(1, "Test Company")}
		existingEmployee := testutils.MockEmployee(1, 1, "Existing User", "john@test.com")

		mockRepo.On("GetAllCompanies").Return(companies, nil)
		mockRepo.On("GetEmployeeByEmail", "john@test.com").Return(&existingEmployee, nil)

		formData := map[string]string{
			"name":             "John Doe",
			"email":            "john@test.com",
			"wa_contact":       "+628123456789",
			"company_id":       "1",
			"password":         "testpass123",
			"confirm_password": "testpass123",
		}
		req := testutils.CreateFormRequest("POST", "/signup", formData)

		ctx, recorder := testutils.CreateTestGinContext("POST", "/signup", nil)
		ctx.Request = req

		handler.signup(ctx)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Email already registered")

		mockRepo.AssertExpectations(t)
	})

	t.Run("signup fails with short password", func(t *testing.T) {
		mockRepo := &mocks.RepositoryMock{}
		handler := NewHandler(mockRepo, nil)

		companies := []models.Company{testutils.MockCompany(1, "Test Company")}
		mockRepo.On("GetAllCompanies").Return(companies, nil)

		formData := map[string]string{
			"name":             "John Doe",
			"email":            "john@test.com",
			"wa_contact":       "+628123456789",
			"company_id":       "1",
			"password":         "123",
			"confirm_password": "123",
		}
		req := testutils.CreateFormRequest("POST", "/signup", formData)

		ctx, recorder := testutils.CreateTestGinContext("POST", "/signup", nil)
		ctx.Request = req

		handler.signup(ctx)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Password must be at least 6 characters")

		mockRepo.AssertExpectations(t)
	})
}

func TestHandler_LoginForm(t *testing.T) {
	testutils.SetupGinTest()

	t.Run("displays login form", func(t *testing.T) {
		handler := NewHandler(nil, nil)

		ctx, recorder := testutils.CreateTestGinContext("GET", "/login", nil)

		handler.loginForm(ctx)

		assert.Equal(t, http.StatusOK, recorder.Code)
		// In a real implementation, you'd check that the HTML contains the login form
	})
}