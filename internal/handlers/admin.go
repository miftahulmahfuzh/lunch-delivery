// internal/handlers/admin.go
package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/miftahulmahfuzh/lunch-delivery/internal/models"
	"github.com/rs/zerolog/log"
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

	// Calculate statistics if we have items
	templateData := gin.H{"items": items}
	if len(items) > 0 {
		var total int = 0
		minPrice := items[0].Price
		maxPrice := items[0].Price

		for _, item := range items {
			total += item.Price
			if item.Price < minPrice {
				minPrice = item.Price
			}
			if item.Price > maxPrice {
				maxPrice = item.Price
			}
		}

		averagePrice := total / len(items)
		templateData["averagePrice"] = averagePrice
		templateData["minPrice"] = minPrice
		templateData["maxPrice"] = maxPrice
	}

	c.HTML(http.StatusOK, "menu_list.html", templateData)
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

	_, err = h.repo.CreateMenuItem(name, price) // Store rupiah directly
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

	err = h.repo.UpdateMenuItem(id, name, price)
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

	// Get employees with existing unpaid orders for this date to notify them
	sessions, _ := h.repo.GetOrderSessionsByDate(date)
	employeesToNotify := make(map[int]bool)
	
	for _, session := range sessions {
		orders, _ := h.repo.GetOrdersBySession(session.ID)
		for _, order := range orders {
			if !order.Paid {
				employeesToNotify[order.EmployeeID] = true
			}
		}
	}

	_, err = h.repo.CreateDailyMenu(date, itemIDs)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	// Create notifications for affected employees with correct company-specific URLs
	for employeeID := range employeesToNotify {
		// Get employee's company ID
		employee, err := h.repo.GetEmployeeByID(employeeID)
		if err != nil {
			log.Warn().Err(err).Int("employee_id", employeeID).Msg("Failed to get employee for notification")
			continue
		}
		
		redirectURL := fmt.Sprintf("/order/%d/%s", employee.CompanyID, date.Format("2006-01-02"))
		err = h.repo.CreateUserNotification(
			employeeID,
			models.NotificationMenuUpdated,
			"Menu Updated",
			"The daily menu has been updated for "+date.Format("January 2, 2006")+". Please review your order as it may be affected.",
			&redirectURL,
		)
		if err != nil {
			log.Warn().Err(err).Int("employee_id", employeeID).Msg("Failed to create menu update notification")
		}
	}

	c.Redirect(http.StatusFound, "/admin/daily-menu")
}

// Order Sessions
func (h *Handler) orderSessionsList(c *gin.Context) {
	// Get date from query parameter, default to today
	dateStr := c.Query("date")
	var targetDate time.Time
	var err error

	if dateStr != "" {
		targetDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			targetDate = time.Now()
		}
	} else {
		targetDate = time.Now()
	}

	sessions, err := h.repo.GetOrderSessionsByDateWithCompany(targetDate)
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
		"sessions":     sessions,
		"companies":    companies,
		"date":         targetDate.Format("2006-01-02"),
		"selectedDate": targetDate,
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

	// Get session details for notifications
	session, err := h.repo.GetOrderSessionWithCompany(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get all orders for this session to notify employees
	orders, err := h.repo.GetOrdersBySession(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.repo.CloseOrderSession(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create notifications for all employees who placed orders
	for _, order := range orders {
		err = h.repo.CreateUserNotification(
			order.EmployeeID,
			models.NotificationSessionClosed,
			"Order Session Closed",
			"The lunch order session for "+session.CompanyName+" has been closed. Your order is now being processed.",
			nil,
		)
		if err != nil {
			log.Warn().Err(err).Int("employee_id", order.EmployeeID).Msg("Failed to create session closed notification")
		}
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

// Update viewSessionOrders in internal/handlers/admin.go
func (h *Handler) viewSessionOrders(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid session ID"})
		return
	}

	// Get session with company name
	session, err := h.repo.GetOrderSessionWithCompany(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}
	if session == nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Session not found"})
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
		"session":      session,
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

	// Get order details for notification
	order, err := h.repo.GetOrderByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.repo.MarkOrderPaid(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create notification for the employee
	err = h.repo.CreateUserNotification(
		order.EmployeeID,
		models.NotificationPaid,
		"Payment Confirmed",
		"Your lunch order payment has been confirmed. Thank you!",
		nil,
	)
	if err != nil {
		// Log error but don't fail the request
		log.Warn().Err(err).Msg("Failed to create payment notification")
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// Add to internal/handlers/admin.go
func (h *Handler) markOrderUnpaid(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	err = h.repo.MarkOrderUnpaid(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// Companies Page
func (h *Handler) updateCompany(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	name := c.PostForm("name")
	address := c.PostForm("address")
	contact := c.PostForm("contact")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Company name required"})
		return
	}

	err = h.repo.UpdateCompany(id, name, address, contact)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) deleteCompany(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.repo.DeleteCompany(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// New handlers for stock empty and employee details functionality
func (h *Handler) getOrderItems(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	items, err := h.repo.GetOrderItemsByOrderID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get stock empty items for this specific order
	stockEmptyItemIDs, err := h.repo.GetStockEmptyItemsForOrder(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create a map for faster lookup
	stockEmptyMap := make(map[int]bool)
	for _, itemID := range stockEmptyItemIDs {
		stockEmptyMap[itemID] = true
	}

	// Add stock empty status to each item
	type ItemWithStatus struct {
		models.MenuItem
		IsStockEmpty bool `json:"is_stock_empty"`
	}

	var itemsWithStatus []ItemWithStatus
	for _, item := range items {
		itemsWithStatus = append(itemsWithStatus, ItemWithStatus{
			MenuItem:     item,
			IsStockEmpty: stockEmptyMap[item.ID],
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"items":   itemsWithStatus,
	})
}

func (h *Handler) markItemsStockEmpty(c *gin.Context) {
	idStr := c.Param("id")
	orderID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var req struct {
		ItemIDs []int `json:"item_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if len(req.ItemIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No items provided"})
		return
	}

	// Get the order to find employee and date
	order, err := h.repo.GetOrderByID(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get order"})
		return
	}

	session, err := h.repo.GetOrderSessionByID(order.SessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get session"})
		return
	}

	// Mark items as stock empty and create notifications
	err = h.repo.MarkItemsStockEmpty(req.ItemIDs, session.Date, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) unmarkItemsStockEmpty(c *gin.Context) {
	idStr := c.Param("id")
	orderID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var req struct {
		ItemIDs []int `json:"item_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if len(req.ItemIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No items provided"})
		return
	}

	// Get the order to find session and date
	order, err := h.repo.GetOrderByID(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get order"})
		return
	}

	session, err := h.repo.GetOrderSessionByID(order.SessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get session"})
		return
	}

	// Unmark items as stock empty
	for _, itemID := range req.ItemIDs {
		err = h.repo.UnmarkItemStockEmpty(itemID, session.Date, orderID)
		if err != nil {
			log.Warn().Err(err).Int("item_id", itemID).Msg("Failed to unmark item as stock empty")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmark some items"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) getEmployeeDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	employee, err := h.repo.GetEmployeeWithCompany(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if employee == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"employee": employee,
	})
}
