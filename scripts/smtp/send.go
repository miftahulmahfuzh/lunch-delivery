package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func loadSMTPConfig() (*SMTPConfig, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables")
	}

	host := getEnvWithDefault("SMTP_HOST", "")
	portStr := getEnvWithDefault("SMTP_PORT", "587")
	username := getEnvWithDefault("SMTP_USERNAME", "")
	password := getEnvWithDefault("SMTP_PASSWORD", "")
	from := getEnvWithDefault("SMTP_FROM", username)

	if host == "" {
		return nil, fmt.Errorf("SMTP_HOST is required")
	}
	if username == "" {
		return nil, fmt.Errorf("SMTP_USERNAME is required")
	}
	if password == "" {
		return nil, fmt.Errorf("SMTP_PASSWORD is required")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid SMTP_PORT: %s", portStr)
	}

	return &SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
	}, nil
}

func sendTestEmail(config *SMTPConfig, toEmail, subject, body string) error {
	// Create message
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s\r\n",
		config.From, toEmail, subject, body)

	// Setup authentication
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)

	// Setup server address
	serverAddr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	// Handle TLS connection for common SMTP ports
	if config.Port == 465 {
		return sendWithTLS(config, serverAddr, auth, toEmail, []byte(message))
	} else {
		return sendWithSTARTTLS(config, serverAddr, auth, toEmail, []byte(message))
	}
}

func sendWithTLS(config *SMTPConfig, serverAddr string, auth smtp.Auth, toEmail string, message []byte) error {
	// Connect with TLS
	tlsconfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         config.Host,
	}

	conn, err := tls.Dial("tcp", serverAddr, tlsconfig)
	if err != nil {
		return fmt.Errorf("TLS dial failed: %w", err)
	}
	defer func() { _ = conn.Close() }()

	client, err := smtp.NewClient(conn, config.Host)
	if err != nil {
		return fmt.Errorf("SMTP client creation failed: %w", err)
	}
	defer func() { _ = client.Close() }()

	if auth != nil {
		if ok, _ := client.Extension("AUTH"); !ok {
			return fmt.Errorf("smtp: server doesn't support AUTH")
		}
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}
	}

	if err = client.Mail(config.From); err != nil {
		return fmt.Errorf("MAIL command failed: %w", err)
	}

	if err = client.Rcpt(toEmail); err != nil {
		return fmt.Errorf("RCPT command failed: %w", err)
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("DATA command failed: %w", err)
	}
	defer func() { _ = writer.Close() }()

	_, err = writer.Write(message)
	if err != nil {
		return fmt.Errorf("writing message failed: %w", err)
	}

	return nil
}

func sendWithSTARTTLS(config *SMTPConfig, serverAddr string, auth smtp.Auth, toEmail string, message []byte) error {
	// Use standard smtp.SendMail for STARTTLS (ports 587, 25)
	return smtp.SendMail(serverAddr, auth, config.From, []string{toEmail}, message)
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	log.Println("🚀 SMTP Test Email Sender")
	log.Println("==========================")

	// Load SMTP configuration
	config, err := loadSMTPConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load SMTP configuration: %v", err)
	}

	log.Printf("📧 SMTP Configuration loaded:")
	log.Printf("   Host: %s:%d", config.Host, config.Port)
	log.Printf("   Username: %s", config.Username)
	log.Printf("   From: %s", config.From)

	// Get target email address
	testEmail := os.Getenv("SMTP_TEST_EMAIL_ADDRESS")
	if testEmail == "" {
		log.Fatal("❌ SMTP_TEST_EMAIL_ADDRESS environment variable is required")
	}

	log.Printf("📤 Sending test email to: %s", testEmail)

	// Create test email content
	subject := "🧪 SMTP Test Email - Lunch Delivery System"
	body := fmt.Sprintf(`Hello!

This is a test email sent from the Lunch Delivery System SMTP configuration.

Test Details:
- Timestamp: %s
- SMTP Server: %s:%d
- From Address: %s
- Test Target: %s

If you receive this email, your SMTP configuration is working correctly! ✅

This email was sent as part of testing the forgot password functionality.

Best regards,
Lunch Delivery System Test Suite
`, time.Now().Format("2006-01-02 15:04:05 MST"), config.Host, config.Port, config.From, testEmail)

	// Send the email
	log.Println("🔄 Sending email...")

	start := time.Now()
	err = sendTestEmail(config, testEmail, subject, body)
	duration := time.Since(start)

	if err != nil {
		log.Printf("❌ Failed to send email: %v", err)
		log.Println("\n🔍 Troubleshooting tips:")
		log.Println("   1. Check your SMTP credentials (username/password)")
		log.Println("   2. Verify SMTP server settings (host/port)")
		log.Println("   3. Enable 'Less secure app access' or use App Passwords for Gmail")
		log.Println("   4. Check firewall settings")
		log.Println("   5. Verify the recipient email address")

		// Print configuration details (without password)
		log.Println("\n📋 Current Configuration:")
		log.Printf("   SMTP_HOST=%s", config.Host)
		log.Printf("   SMTP_PORT=%d", config.Port)
		log.Printf("   SMTP_USERNAME=%s", config.Username)
		log.Printf("   SMTP_FROM=%s", config.From)
		log.Printf("   SMTP_PASSWORD=%s", strings.Repeat("*", len(config.Password)))

		os.Exit(1)
	}

	log.Printf("✅ Email sent successfully! (took %v)", duration)
	log.Printf("📬 Check your inbox at: %s", testEmail)
	log.Println("\n🎉 SMTP configuration is working correctly!")
}
