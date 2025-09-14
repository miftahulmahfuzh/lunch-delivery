// internal/handlers/orders.go
package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/miftahulmahfuzh/lunch-delivery/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Helper functions for date calculations
func getStartOfWeek(date time.Time) time.Time {
	day := date.Weekday()
	diff := date.AddDate(0, 0, -int(day)+1) // Monday as start of week
	if day == time.Sunday {
		diff = date.AddDate(0, 0, -6) // If Sunday, go back to Monday
	}
	return time.Date(diff.Year(), diff.Month(), diff.Day(), 0, 0, 0, 0, diff.Location())
}

func getEndOfWeek(date time.Time) time.Time {
	startOfWeek := getStartOfWeek(date)
	return startOfWeek.AddDate(0, 0, 6) // Sunday as end of week
}

// orderRedirect handles the generic /order route and redirects to the proper order form
func (h *Handler) orderRedirect(c *gin.Context) {
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

	// Get the employee to find their company
	employee, err := h.repo.GetEmployeeByID(userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get employee")
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Database error"})
		return
	}

	// Get today's date
	today := time.Now()
	todayStr := today.Format("2006-01-02")

	// Check if there's an order session for today
	_, err = h.repo.GetOrderSession(employee.CompanyID, today)
	if err != nil {
		// No session for today, redirect to my-orders with an informational message
		c.Redirect(http.StatusFound, "/my-orders")
		return
	}

	// Redirect to the proper order form
	c.Redirect(http.StatusFound, "/order/"+strconv.Itoa(employee.CompanyID)+"/"+todayStr)
}

// Update orderForm handler in internal/handlers/orders.go
func (h *Handler) orderForm(c *gin.Context) {
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

	// Get company info
	company, err := h.repo.GetCompanyByID(companyID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}
	if company == nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Company not found"})
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

	// Prevent editing if order is already paid
	if existingOrder != nil && existingOrder.Paid {
		c.HTML(http.StatusForbidden, "error.html", gin.H{"error": "Cannot modify order: Payment has already been processed"})
		return
	}

	// Check if menu was reset and user needs notification
	showResetNotification := false
	if dailyMenu.NutritionistReset {
		// Check if this user has used nutritionist selection for today and order is unpaid
		unpaidUsers, _ := h.nutritionistService.GetUsersNeedingNotification(date)
		for _, userSel := range unpaidUsers {
			if userSel.EmployeeID == userID {
				showResetNotification = true
				break
			}
		}
	}

	// Get stock empty items for this specific user on this date
	stockEmptyItems, err := h.repo.GetStockEmptyItemsForUser(userID, date)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get user-specific stock empty items")
		stockEmptyItems = []int{} // Continue with empty list
	}

	c.HTML(http.StatusOK, "order_form.html", gin.H{
		"title":                 "Place Your Order",
		"menu_items":            menuItems,
		"session":               session,
		"company":               company,
		"existing_order":        existingOrder,
		"user_id":               userID,
		"show_reset_notification": showResetNotification,
		"stock_empty_items":     stockEmptyItems,
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

	// Check if user already has a paid order for this session
	orders, err := h.repo.GetOrdersBySession(sessionID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	for _, order := range orders {
		if order.EmployeeID == userID && order.Paid {
			c.HTML(http.StatusForbidden, "error.html", gin.H{"error": "Cannot modify order: Payment has already been processed"})
			return
		}
	}

	_, err = h.repo.CreateIndividualOrder(sessionID, userID, itemIDs, totalPrice)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/my-orders")
}

// Replace the empty myOrders handler in internal/handlers/orders.go
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

	// Get employee info
	employee, err := h.repo.GetEmployeeByID(userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}
	if employee == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Get company info
	company, err := h.repo.GetCompanyByID(employee.CompanyID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	// Get today's session for this company
	today := time.Now()
	todaySession, _ := h.repo.GetOrderSession(employee.CompanyID, today)

	var todayOrder *models.IndividualOrder
	var todayOrderItems []models.MenuItem
	var stockEmptyItems []models.MenuItem
	var adjustedTotalPrice int
	var stockReductionAmount int

	if todaySession != nil {
		// Get user's order for today
		orders, err := h.repo.GetOrdersBySession(todaySession.ID)
		if err == nil {
			for _, order := range orders {
				if order.EmployeeID == userID {
					// Get fresh order data from database to ensure payment status is current
					freshOrder, freshErr := h.repo.GetOrderByID(order.ID)
					if freshErr == nil && freshOrder != nil {
						todayOrder = freshOrder
					} else {
						todayOrder = &order
					}
					break
				}
			}
		}

		// Get menu items for today's order
		if todayOrder != nil && len(todayOrder.MenuItemIDs) > 0 {
			var itemIDs []int64
			for _, id := range todayOrder.MenuItemIDs {
				itemIDs = append(itemIDs, id)
			}
			allOrderItems, _ := h.repo.GetMenuItemsByIDs(itemIDs)
			
			// Get stock empty items for this user today
			stockEmptyItemIDs, _ := h.repo.GetStockEmptyItemsForUser(userID, today)
			stockEmptyMap := make(map[int]bool)
			for _, itemID := range stockEmptyItemIDs {
				stockEmptyMap[itemID] = true
			}
			
			// Filter out stock empty items from display and calculate adjusted totals
			adjustedTotalPrice = todayOrder.TotalPrice
			stockReductionAmount = 0
			for _, item := range allOrderItems {
				if stockEmptyMap[item.ID] {
					stockEmptyItems = append(stockEmptyItems, item)
					adjustedTotalPrice -= item.Price
					stockReductionAmount += item.Price
				} else {
					todayOrderItems = append(todayOrderItems, item)
				}
			}
		}
	}

	// Get recent orders based on date range
	var recentOrders []models.RecentOrder
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	
	if startDateStr != "" && endDateStr != "" {
		// Parse dates from query parameters
		startDate, err1 := time.Parse("2006-01-02", startDateStr)
		endDate, err2 := time.Parse("2006-01-02", endDateStr)
		
		if err1 == nil && err2 == nil {
			recentOrders, _ = h.repo.GetRecentOrdersByEmployee(userID, startDate, endDate)
		}
	} else {
		// Default to "this week" if no parameters provided
		today := time.Now()
		startOfWeek := getStartOfWeek(today)
		endOfWeek := getEndOfWeek(today)
		recentOrders, _ = h.repo.GetRecentOrdersByEmployee(userID, startOfWeek, endOfWeek)
	}

	// Get user notifications
	notifications, _ := h.repo.GetUserNotifications(userID, 10) // Get latest 10 notifications

	c.HTML(http.StatusOK, "my_orders.html", gin.H{
		"employee":             employee,
		"company":              company,
		"todaySession":         todaySession,
		"todayOrder":           todayOrder,
		"todayOrderItems":      todayOrderItems,
		"stockEmptyItems":      stockEmptyItems,
		"adjustedTotalPrice":   adjustedTotalPrice,
		"stockReductionAmount": stockReductionAmount,
		"recentOrders":         recentOrders,
		"notifications":        notifications,
		"startDate":            startDateStr,
		"endDate":              endDateStr,
	})
}

// Nutritionist selection handler
func (h *Handler) nutritionistSelect(c *gin.Context) {
	userIDStr, err := c.Cookie("user_id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not logged in"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	dateStr := c.Param("date")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	// Get daily menu
	dailyMenu, err := h.repo.GetDailyMenuByDate(date)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get daily menu")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get menu"})
		return
	}
	if dailyMenu == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No menu available for this date"})
		return
	}

	// Get available menu items
	var itemIDs []int64
	for _, id := range dailyMenu.MenuItemIDs {
		itemIDs = append(itemIDs, id)
	}

	menuItems, err := h.repo.GetMenuItemsByIDs(itemIDs)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get menu items")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get menu items"})
		return
	}

	if len(menuItems) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No menu items available"})
		return
	}

	// Call nutritionist service with user ID to handle user-specific stock constraints
	ctx := c.Request.Context()
	selection, err := h.nutritionistService.GetNutritionistSelection(ctx, date, menuItems, userID)
	if err != nil {
		log.Error().Err(err).Msg("Nutritionist selection failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Nutritionist selection failed"})
		return
	}

	// Track that this user used nutritionist selection
	if err := h.nutritionistService.TrackUserSelection(userID, date, nil); err != nil {
		log.Warn().Err(err).Msg("Failed to track user selection, but continuing")
	}

	log.Info().
		Interface("selected_indices", selection.SelectedIndices).
		Str("reasoning", selection.Reasoning).
		Int("user_id", userID).
		Msg("Nutritionist selection successful")

	c.JSON(http.StatusOK, gin.H{
		"success":           true,
		"selected_indices":  selection.SelectedIndices,
		"reasoning":         selection.Reasoning,
		"nutritional_summary": selection.NutritionalSummary,
	})
}

// Notification handlers
func (h *Handler) markNotificationRead(c *gin.Context) {
	userIDStr, err := c.Cookie("user_id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not logged in"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	idStr := c.Param("id")
	notificationID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	// Verify notification belongs to user (security check)
	notifications, err := h.repo.GetUserNotifications(userID, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	found := false
	for _, notification := range notifications {
		if notification.ID == notificationID {
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusForbidden, gin.H{"error": "Notification not found"})
		return
	}

	err = h.repo.MarkNotificationRead(notificationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) deleteNotification(c *gin.Context) {
	userIDStr, err := c.Cookie("user_id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not logged in"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	idStr := c.Param("id")
	notificationID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	// Verify notification belongs to user (security check)
	notifications, err := h.repo.GetUserNotifications(userID, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	found := false
	for _, notification := range notifications {
		if notification.ID == notificationID {
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusForbidden, gin.H{"error": "Notification not found"})
		return
	}

	err = h.repo.DeleteUserNotification(notificationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) clearAllNotifications(c *gin.Context) {
	userIDStr, err := c.Cookie("user_id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not logged in"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.repo.DeleteAllUserNotifications(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) clearStockEmptyNotifications(c *gin.Context) {
	userIDStr, err := c.Cookie("user_id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not logged in"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.repo.DeleteUserNotificationsByType(userID, "STOCK_EMPTY")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) clearMenuRelatedNotifications(c *gin.Context) {
	userIDStr, err := c.Cookie("user_id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not logged in"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.repo.DeleteUserNotificationsByTypes(userID, []string{"STOCK_EMPTY", "MENU_UPDATED"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
