// internal/middleware/auth.go
package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr, err := c.Cookie("user_id")
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simple admin check - you can make this more sophisticated
		adminKey := c.GetHeader("X-Admin-Key")
		if adminKey != "admin123" { // Change this to something secure
			c.HTML(http.StatusForbidden, "error.html", gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
