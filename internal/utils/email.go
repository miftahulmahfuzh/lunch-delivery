package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strconv"
)

type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type EmailService struct {
	config EmailConfig
}

func NewEmailService() *EmailService {
	port, _ := strconv.Atoi(getEnvWithDefault("SMTP_PORT", "587"))

	config := EmailConfig{
		Host:     getEnvWithDefault("SMTP_HOST", "smtp.gmail.com"),
		Port:     port,
		Username: getEnvWithDefault("SMTP_USERNAME", ""),
		Password: getEnvWithDefault("SMTP_PASSWORD", ""),
		From:     getEnvWithDefault("SMTP_FROM", ""),
	}

	return &EmailService{config: config}
}

func (e *EmailService) SendPasswordResetEmail(toEmail, resetToken, baseURL string) error {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", baseURL, resetToken)

	subject := "Password Reset - Lunch Delivery System"
	body := fmt.Sprintf(`
Hello,

You have requested to reset your password for the Lunch Delivery System.

Please click the link below to reset your password:
%s

This link will expire in 1 hour.

If you did not request this password reset, please ignore this email.

Best regards,
Lunch Delivery Team
`, resetURL)

	return e.sendEmail(toEmail, subject, body)
}

func (e *EmailService) sendEmail(to, subject, body string) error {
	if e.config.Username == "" || e.config.Password == "" {
		log.Printf("SMTP credentials not configured, would send email to %s with subject: %s", to, subject)
		log.Printf("Email body:\n%s", body)
		return nil
	}

	from := e.config.From
	if from == "" {
		from = e.config.Username
	}

	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n", to, subject, body))

	auth := smtp.PlainAuth("", e.config.Username, e.config.Password, e.config.Host)
	addr := fmt.Sprintf("%s:%d", e.config.Host, e.config.Port)

	err := smtp.SendMail(addr, auth, from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Password reset email sent to %s", to)
	return nil
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
