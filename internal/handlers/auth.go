// internal/handlers/auth.go
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) loginForm(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func (h *Handler) login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	if email == "" || password == "" {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": "Email and password required"})
		return
	}

	employee, err := h.repo.GetEmployeeByEmail(email)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}
	if employee == nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid credentials"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(employee.PasswordHash), []byte(password))
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid credentials"})
		return
	}

	// Set session
	c.SetCookie("user_id", strconv.Itoa(employee.ID), 3600*24, "/", "", false, true)
	c.SetCookie("company_id", strconv.Itoa(employee.CompanyID), 3600*24, "/", "", false, true)

	c.Redirect(http.StatusFound, "/my-orders")
}

func (h *Handler) logout(c *gin.Context) {
	c.SetCookie("user_id", "", -1, "/", "", false, true)
	c.SetCookie("company_id", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}

func (h *Handler) signupForm(c *gin.Context) {
	companies, err := h.repo.GetAllCompanies()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "signup.html", gin.H{
		"companies": companies,
	})
}

func (h *Handler) signup(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")
	waContact := c.PostForm("wa_contact")
	companyIDStr := c.PostForm("company_id")
	password := c.PostForm("password")
	confirmPassword := c.PostForm("confirm_password")

	// Get companies for form redisplay
	companies, _ := h.repo.GetAllCompanies()

	// Form data for redisplay on error (excluding passwords for security)
	formData := gin.H{
		"name":       name,
		"email":      email,
		"wa_contact": waContact,
		"company_id": companyIDStr,
	}

	// Validation
	if name == "" || email == "" || waContact == "" || companyIDStr == "" || password == "" {
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{
			"error":     "All fields are required",
			"companies": companies,
			"form":      formData,
		})
		return
	}

	if password != confirmPassword {
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{
			"error":     "Passwords don't match",
			"companies": companies,
			"form":      formData,
		})
		return
	}

	if len(password) < 6 {
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{
			"error":     "Password must be at least 6 characters",
			"companies": companies,
			"form":      formData,
		})
		return
	}

	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{
			"error":     "Invalid company selection",
			"companies": companies,
			"form":      formData,
		})
		return
	}

	// Check if email already exists
	existingEmployee, err := h.repo.GetEmployeeByEmail(email)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "signup.html", gin.H{
			"error":     "System error. Please try again.",
			"companies": companies,
			"form":      formData,
		})
		return
	}
	if existingEmployee != nil {
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{
			"error":     "Email already registered",
			"companies": companies,
			"form":      formData,
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "signup.html", gin.H{
			"error":     "System error. Please try again.",
			"companies": companies,
			"form":      formData,
		})
		return
	}

	// Create employee
	_, err = h.repo.CreateEmployee(companyID, name, email, waContact, string(hashedPassword))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "signup.html", gin.H{
			"error":     "Failed to create account. Please try again.",
			"companies": companies,
			"form":      formData,
		})
		return
	}

	c.HTML(http.StatusOK, "signup.html", gin.H{
		"success":   "Account created successfully! You can now sign in.",
		"companies": companies,
	})
}
