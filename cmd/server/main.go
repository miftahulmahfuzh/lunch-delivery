// cmd/server/main.go
package main

import (
	"html/template"
	"log"
	"strings"

	"github.com/miftahulmahfuzh/lunch-delivery/internal/config"
	"github.com/miftahulmahfuzh/lunch-delivery/internal/database"
	"github.com/miftahulmahfuzh/lunch-delivery/internal/handlers"
	"github.com/miftahulmahfuzh/lunch-delivery/internal/models"
	"github.com/miftahulmahfuzh/lunch-delivery/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Config loading failed:", err)
	}

	// Database connection
	db, err := database.NewConnection(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	repo := models.NewRepository(db.DB)

	// Initialize nutritionist service
	nutritionistService, err := services.NewNutritionistService(cfg, repo)
	if err != nil {
		log.Fatal("Nutritionist service initialization failed:", err)
	}

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

	handlers.SetupRoutes(r, repo, nutritionistService)

	log.Println("Server starting on :8080")
	r.Run(":8080")
}
