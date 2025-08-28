// internal/handlers/admin.go
package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) adminDashboard(c *gin.Context) {
	today := time.Now()
	sessions, err := h.repo.GetOrderSessionsByDateWithCompany(today)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "admin_dashboard.html", gin.H{
		"sessions": sessions,
		"date":     today.Format("2006-01-02"),
	})
}

// Menu Items
func (h *Handler) menuList(c *gin.Context) {
	items, err := h.repo.GetAllMenuItems()
	if err != nil {
		log.Printf("Error getting menu items: %v", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	log.Printf("Found %d menu items", len(items))
	c.HTML(http.StatusOK, "menu_list.html", gin.H{"items": items})
}

func (h *Handler) createMenuItem(c *gin.Context) {
	name := c.PostForm("name")
	priceStr := c.PostForm("price")

	if name == "" || priceStr == "" {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Name and price required"})
		return
	}

	price, err := strconv.Atoi(priceStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid price"})
		return
	}

	_, err = h.repo.CreateMenuItem(name, price*100) // Convert to cents
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/admin/menu")
}

func (h *Handler) updateMenuItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	name := c.PostForm("name")
	priceStr := c.PostForm("price")

	if name == "" || priceStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and price required"})
		return
	}

	price, err := strconv.Atoi(priceStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price"})
		return
	}

	err = h.repo.UpdateMenuItem(id, name, price*100)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) deleteMenuItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.repo.DeleteMenuItem(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// Companies
func (h *Handler) companiesList(c *gin.Context) {
	companies, err := h.repo.GetAllCompanies()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "companies_list.html", gin.H{"companies": companies})
}

func (h *Handler) createCompany(c *gin.Context) {
	name := c.PostForm("name")
	address := c.PostForm("address")
	contact := c.PostForm("contact")

	if name == "" {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Company name required"})
		return
	}

	_, err := h.repo.CreateCompany(name, address, contact)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/admin/companies")
}

func (h *Handler) companyEmployees(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid company ID"})
		return
	}

	company, err := h.repo.GetCompanyByID(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}
	if company == nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Company not found"})
		return
	}

	employees, err := h.repo.GetEmployeesByCompany(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "company_employees.html", gin.H{
		"company":   company,
		"employees": employees,
	})
}

// Daily Menu
func (h *Handler) dailyMenuForm(c *gin.Context) {
	items, err := h.repo.GetAllMenuItems()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	today := time.Now()
	existingMenu, _ := h.repo.GetDailyMenuByDate(today)

	c.HTML(http.StatusOK, "daily_menu_form.html", gin.H{
		"items":    items,
		"existing": existingMenu,
		"date":     today.Format("2006-01-02"),
	})
}

func (h *Handler) createDailyMenu(c *gin.Context) {
	dateStr := c.PostForm("date")
	menuItems := c.PostFormArray("menu_items")

	if dateStr == "" || len(menuItems) == 0 {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Date and menu items required"})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid date format"})
		return
	}

	var itemIDs []int64
	for _, idStr := range menuItems {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid menu item ID"})
			return
		}
		itemIDs = append(itemIDs, id)
	}

	_, err = h.repo.CreateDailyMenu(date, itemIDs)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/admin/daily-menu")
}

// Order Sessions
func (h *Handler) orderSessionsList(c *gin.Context) {
	today := time.Now()
	sessions, err := h.repo.GetOrderSessionsByDateWithCompany(today)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	companies, err := h.repo.GetAllCompanies()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "order_sessions.html", gin.H{
		"sessions":  sessions,
		"companies": companies,
		"date":      today.Format("2006-01-02"),
	})
}

func (h *Handler) createOrderSession(c *gin.Context) {
	companyIDStr := c.PostForm("company_id")
	dateStr := c.PostForm("date")

	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid company ID"})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid date format"})
		return
	}

	_, err = h.repo.CreateOrderSession(companyID, date)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/admin/sessions")
}

func (h *Handler) closeOrderSession(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	err = h.repo.CloseOrderSession(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) reopenOrderSession(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	err = h.repo.ReopenOrderSession(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) viewSessionOrders(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid session ID"})
		return
	}

	orders, err := h.repo.GetOrdersBySessionWithDetails(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	// Calculate summary
	totalRevenue := 0
	paidCount := 0
	for _, order := range orders {
		totalRevenue += order.TotalPrice
		if order.Paid {
			paidCount++
		}
	}

	c.HTML(http.StatusOK, "session_orders.html", gin.H{
		"orders":       orders,
		"session_id":   id,
		"totalRevenue": totalRevenue,
		"paidCount":    paidCount,
	})
}

func (h *Handler) markOrderPaid(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	err = h.repo.MarkOrderPaid(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
