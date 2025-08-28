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
