package config

import (
	"os"
	"testing"
	"time"

	"github.com/miftahulmahfuzh/lunch-delivery/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected *Config
	}{
		{
			name:    "default values when no environment variables are set",
			envVars: map[string]string{},
			expected: &Config{
				DBHost:                 "localhost",
				DBPort:                 "5432",
				DBUser:                 "lunch_user",
				DBPassword:             "1234",
				DBName:                 "lunch_delivery",
				LLMType:                "DEEPSEEK_TENCENT",
				DeepseekTencentAPIKey:  "",
				DeepseekTencentModel:   "deepseek-v3",
				DeepseekTencentBaseURL: "https://api.lkeap.tencentcloud.com/v1",
				LLMRequestTimeout:      5 * time.Minute,
			},
		},
		{
			name: "uses environment variables when set",
			envVars: map[string]string{
				"DB_HOST":                   "custom-host",
				"DB_PORT":                   "3306",
				"DB_USER":                   "custom_user",
				"DB_PASSWORD":               "custom_pass",
				"DB_NAME":                   "custom_db",
				"LLM_TYPE":                  "CUSTOM_LLM",
				"DEEPSEEK_TENCENT_API_KEY":  "test-api-key",
				"DEEPSEEK_TENCENT_MODEL":    "custom-model",
				"DEEPSEEK_TENCENT_BASE_URL": "https://custom-api.com/v1",
			},
			expected: &Config{
				DBHost:                 "custom-host",
				DBPort:                 "3306",
				DBUser:                 "custom_user",
				DBPassword:             "custom_pass",
				DBName:                 "custom_db",
				LLMType:                "CUSTOM_LLM",
				DeepseekTencentAPIKey:  "test-api-key",
				DeepseekTencentModel:   "custom-model",
				DeepseekTencentBaseURL: "https://custom-api.com/v1",
				LLMRequestTimeout:      5 * time.Minute,
			},
		},
		{
			name: "partial environment variables with defaults",
			envVars: map[string]string{
				"DB_HOST":                  "partial-host",
				"DEEPSEEK_TENCENT_API_KEY": "partial-key",
			},
			expected: &Config{
				DBHost:                 "partial-host",
				DBPort:                 "5432",
				DBUser:                 "lunch_user",
				DBPassword:             "1234",
				DBName:                 "lunch_delivery",
				LLMType:                "DEEPSEEK_TENCENT",
				DeepseekTencentAPIKey:  "partial-key",
				DeepseekTencentModel:   "deepseek-v3",
				DeepseekTencentBaseURL: "https://api.lkeap.tencentcloud.com/v1",
				LLMRequestTimeout:      5 * time.Minute,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment
			cleanup := testutils.SetTestEnv(tt.envVars)
			defer cleanup()

			// Test Load function
			cfg, err := Load()

			// Assertions
			require.NoError(t, err)
			assert.Equal(t, tt.expected.DBHost, cfg.DBHost)
			assert.Equal(t, tt.expected.DBPort, cfg.DBPort)
			assert.Equal(t, tt.expected.DBUser, cfg.DBUser)
			assert.Equal(t, tt.expected.DBPassword, cfg.DBPassword)
			assert.Equal(t, tt.expected.DBName, cfg.DBName)
			assert.Equal(t, tt.expected.LLMType, cfg.LLMType)
			assert.Equal(t, tt.expected.DeepseekTencentAPIKey, cfg.DeepseekTencentAPIKey)
			assert.Equal(t, tt.expected.DeepseekTencentModel, cfg.DeepseekTencentModel)
			assert.Equal(t, tt.expected.DeepseekTencentBaseURL, cfg.DeepseekTencentBaseURL)
			assert.Equal(t, tt.expected.LLMRequestTimeout, cfg.LLMRequestTimeout)
		})
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
		setEnv       bool
	}{
		{
			name:         "returns default when environment variable is not set",
			key:          "TEST_VAR_NOT_SET",
			defaultValue: "default_value",
			expected:     "default_value",
			setEnv:       false,
		},
		{
			name:         "returns environment value when set",
			key:          "TEST_VAR_SET",
			defaultValue: "default_value",
			envValue:     "env_value",
			expected:     "env_value",
			setEnv:       true,
		},
		{
			name:         "returns default when environment variable is empty string",
			key:          "TEST_VAR_EMPTY",
			defaultValue: "default_value",
			envValue:     "",
			expected:     "default_value",
			setEnv:       true,
		},
		{
			name:         "handles empty default value",
			key:          "TEST_VAR_NO_DEFAULT",
			defaultValue: "",
			envValue:     "some_value",
			expected:     "some_value",
			setEnv:       true,
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

			// Test getEnv function
			result := getEnv(tt.key, tt.defaultValue)

			// Assertion
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLoadWithDotEnvFile(t *testing.T) {
	// Create temporary directory
	tempDir, cleanup := testutils.CreateTempDir(t)
	defer cleanup()

	// Change to temp directory
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Errorf("Failed to change back to original directory: %v", err)
		}
	}()

	err = os.Chdir(tempDir)
	require.NoError(t, err)

	t.Run("loads .env file when it exists", func(t *testing.T) {
		// Create .env file
		envContent := `DB_HOST=dotenv-host
DB_PORT=9999
DB_USER=dotenv_user
DEEPSEEK_TENCENT_API_KEY=dotenv-api-key
`
		err := os.WriteFile(".env", []byte(envContent), 0644)
		require.NoError(t, err)

		// Test Load function
		cfg, err := Load()
		require.NoError(t, err)

		// Assertions - should use values from .env file
		assert.Equal(t, "dotenv-host", cfg.DBHost)
		assert.Equal(t, "9999", cfg.DBPort)
		assert.Equal(t, "dotenv_user", cfg.DBUser)
		assert.Equal(t, "dotenv-api-key", cfg.DeepseekTencentAPIKey)

		// Cleanup
		_ = os.Remove(".env")
	})

	t.Run("handles missing .env file gracefully", func(t *testing.T) {
		// Ensure no .env file exists and clean environment
		_ = os.Remove(".env")
		cleanup := testutils.SetTestEnv(map[string]string{
			"DB_HOST":     "",
			"DB_PORT":     "",
			"DB_USER":     "",
			"DB_PASSWORD": "",
			"DB_NAME":     "",
		})
		defer cleanup()

		// Test Load function
		cfg, err := Load()
		require.NoError(t, err)

		// Assertions - should use default values
		assert.Equal(t, "localhost", cfg.DBHost)
		assert.Equal(t, "5432", cfg.DBPort)
		assert.Equal(t, "lunch_user", cfg.DBUser)
	})

	t.Run("environment variables override .env file", func(t *testing.T) {
		// Create .env file
		envContent := `DB_HOST=dotenv-host
DB_PORT=9999
`
		err := os.WriteFile(".env", []byte(envContent), 0644)
		require.NoError(t, err)

		// Set environment variables
		cleanup := testutils.SetTestEnv(map[string]string{
			"DB_HOST": "env-override",
			"DB_USER": "env-user",
		})
		defer cleanup()

		// Test Load function
		cfg, err := Load()
		require.NoError(t, err)

		// Assertions - env vars should override .env file
		assert.Equal(t, "env-override", cfg.DBHost) // from env var
		assert.Equal(t, "9999", cfg.DBPort)         // from .env file
		assert.Equal(t, "env-user", cfg.DBUser)     // from env var

		// Cleanup
		_ = os.Remove(".env")
	})
}

func TestConfig_ValidationScenarios(t *testing.T) {
	t.Run("config with all fields populated", func(t *testing.T) {
		envVars := testutils.MockEnvironment()
		cleanup := testutils.SetTestEnv(envVars)
		defer cleanup()

		cfg, err := Load()
		require.NoError(t, err)

		// Validate all fields are populated as expected
		assert.NotEmpty(t, cfg.DBHost)
		assert.NotEmpty(t, cfg.DBPort)
		assert.NotEmpty(t, cfg.DBUser)
		assert.NotEmpty(t, cfg.DBPassword)
		assert.NotEmpty(t, cfg.DBName)
		assert.NotEmpty(t, cfg.LLMType)
		assert.NotEmpty(t, cfg.DeepseekTencentModel)
		assert.NotEmpty(t, cfg.DeepseekTencentBaseURL)
		assert.True(t, cfg.LLMRequestTimeout > 0)
	})

	t.Run("config with minimal required fields", func(t *testing.T) {
		// Remove any .env file
		_ = os.Remove(".env")

		// Clear all env vars explicitly
		cleanup := testutils.SetTestEnv(map[string]string{
			"DB_HOST":                   "",
			"DB_PORT":                   "",
			"DB_USER":                   "",
			"DB_PASSWORD":               "",
			"DB_NAME":                   "",
			"LLM_TYPE":                  "",
			"DEEPSEEK_TENCENT_API_KEY":  "",
			"DEEPSEEK_TENCENT_MODEL":    "",
			"DEEPSEEK_TENCENT_BASE_URL": "",
			"SMTP_HOST":                 "",
			"SMTP_PORT":                 "",
			"SMTP_USERNAME":             "",
			"SMTP_PASSWORD":             "",
			"SMTP_FROM":                 "",
		})
		defer cleanup()

		cfg, err := Load()
		require.NoError(t, err)

		// Should still create valid config with defaults
		assert.NotNil(t, cfg)
		assert.Equal(t, "localhost", cfg.DBHost)
		assert.Equal(t, "", cfg.DeepseekTencentAPIKey) // This can be empty
	})
}

// Benchmark tests for config loading performance
func BenchmarkLoad(b *testing.B) {
	cleanup := testutils.SetTestEnv(testutils.MockEnvironment())
	defer cleanup()

	for i := 0; i < b.N; i++ {
		_, err := Load()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGetEnv(b *testing.B) {
	_ = os.Setenv("BENCH_TEST_VAR", "bench_value")
	defer func() { _ = os.Unsetenv("BENCH_TEST_VAR") }()

	for i := 0; i < b.N; i++ {
		_ = getEnv("BENCH_TEST_VAR", "default")
	}
}
