// internal/handlers/orders.go
package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/miftahulmahfuzh/lunch-delivery/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) orderForm(c *gin.Context) {
	// Get user from cookie
	userIDStr, err := c.Cookie("user_id")
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	companyIDStr := c.Param("company")
	dateStr := c.Param("date")

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

	// Get order session
	session, err := h.repo.GetOrderSession(companyID, date)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}
	if session == nil || session.Status != models.StatusOpen {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Order session not available"})
		return
	}

	// Get daily menu
	dailyMenu, err := h.repo.GetDailyMenuByDate(date)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}
	if dailyMenu == nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "No menu available for this date"})
		return
	}

	// Get available menu items
	var itemIDs []int64
	for _, id := range dailyMenu.MenuItemIDs {
		itemIDs = append(itemIDs, id)
	}

	menuItems, err := h.repo.GetMenuItemsByIDs(itemIDs)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	// Check if user already has an order
	orders, err := h.repo.GetOrdersBySession(session.ID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	var existingOrder *models.IndividualOrder
	for _, order := range orders {
		if order.EmployeeID == userID {
			existingOrder = &order
			break
		}
	}

	c.HTML(http.StatusOK, "order_form.html", gin.H{
		"menu_items":     menuItems,
		"session":        session,
		"existing_order": existingOrder,
		"user_id":        userID,
	})
}

func (h *Handler) submitOrder(c *gin.Context) {
	userIDStr, err := c.Cookie("user_id")
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	sessionIDStr := c.PostForm("session_id")
	menuItems := c.PostFormArray("menu_items")

	sessionID, err := strconv.Atoi(sessionIDStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid session ID"})
		return
	}

	if len(menuItems) == 0 {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "No menu items selected"})
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

	// Calculate total price
	items, err := h.repo.GetMenuItemsByIDs(itemIDs)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	totalPrice := 0
	for _, item := range items {
		totalPrice += item.Price
	}

	_, err = h.repo.CreateIndividualOrder(sessionID, userID, itemIDs, totalPrice)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/my-orders")
}

func (h *Handler) myOrders(c *gin.Context) {
	userIDStr, err := c.Cookie("user_id")
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// This is simplified - in reality you'd want to get recent orders
	c.HTML(http.StatusOK, "my_orders.html", gin.H{
		"user_id": userID,
	})
}
