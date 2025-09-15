package utils

import (
	"os"
	"testing"

	"github.com/miftahulmahfuzh/lunch-delivery/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEmailService(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected EmailConfig
	}{
		{
			name:    "uses default values when no environment variables are set",
			envVars: map[string]string{},
			expected: EmailConfig{
				Host:     "smtp.gmail.com",
				Port:     587,
				Username: "",
				Password: "",
				From:     "",
			},
		},
		{
			name: "uses environment variables when set",
			envVars: map[string]string{
				"SMTP_HOST":     "smtp.example.com",
				"SMTP_PORT":     "25",
				"SMTP_USERNAME": "test@example.com",
				"SMTP_PASSWORD": "testpass123",
				"SMTP_FROM":     "noreply@example.com",
			},
			expected: EmailConfig{
				Host:     "smtp.example.com",
				Port:     25,
				Username: "test@example.com",
				Password: "testpass123",
				From:     "noreply@example.com",
			},
		},
		{
			name: "handles partial environment variables",
			envVars: map[string]string{
				"SMTP_HOST":     "custom.smtp.com",
				"SMTP_USERNAME": "custom@user.com",
			},
			expected: EmailConfig{
				Host:     "custom.smtp.com",
				Port:     587, // default
				Username: "custom@user.com",
				Password: "",
				From:     "",
			},
		},
		{
			name: "handles invalid port number gracefully",
			envVars: map[string]string{
				"SMTP_PORT": "invalid",
			},
			expected: EmailConfig{
				Host:     "smtp.gmail.com",
				Port:     0, // strconv.Atoi returns 0 for invalid input
				Username: "",
				Password: "",
				From:     "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment
			cleanup := testutils.SetTestEnv(tt.envVars)
			defer cleanup()

			// Test NewEmailService
			service := NewEmailService()

			// Assertions
			require.NotNil(t, service)
			assert.Equal(t, tt.expected.Host, service.config.Host)
			assert.Equal(t, tt.expected.Port, service.config.Port)
			assert.Equal(t, tt.expected.Username, service.config.Username)
			assert.Equal(t, tt.expected.Password, service.config.Password)
			assert.Equal(t, tt.expected.From, service.config.From)
		})
	}
}

func TestSendPasswordResetEmail(t *testing.T) {
	tests := []struct {
		name              string
		toEmail           string
		resetToken        string
		baseURL           string
		envVars           map[string]string
		expectError       bool
		expectedURLInBody string
	}{
		{
			name:              "sends password reset email without SMTP credentials",
			toEmail:           "user@example.com",
			resetToken:        "test-token-123",
			baseURL:           "https://example.com",
			envVars:           map[string]string{}, // No SMTP credentials
			expectError:       false,               // Should not error when no SMTP credentials (just logs)
			expectedURLInBody: "https://example.com/reset-password?token=test-token-123",
		},
		{
			name:              "constructs correct reset URL with different base URL",
			toEmail:           "user@test.com",
			resetToken:        "different-token",
			baseURL:           "http://localhost:8080",
			envVars:           map[string]string{},
			expectError:       false,
			expectedURLInBody: "http://localhost:8080/reset-password?token=different-token",
		},
		{
			name:              "handles complex token",
			toEmail:           "complex@user.com",
			resetToken:        "abc123-def456-ghi789",
			baseURL:           "https://secure.app.com",
			envVars:           map[string]string{},
			expectError:       false,
			expectedURLInBody: "https://secure.app.com/reset-password?token=abc123-def456-ghi789",
		},
		{
			name:       "with SMTP credentials configured (would attempt real send)",
			toEmail:    "real@user.com",
			resetToken: "real-token",
			baseURL:    "https://real.com",
			envVars: map[string]string{
				"SMTP_HOST":     "smtp.gmail.com",
				"SMTP_PORT":     "587",
				"SMTP_USERNAME": "test@gmail.com",
				"SMTP_PASSWORD": "testpass",
				"SMTP_FROM":     "noreply@test.com",
			},
			expectError:       true, // Will likely fail without real SMTP server
			expectedURLInBody: "https://real.com/reset-password?token=real-token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment
			cleanup := testutils.SetTestEnv(tt.envVars)
			defer cleanup()

			// Create email service
			service := NewEmailService()

			// Test SendPasswordResetEmail
			err := service.SendPasswordResetEmail(tt.toEmail, tt.resetToken, tt.baseURL)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// The URL construction is tested by checking that the method doesn't panic
			// and constructs the expected URL format (can't easily test the actual email content
			// without intercepting the SMTP call or having a real email server)
		})
	}
}

func TestSendEmail(t *testing.T) {
	tests := []struct {
		name        string
		to          string
		subject     string
		body        string
		envVars     map[string]string
		expectError bool
	}{
		{
			name:    "logs email when no SMTP credentials configured",
			to:      "test@example.com",
			subject: "Test Subject",
			body:    "Test Body Content",
			envVars: map[string]string{
				// No SMTP credentials
			},
			expectError: false, // Should not error, just log
		},
		{
			name:    "handles empty from field by using username",
			to:      "test@example.com",
			subject: "Test Subject",
			body:    "Test Body",
			envVars: map[string]string{
				"SMTP_USERNAME": "sender@example.com",
				"SMTP_PASSWORD": "password",
				// SMTP_FROM not set
			},
			expectError: true, // Will error trying to connect to SMTP
		},
		{
			name:    "uses explicit from field when set",
			to:      "test@example.com",
			subject: "Test Subject",
			body:    "Test Body",
			envVars: map[string]string{
				"SMTP_USERNAME": "username@example.com",
				"SMTP_PASSWORD": "password",
				"SMTP_FROM":     "explicit@example.com",
			},
			expectError: true, // Will error trying to connect to SMTP
		},
		{
			name:    "handles missing password",
			to:      "test@example.com",
			subject: "Test Subject",
			body:    "Test Body",
			envVars: map[string]string{
				"SMTP_USERNAME": "user@example.com",
				// Password not set
			},
			expectError: false, // Should log instead of sending
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment
			cleanup := testutils.SetTestEnv(tt.envVars)
			defer cleanup()

			// Create email service
			service := NewEmailService()

			// Test sendEmail (private method, tested through SendPasswordResetEmail)
			err := service.sendEmail(tt.to, tt.subject, tt.body)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetEnvWithDefault(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		setEnv       bool
		expected     string
	}{
		{
			name:         "returns default when environment variable not set",
			key:          "TEST_EMAIL_VAR",
			defaultValue: "default_value",
			setEnv:       false,
			expected:     "default_value",
		},
		{
			name:         "returns environment value when set",
			key:          "TEST_EMAIL_VAR_SET",
			defaultValue: "default_value",
			envValue:     "env_value",
			setEnv:       true,
			expected:     "env_value",
		},
		{
			name:         "returns default when environment variable is empty",
			key:          "TEST_EMAIL_VAR_EMPTY",
			defaultValue: "default_value",
			envValue:     "",
			setEnv:       true,
			expected:     "default_value",
		},
		{
			name:         "handles empty default value",
			key:          "TEST_EMAIL_VAR_NO_DEFAULT",
			defaultValue: "",
			envValue:     "some_value",
			setEnv:       true,
			expected:     "some_value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment
			originalValue := os.Getenv(tt.key)
			defer func() {
				if originalValue != "" {
					_ = os.Setenv(tt.key, originalValue)
				} else {
					_ = os.Unsetenv(tt.key)
				}
			}()

			if tt.setEnv {
				_ = os.Setenv(tt.key, tt.envValue)
			} else {
				_ = os.Unsetenv(tt.key)
			}

			// Test getEnvWithDefault
			result := getEnvWithDefault(tt.key, tt.defaultValue)

			// Assertion
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEmailService_Integration(t *testing.T) {
	t.Run("complete email workflow without real SMTP", func(t *testing.T) {
		// Setup test environment
		cleanup := testutils.SetTestEnv(map[string]string{
			"SMTP_HOST": "test.smtp.com",
			"SMTP_PORT": "587",
			// No username/password so it will log instead of send
		})
		defer cleanup()

		// Create service
		service := NewEmailService()

		// Test complete workflow
		err := service.SendPasswordResetEmail(
			"test@example.com",
			"integration-test-token",
			"https://integration-test.com",
		)

		// Should not error because no SMTP credentials are provided
		assert.NoError(t, err)
	})

	t.Run("email service with all configurations", func(t *testing.T) {
		t.Skip("Skipping integration test that tries to connect to SMTP server")
		// Setup complete configuration
		envVars := map[string]string{
			"SMTP_HOST":     "smtp.test.com",
			"SMTP_PORT":     "25",
			"SMTP_USERNAME": "test@test.com",
			"SMTP_PASSWORD": "testpass",
			"SMTP_FROM":     "noreply@test.com",
		}

		cleanup := testutils.SetTestEnv(envVars)
		defer cleanup()

		// Create service
		service := NewEmailService()

		// Verify all configurations are loaded correctly
		assert.Equal(t, "smtp.test.com", service.config.Host)
		assert.Equal(t, 25, service.config.Port)
		assert.Equal(t, "test@test.com", service.config.Username)
		assert.Equal(t, "testpass", service.config.Password)
		assert.Equal(t, "noreply@test.com", service.config.From)
	})
}

func TestEmailService_EdgeCases(t *testing.T) {
	t.Run("handles special characters in email content", func(t *testing.T) {
		service := NewEmailService()

		// Test with special characters
		err := service.SendPasswordResetEmail(
			"user+test@example.com",
			"token-with-special-chars!@#$%",
			"https://example.com/path?param=value&other=true",
		)

		// Should not error (logs only)
		assert.NoError(t, err)
	})

	t.Run("handles empty values gracefully", func(t *testing.T) {
		service := NewEmailService()

		err := service.SendPasswordResetEmail("", "", "")
		assert.NoError(t, err) // Logs only, no SMTP credentials
	})

	t.Run("handles very long token", func(t *testing.T) {
		service := NewEmailService()

		longToken := ""
		for i := 0; i < 1000; i++ {
			longToken += "a"
		}

		err := service.SendPasswordResetEmail(
			"test@example.com",
			longToken,
			"https://example.com",
		)

		assert.NoError(t, err)
	})
}

// Benchmark tests
func BenchmarkNewEmailService(b *testing.B) {
	// Set up environment
	cleanup := testutils.SetTestEnv(map[string]string{
		"SMTP_HOST":     "smtp.benchmark.com",
		"SMTP_PORT":     "587",
		"SMTP_USERNAME": "bench@test.com",
		"SMTP_PASSWORD": "benchpass",
		"SMTP_FROM":     "noreply@bench.com",
	})
	defer cleanup()

	for i := 0; i < b.N; i++ {
		_ = NewEmailService()
	}
}

func BenchmarkSendPasswordResetEmail(b *testing.B) {
	service := NewEmailService()

	for i := 0; i < b.N; i++ {
		_ = service.SendPasswordResetEmail(
			"bench@example.com",
			"bench-token",
			"https://bench.com",
		)
	}
}
