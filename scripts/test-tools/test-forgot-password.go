package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/miftahulmahfuzh/lunch-delivery/internal/utils"
)

func main() {
	log.Println("🔐 Forgot Password Email Test")
	log.Println("==============================")

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables")
	}

	// Get test email address
	testEmail := os.Getenv("SMTP_TEST_EMAIL_ADDRESS")
	if testEmail == "" {
		log.Fatal("❌ SMTP_TEST_EMAIL_ADDRESS environment variable is required")
	}

	log.Printf("📤 Testing forgot password email to: %s", testEmail)

	// Generate a test token
	token, err := utils.GeneratePasswordResetToken()
	if err != nil {
		log.Fatalf("❌ Failed to generate token: %v", err)
	}

	log.Printf("🔑 Generated test token: %s", token[:50]+"...")

	// Create email service
	emailService := utils.NewEmailService()
	baseURL := "http://localhost:8080"

	// Send forgot password email
	log.Println("📧 Sending forgot password email...")

	start := time.Now()
	err = emailService.SendPasswordResetEmail(testEmail, token, baseURL)
	duration := time.Since(start)

	if err != nil {
		log.Printf("❌ Failed to send email: %v", err)
		log.Println("\n🔍 This is expected if SMTP credentials are not configured.")
		log.Println("📋 Check the console log above - the email content should be displayed.")
		log.Println("\n💡 To send real emails:")
		log.Println("   1. Set up Gmail App Password (see scripts/smtp/setup-gmail.md)")
		log.Println("   2. Or configure a different email provider in .env")
		os.Exit(1)
	}

	log.Printf("✅ Forgot password email sent successfully! (took %v)", duration)
	log.Printf("📬 Check your inbox at: %s", testEmail)
	log.Printf("🔗 Reset link: %s/reset-password?token=%s", baseURL, token)
	log.Println("\n🎉 Forgot password email functionality is working!")
}
