package testutils

import (
	"time"
)

// Common test data and utilities without package dependencies

// Test Environment Variables for Config Testing
func MockEnvironment() map[string]string {
	return map[string]string{
		"DB_HOST":                     "test-host",
		"DB_PORT":                     "5433",
		"DB_USER":                     "test-user",
		"DB_PASSWORD":                 "test-password",
		"DB_NAME":                     "test-database",
		"LLM_TYPE":                    "TEST_LLM",
		"DEEPSEEK_TENCENT_API_KEY":    "test-key-123",
		"DEEPSEEK_TENCENT_MODEL":      "test-model",
		"DEEPSEEK_TENCENT_BASE_URL":   "https://test-api.com",
		"SMTP_HOST":                   "smtp.test.com",
		"SMTP_PORT":                   "587",
		"SMTP_USERNAME":               "test@test.com",
		"SMTP_PASSWORD":               "test-smtp-pass",
		"SMTP_FROM":                   "noreply@test.com",
	}
}

// Common Test Dates
func TestDate() time.Time {
	return time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
}

func TestDateString() string {
	return "2024-01-15"
}

// Mock Configuration - returns a generic interface{} to avoid import cycles
func MockConfig() interface{} {
	return map[string]interface{}{
		"DBHost":                   "test-host",
		"DBPort":                   "5433",
		"DBUser":                   "test-user",
		"DBPassword":               "test-password",
		"DBName":                   "test-database",
		"LLMType":                  "TEST_LLM",
		"DeepSeekTencentAPIKey":    "test-key-123",
		"DeepSeekTencentModel":     "test-model",
		"DeepSeekTencentBaseURL":   "https://test-api.com",
		"SMTPHost":                 "smtp.test.com",
		"SMTPPort":                 "587",
		"SMTPUsername":             "test@test.com",
		"SMTPPassword":             "test-smtp-pass",
		"SMTPFrom":                 "noreply@test.com",
	}
}

// Mock Menu Items for testing
func MockMenuItems() []interface{} {
	return []interface{}{
		map[string]interface{}{"ID": 1, "Name": "Nasi Gudeg", "Price": 25000, "Active": true},
		map[string]interface{}{"ID": 2, "Name": "Ayam Bakar", "Price": 30000, "Active": true},
		map[string]interface{}{"ID": 3, "Name": "Sayur Lodeh", "Price": 15000, "Active": true},
		map[string]interface{}{"ID": 4, "Name": "Tempe Goreng", "Price": 10000, "Active": true},
		map[string]interface{}{"ID": 5, "Name": "Es Teh Manis", "Price": 5000, "Active": true},
	}
}

// Mock Company for testing
func MockCompany() interface{} {
	return map[string]interface{}{
		"ID":      1,
		"Name":    "Test Company",
		"Address": "123 Test St",
		"Contact": "test@company.com",
	}
}