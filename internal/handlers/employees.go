// internal/handlers/employees.go
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Update createEmployee in internal/handlers/employees.go
func (h *Handler) createEmployee(c *gin.Context) {
	companyIDStr := c.PostForm("company_id")
	name := c.PostForm("name")
	email := c.PostForm("email")
	waContact := c.PostForm("wa_contact")
	password := c.PostForm("password")

	if companyIDStr == "" || name == "" || email == "" || waContact == "" || password == "" {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "All fields required"})
		return
	}

	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid company ID"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Password hashing failed"})
		return
	}

	_, err = h.repo.CreateEmployee(companyID, name, email, waContact, string(hashedPassword))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/admin/companies/"+companyIDStr+"/employees")
}
