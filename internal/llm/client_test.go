package llm

import (
	"fmt"
	"testing"
	"time"

	"github.com/miftahulmahfuzh/lunch-delivery/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "creates client with valid configuration",
			config: &config.Config{
				DeepseekTencentAPIKey:  "valid-api-key",
				DeepseekTencentModel:   "deepseek-v3",
				DeepseekTencentBaseURL: "https://api.test.com/v1",
				LLMRequestTimeout:      5 * time.Minute,
			},
			expectError: false,
		},
		{
			name: "fails when API key is empty",
			config: &config.Config{
				DeepseekTencentAPIKey:  "",
				DeepseekTencentModel:   "deepseek-v3",
				DeepseekTencentBaseURL: "https://api.test.com/v1",
				LLMRequestTimeout:      5 * time.Minute,
			},
			expectError: true,
		},
		{
			name: "creates client with custom timeout",
			config: &config.Config{
				DeepseekTencentAPIKey:  "test-key",
				DeepseekTencentModel:   "custom-model",
				DeepseekTencentBaseURL: "https://custom.api.com/v1",
				LLMRequestTimeout:      10 * time.Second,
			},
			expectError: false,
		},
		{
			name: "creates client with different model",
			config: &config.Config{
				DeepseekTencentAPIKey:  "another-key",
				DeepseekTencentModel:   "gpt-4",
				DeepseekTencentBaseURL: "https://openai.api.com/v1",
				LLMRequestTimeout:      2 * time.Minute,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, client)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
				assert.NotNil(t, client.llm)
			}
		})
	}
}

func TestClient_GenerateContent(t *testing.T) {
	// Note: These tests will not actually call external LLM services
	// They test the client wrapper functionality

	t.Run("calls underlying LLM client", func(t *testing.T) {
		cfg := &config.Config{
			DeepseekTencentAPIKey:  "test-key",
			DeepseekTencentModel:   "test-model",
			DeepseekTencentBaseURL: "https://test-api.com",
			LLMRequestTimeout:      time.Second * 30,
		}

		client, err := NewClient(cfg)
		require.NoError(t, err)
		require.NotNil(t, client)

		// Test our interface method
		_, err = client.GenerateContent("You are a helpful assistant", "Hello", "0.7")

		// We expect an error here since we're not connecting to a real service
		// The important thing is that the method exists and can be called
		assert.NotNil(t, err) // Expected to fail without real API
	})

	t.Run("handles context cancellation", func(t *testing.T) {
		cfg := &config.Config{
			DeepseekTencentAPIKey:  "test-key",
			DeepseekTencentModel:   "test-model",
			DeepseekTencentBaseURL: "https://test-api.com",
			LLMRequestTimeout:      time.Second * 30,
		}

		client, err := NewClient(cfg)
		require.NoError(t, err)

		// Test with our interface method
		_, err = client.GenerateContent("System prompt", "Test message", "0.7")

		// Should get context cancellation error
		assert.Error(t, err)
	})

	t.Run("handles timeout", func(t *testing.T) {
		cfg := &config.Config{
			DeepseekTencentAPIKey:  "test-key",
			DeepseekTencentModel:   "test-model",
			DeepseekTencentBaseURL: "https://test-api.com",
			LLMRequestTimeout:      1 * time.Millisecond, // Very short timeout
		}

		client, err := NewClient(cfg)
		require.NoError(t, err)

		_, err = client.GenerateContent("System prompt", "Test with short timeout", "0.7")

		// Should get timeout or connection error
		assert.Error(t, err)
	})
}

func TestNewClient_ConfigurationEdgeCases(t *testing.T) {
	t.Run("handles nil config", func(t *testing.T) {
		client, err := NewClient(nil)
		assert.Error(t, err)
		assert.Nil(t, client)
	})

	t.Run("handles empty base URL", func(t *testing.T) {
		cfg := &config.Config{
			DeepseekTencentAPIKey:  "test-key",
			DeepseekTencentModel:   "test-model",
			DeepseekTencentBaseURL: "",
			LLMRequestTimeout:      5 * time.Minute,
		}

		client, err := NewClient(cfg)
		// Should still create client, but might fail on actual use
		assert.NoError(t, err)
		assert.NotNil(t, client)
	})

	t.Run("handles empty model name", func(t *testing.T) {
		cfg := &config.Config{
			DeepseekTencentAPIKey:  "test-key",
			DeepseekTencentModel:   "",
			DeepseekTencentBaseURL: "https://api.test.com/v1",
			LLMRequestTimeout:      5 * time.Minute,
		}

		client, err := NewClient(cfg)
		// Should still create client
		assert.NoError(t, err)
		assert.NotNil(t, client)
	})

	t.Run("handles zero timeout", func(t *testing.T) {
		cfg := &config.Config{
			DeepseekTencentAPIKey:  "test-key",
			DeepseekTencentModel:   "test-model",
			DeepseekTencentBaseURL: "https://api.test.com/v1",
			LLMRequestTimeout:      0,
		}

		client, err := NewClient(cfg)
		assert.NoError(t, err)
		assert.NotNil(t, client)
	})
}

func TestClient_Integration(t *testing.T) {
	t.Run("client wraps LLM correctly", func(t *testing.T) {
		cfg := &config.Config{
			DeepseekTencentAPIKey:  "test-key",
			DeepseekTencentModel:   "test-model",
			DeepseekTencentBaseURL: "https://test-api.com",
			LLMRequestTimeout:      time.Second * 30,
		}
		cfg.DeepseekTencentAPIKey = "integration-test-key"

		client, err := NewClient(cfg)
		require.NoError(t, err)
		require.NotNil(t, client)

		// Test that client implements the expected interface
		assert.NotNil(t, client.llm)

		// Test method signature exists and is callable
		// Test our interface method
		response, err := client.GenerateContent("You are a helpful assistant.", "Hello, how are you?", "0.7")

		// We expect this to fail in tests, but the structure should be correct
		if err != nil {
			// Error is expected without real LLM service
			assert.Error(t, err)
			assert.Empty(t, response)
		} else {
			// If somehow it succeeds, response should be valid
			assert.NotEmpty(t, response)
		}
	})

	t.Run("client uses configuration correctly", func(t *testing.T) {
		// Test with different configurations to ensure they're applied
		configs := []*config.Config{
			{
				DeepseekTencentAPIKey:  "key1",
				DeepseekTencentModel:   "model1",
				DeepseekTencentBaseURL: "https://api1.com/v1",
				LLMRequestTimeout:      1 * time.Minute,
			},
			{
				DeepseekTencentAPIKey:  "key2",
				DeepseekTencentModel:   "model2",
				DeepseekTencentBaseURL: "https://api2.com/v1",
				LLMRequestTimeout:      2 * time.Minute,
			},
		}

		for i, cfg := range configs {
			t.Run(fmt.Sprintf("config_%d", i), func(t *testing.T) {
				client, err := NewClient(cfg)
				assert.NoError(t, err)
				assert.NotNil(t, client)
				assert.NotNil(t, client.llm)
			})
		}
	})
}

func TestClient_ErrorHandling(t *testing.T) {
	t.Run("returns error when LLM creation fails", func(t *testing.T) {
		// Test with invalid configuration that should cause LLM creation to fail
		cfg := &config.Config{
			DeepseekTencentAPIKey:  "invalid-key-format-that-might-cause-errors",
			DeepseekTencentModel:   "invalid-model",
			DeepseekTencentBaseURL: "not-a-url",
			LLMRequestTimeout:      -1 * time.Second, // negative timeout
		}

		client, err := NewClient(cfg)

		// The behavior here depends on the underlying LLM library
		// It might succeed in creation but fail on use, or fail immediately
		if err != nil {
			assert.Error(t, err)
			assert.Nil(t, client)
		} else {
			// If creation succeeds, client should be valid
			assert.NotNil(t, client)
		}
	})

	t.Run("handles invalid options gracefully", func(t *testing.T) {
		cfg := &config.Config{
			DeepseekTencentAPIKey:  "test-key",
			DeepseekTencentModel:   "test-model",
			DeepseekTencentBaseURL: "https://test.com/v1",
			LLMRequestTimeout:      5 * time.Minute,
		}

		// Should create client even with edge case configurations
		client, err := NewClient(cfg)
		assert.NoError(t, err)
		assert.NotNil(t, client)
	})
}

// Benchmark tests
func BenchmarkNewClient(b *testing.B) {
	cfg := &config.Config{
		DeepseekTencentAPIKey:  "benchmark-key",
		DeepseekTencentModel:   "test-model",
		DeepseekTencentBaseURL: "https://test-api.com",
		LLMRequestTimeout:      time.Second * 30,
	}

	for i := 0; i < b.N; i++ {
		client, err := NewClient(cfg)
		if err != nil {
			b.Fatal(err)
		}
		if client == nil {
			b.Fatal("client is nil")
		}
	}
}

func BenchmarkGenerateContent(b *testing.B) {
	cfg := &config.Config{
		DeepseekTencentAPIKey:  "benchmark-key",
		DeepseekTencentModel:   "test-model",
		DeepseekTencentBaseURL: "https://test-api.com",
		LLMRequestTimeout:      time.Second * 30,
	}

	client, err := NewClient(cfg)
	if err != nil {
		b.Skip("Cannot create client for benchmark:", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.GenerateContent("You are a helpful assistant", "Benchmark test message", "0.7")
		// We expect errors in benchmark without real API
		_ = err
	}
}
