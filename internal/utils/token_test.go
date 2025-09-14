package utils

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeneratePasswordResetToken(t *testing.T) {
	t.Run("generates valid token format", func(t *testing.T) {
		token, err := GeneratePasswordResetToken()

		require.NoError(t, err)
		assert.NotEmpty(t, token)

		// Token should have format: uuid-hex-timestamp
		parts := strings.Split(token, "-")
		assert.GreaterOrEqual(t, len(parts), 7) // UUID has 5 parts, plus hex, plus timestamp

		// Check UUID format (first 5 parts should form a valid UUID)
		uuidPart := strings.Join(parts[:5], "-")
		uuidPattern := `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`
		matched, err := regexp.MatchString(uuidPattern, uuidPart)
		require.NoError(t, err)
		assert.True(t, matched, "UUID part should match UUID format: %s", uuidPart)

		// Check hex part (6th part should be 32 hex characters)
		hexPart := parts[5]
		assert.Len(t, hexPart, 32, "Hex part should be 32 characters")
		hexPattern := `^[0-9a-f]{32}$`
		matched, err = regexp.MatchString(hexPattern, hexPart)
		require.NoError(t, err)
		assert.True(t, matched, "Hex part should be valid hex: %s", hexPart)

		// Check timestamp part (last part should be a valid timestamp)
		timestampPart := parts[len(parts)-1]
		assert.NotEmpty(t, timestampPart)
		assert.Regexp(t, `^\d+$`, timestampPart, "Timestamp part should be numeric")
	})

	t.Run("generates unique tokens", func(t *testing.T) {
		tokens := make(map[string]bool)
		numTokens := 1000

		for i := 0; i < numTokens; i++ {
			token, err := GeneratePasswordResetToken()
			require.NoError(t, err)

			// Check token is unique
			assert.False(t, tokens[token], "Token should be unique: %s", token)
			tokens[token] = true
		}

		assert.Len(t, tokens, numTokens, "All tokens should be unique")
	})

	t.Run("generates tokens with current timestamp", func(t *testing.T) {
		beforeTime := time.Now().Unix()

		token, err := GeneratePasswordResetToken()
		require.NoError(t, err)

		afterTime := time.Now().Unix()

		// Extract timestamp from token
		parts := strings.Split(token, "-")
		timestampStr := parts[len(parts)-1]

		// Convert to int64 for comparison
		var timestamp int64
		_, err = fmt.Sscanf(timestampStr, "%d", &timestamp)
		require.NoError(t, err)

		// Timestamp should be between before and after times
		assert.GreaterOrEqual(t, timestamp, beforeTime, "Token timestamp should be >= beforeTime")
		assert.LessOrEqual(t, timestamp, afterTime, "Token timestamp should be <= afterTime")
	})

	t.Run("token components are properly formatted", func(t *testing.T) {
		token, err := GeneratePasswordResetToken()
		require.NoError(t, err)

		parts := strings.Split(token, "-")

		// UUID component (first 5 parts)
		uuidPart := strings.Join(parts[:5], "-")
		assert.Len(t, uuidPart, 36, "UUID should be 36 characters including hyphens")

		// Hex component (6th part)
		hexPart := parts[5]
		assert.Len(t, hexPart, 32, "Hex part should be exactly 32 characters")

		// All parts should contain only valid characters
		for i, part := range parts[:6] { // UUID + hex parts
			if i < 5 { // UUID parts
				assert.Regexp(t, `^[0-9a-f-]+$`, part, "UUID part %d should contain only hex chars and hyphens", i)
			} else { // Hex part
				assert.Regexp(t, `^[0-9a-f]+$`, part, "Hex part should contain only hex characters")
			}
		}

		// Timestamp part should be numeric
		timestampPart := parts[len(parts)-1]
		assert.Regexp(t, `^\d+$`, timestampPart, "Timestamp should be numeric")
	})

	t.Run("multiple calls produce different tokens with different timestamps", func(t *testing.T) {
		token1, err := GeneratePasswordResetToken()
		require.NoError(t, err)

		// Small delay to ensure different timestamp
		time.Sleep(time.Millisecond * 10)

		token2, err := GeneratePasswordResetToken()
		require.NoError(t, err)

		assert.NotEqual(t, token1, token2, "Consecutive tokens should be different")

		// Extract timestamps
		parts1 := strings.Split(token1, "-")
		parts2 := strings.Split(token2, "-")

		timestamp1 := parts1[len(parts1)-1]
		timestamp2 := parts2[len(parts2)-1]

		// Timestamps should be different (or at least token2 >= token1)
		var ts1, ts2 int64
		if _, err := fmt.Sscanf(timestamp1, "%d", &ts1); err != nil {
			t.Errorf("Failed to parse timestamp1: %v", err)
		}
		if _, err := fmt.Sscanf(timestamp2, "%d", &ts2); err != nil {
			t.Errorf("Failed to parse timestamp2: %v", err)
		}

		assert.GreaterOrEqual(t, ts2, ts1, "Second token timestamp should be >= first token timestamp")
	})
}

func TestGeneratePasswordResetToken_EdgeCases(t *testing.T) {
	t.Run("handles rapid successive calls", func(t *testing.T) {
		tokens := make([]string, 100)

		// Generate many tokens rapidly
		for i := 0; i < 100; i++ {
			token, err := GeneratePasswordResetToken()
			require.NoError(t, err)
			tokens[i] = token
		}

		// All should be unique
		tokenMap := make(map[string]bool)
		for _, token := range tokens {
			assert.False(t, tokenMap[token], "Token should be unique: %s", token)
			tokenMap[token] = true
		}

		assert.Len(t, tokenMap, 100, "All rapid tokens should be unique")
	})

	t.Run("token length is consistent", func(t *testing.T) {
		var lengths []int

		for i := 0; i < 10; i++ {
			token, err := GeneratePasswordResetToken()
			require.NoError(t, err)
			lengths = append(lengths, len(token))
		}

		// All tokens should have the same length (within a reasonable range due to timestamp)
		minLength := lengths[0]
		maxLength := lengths[0]

		for _, length := range lengths {
			if length < minLength {
				minLength = length
			}
			if length > maxLength {
				maxLength = length
			}
		}

		// Length difference should be small (only timestamp can vary)
		assert.LessOrEqual(t, maxLength-minLength, 2, "Token lengths should be very similar")
	})

	t.Run("token contains no spaces or invalid characters", func(t *testing.T) {
		// Compile the regex once outside the loop for better performance
		validPattern := regexp.MustCompile(`^[0-9a-f-]+$`)

		for i := 0; i < 10; i++ {
			token, err := GeneratePasswordResetToken()
			require.NoError(t, err)

			assert.NotContains(t, token, " ", "Token should not contain spaces")
			assert.NotContains(t, token, "\n", "Token should not contain newlines")
			assert.NotContains(t, token, "\t", "Token should not contain tabs")

			// Should only contain valid characters
			matched := validPattern.MatchString(token)
			assert.True(t, matched, "Token should only contain hex characters and hyphens: %s", token)
		}
	})
}

func TestGeneratePasswordResetToken_Security(t *testing.T) {
	t.Run("tokens have sufficient entropy", func(t *testing.T) {
		const numTokens = 1000
		tokens := make([]string, numTokens)

		for i := 0; i < numTokens; i++ {
			token, err := GeneratePasswordResetToken()
			require.NoError(t, err)
			tokens[i] = token
		}

		// Check for collisions (should be extremely rare)
		tokenSet := make(map[string]bool)
		collisions := 0

		for _, token := range tokens {
			if tokenSet[token] {
				collisions++
			}
			tokenSet[token] = true
		}

		assert.Equal(t, 0, collisions, "Should have no collisions in %d tokens", numTokens)
		assert.Len(t, tokenSet, numTokens, "All tokens should be unique")
	})

	t.Run("tokens are not predictable", func(t *testing.T) {
		// Generate several tokens and check they don't follow a predictable pattern
		tokens := make([]string, 5)
		for i := range tokens {
			token, err := GeneratePasswordResetToken()
			require.NoError(t, err)
			tokens[i] = token
		}

		// Check that consecutive tokens don't have predictable differences
		for i := 1; i < len(tokens); i++ {
			parts1 := strings.Split(tokens[i-1], "-")
			parts2 := strings.Split(tokens[i], "-")

			// UUID parts should be completely different (first 5 parts)
			for j := 0; j < 5; j++ {
				assert.NotEqual(t, parts1[j], parts2[j], "UUID parts should be different between tokens")
			}

			// Hex parts should be different
			assert.NotEqual(t, parts1[5], parts2[5], "Hex parts should be different between tokens")

			// Only timestamps might be the same (if generated in same second)
		}
	})

	t.Run("tokens have appropriate length for security", func(t *testing.T) {
		token, err := GeneratePasswordResetToken()
		require.NoError(t, err)

		// Remove hyphens to count actual entropy characters
		entropyChars := strings.ReplaceAll(token, "-", "")

		// Should have at least 64 characters of entropy (UUID=32, hex=32, timestamp varies)
		assert.GreaterOrEqual(t, len(entropyChars), 64, "Token should have sufficient entropy characters")
	})
}

func TestGeneratePasswordResetToken_Performance(t *testing.T) {
	t.Run("generation is reasonably fast", func(t *testing.T) {
		start := time.Now()

		for i := 0; i < 100; i++ {
			_, err := GeneratePasswordResetToken()
			require.NoError(t, err)
		}

		duration := time.Since(start)

		// Should be able to generate 100 tokens in less than 100ms
		assert.Less(t, duration, 100*time.Millisecond, "Token generation should be fast")
	})
}

// Benchmark tests
func BenchmarkGeneratePasswordResetToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GeneratePasswordResetToken()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGeneratePasswordResetToken_Parallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := GeneratePasswordResetToken()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

