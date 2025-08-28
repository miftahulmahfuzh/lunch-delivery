// internal/handlers/handlers.go
package handlers

import (
	"net/http"

	"github.com/miftahulmahfuzh/lunch-delivery/internal/middleware"
	"github.com/miftahulmahfuzh/lunch-delivery/internal/models"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo *models.Repository
}

func NewHandler(repo *models.Repository) *Handler {
	return &Handler{repo: repo}
}

// internal/handlers/handlers.go
func SetupRoutes(r *gin.Engine, repo *models.Repository) {
	h := NewHandler(repo)

	// Root redirect
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/login")
	})

	// Public routes
	r.GET("/login", h.loginForm)
	r.POST("/login", h.login)

	// Protected customer routes
	customer := r.Group("/")
	customer.Use(middleware.RequireAuth())
	{
		customer.GET("/logout", h.logout)
		customer.GET("/order/:company/:date", h.orderForm)
		customer.POST("/order", h.submitOrder)
		customer.GET("/my-orders", h.myOrders)
	}

	// Admin routes
	admin := r.Group("/admin")
	{
		admin.GET("/", h.adminDashboard)
		admin.GET("/menu", h.menuList)
		admin.POST("/menu", h.createMenuItem)
		admin.PUT("/menu/:id", h.updateMenuItem)
		admin.DELETE("/menu/:id", h.deleteMenuItem)
		admin.GET("/companies", h.companiesList)
		admin.POST("/companies", h.createCompany)
		admin.GET("/companies/:id/employees", h.companyEmployees)
		admin.POST("/employees", h.createEmployee)
		admin.GET("/daily-menu", h.dailyMenuForm)
		admin.POST("/daily-menu", h.createDailyMenu)
		admin.GET("/sessions", h.orderSessionsList)
		admin.POST("/sessions", h.createOrderSession)
		admin.POST("/sessions/:id/close", h.closeOrderSession)
		admin.POST("/sessions/:id/reopen", h.reopenOrderSession)
		admin.GET("/sessions/:id/orders", h.viewSessionOrders)
		admin.POST("/orders/:id/paid", h.markOrderPaid)
		admin.POST("/orders/:id/unpaid", h.markOrderUnpaid)
		admin.PUT("/companies/:id", h.updateCompany)
		admin.DELETE("/companies/:id", h.deleteCompany)
	}
}
