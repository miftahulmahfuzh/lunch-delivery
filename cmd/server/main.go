// cmd/server/main.go
package main

import (
	"html/template"
	"log"
	"strings"

	"github.com/miftahulmahfuzh/lunch-delivery/internal/database"
	"github.com/miftahulmahfuzh/lunch-delivery/internal/handlers"
	"github.com/miftahulmahfuzh/lunch-delivery/internal/models"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.NewConnection("localhost", "5432", "lunch_user", "1234", "lunch_delivery")
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	repo := models.NewRepository(db.DB)

	r := gin.Default()
	// Add template functions
	funcMap := template.FuncMap{
		"divideBy100": func(n int) int {
			return n / 100
		},
		"lower": strings.ToLower,
	}

	r.SetFuncMap(funcMap)
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	handlers.SetupRoutes(r, repo)

	log.Println("Server starting on :8080")
	r.Run(":8080")
}
