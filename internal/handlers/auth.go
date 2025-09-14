// internal/handlers/auth.go
package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/miftahulmahfuzh/lunch-delivery/internal/utils"
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

func (h *Handler) forgotPasswordForm(c *gin.Context) {
	c.HTML(http.StatusOK, "forgot_password.html", nil)
}

func (h *Handler) forgotPassword(c *gin.Context) {
	email := c.PostForm("email")

	if email == "" {
		c.HTML(http.StatusBadRequest, "forgot_password.html", gin.H{
			"error": "Email address is required",
		})
		return
	}

	employee, err := h.repo.GetEmployeeByEmail(email)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "forgot_password.html", gin.H{
			"error": "System error. Please try again.",
		})
		return
	}

	if employee != nil {
		token, err := utils.GeneratePasswordResetToken()
		if err != nil {
			c.HTML(http.StatusInternalServerError, "forgot_password.html", gin.H{
				"error": "System error. Please try again.",
			})
			return
		}

		expiresAt := time.Now().Add(1 * time.Hour)
		err = h.repo.CreatePasswordResetToken(employee.ID, token, expiresAt)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "forgot_password.html", gin.H{
				"error": "System error. Please try again.",
			})
			return
		}

		emailService := utils.NewEmailService()
		baseURL := c.Request.Header.Get("Origin")
		if baseURL == "" {
			baseURL = "http://localhost:8080"
		}

		err = emailService.SendPasswordResetEmail(email, token, baseURL)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "forgot_password.html", gin.H{
				"error": "Failed to send reset email. Please try again.",
			})
			return
		}
	}

	c.HTML(http.StatusOK, "forgot_password.html", gin.H{
		"success": "If an account with that email exists, a password reset link has been sent.",
	})
}

func (h *Handler) resetPasswordForm(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Invalid or missing reset token",
		})
		return
	}

	resetToken, err := h.repo.GetPasswordResetToken(token)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "System error. Please try again.",
		})
		return
	}

	if resetToken == nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Invalid or expired reset token",
		})
		return
	}

	c.HTML(http.StatusOK, "reset_password.html", gin.H{
		"token": token,
	})
}

func (h *Handler) resetPassword(c *gin.Context) {
	token := c.PostForm("token")
	password := c.PostForm("password")
	confirmPassword := c.PostForm("confirm_password")

	if token == "" || password == "" || confirmPassword == "" {
		c.HTML(http.StatusBadRequest, "reset_password.html", gin.H{
			"error": "All fields are required",
			"token": token,
		})
		return
	}

	if password != confirmPassword {
		c.HTML(http.StatusBadRequest, "reset_password.html", gin.H{
			"error": "Passwords don't match",
			"token": token,
		})
		return
	}

	if len(password) < 6 {
		c.HTML(http.StatusBadRequest, "reset_password.html", gin.H{
			"error": "Password must be at least 6 characters",
			"token": token,
		})
		return
	}

	resetToken, err := h.repo.GetPasswordResetToken(token)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "reset_password.html", gin.H{
			"error": "System error. Please try again.",
			"token": token,
		})
		return
	}

	if resetToken == nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Invalid or expired reset token",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "reset_password.html", gin.H{
			"error": "System error. Please try again.",
			"token": token,
		})
		return
	}

	err = h.repo.UpdateEmployeePassword(resetToken.EmployeeID, string(hashedPassword))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "reset_password.html", gin.H{
			"error": "Failed to update password. Please try again.",
			"token": token,
		})
		return
	}

	err = h.repo.MarkPasswordResetTokenAsUsed(token)
	if err != nil {
		// Log error but don't fail the request since password was already updated
	}

	c.HTML(http.StatusOK, "login.html", gin.H{
		"success": "Password reset successfully! You can now sign in with your new password.",
	})
}
