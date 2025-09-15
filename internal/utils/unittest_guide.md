# Utils Package Unit Tests Guide

## Overview

This guide explains the unit tests for the `internal/utils` package, which provides utility functions for the lunch delivery system. The package contains essential utility functions for email services and token generation.

## Test Structure

### Files Tested
- `email.go` - Email service functionality with SMTP configuration
- `token.go` - Password reset token generation utilities

### Test Files
- `email_test.go` - Email service functionality testing
- `token_test.go` - Token generation functionality testing

## Test Categories

### 1. Email Service Tests (`email_test.go`)

#### A. Email Service Creation (`TestNewEmailService`)
Tests the email service initialization with environment variable configuration:

**Test Scenarios**:
- Uses default values when no environment variables are set
- Uses environment variables when set (SMTP_HOST, SMTP_PORT, SMTP_USERNAME, SMTP_PASSWORD, SMTP_FROM)
- Handles partial environment variables (falls back to defaults)
- Handles invalid port number gracefully (strconv.Atoi error handling)

**Environment Variables Tested**:
- `SMTP_HOST` - SMTP server host (default: smtp.gmail.com)
- `SMTP_PORT` - SMTP server port (default: 587)
- `SMTP_USERNAME` - SMTP authentication username
- `SMTP_PASSWORD` - SMTP authentication password
- `SMTP_FROM` - Email sender address

#### B. Password Reset Email (`TestSendPasswordResetEmail`)
Tests the password reset email functionality:

**Test Scenarios**:
- Sends password reset email without SMTP credentials (logs only)
- Constructs correct reset URL with different base URLs
- Handles complex tokens with special characters
- Attempts real SMTP connection when credentials are configured (expects error without real server)

**URL Construction Testing**:
- Base URL: `https://example.com` + Token: `test-token-123` → `https://example.com/reset-password?token=test-token-123`
- Different base URLs and token formats
- Special characters in tokens and URLs

#### C. Email Sending (`TestSendEmail`)
Tests the internal email sending functionality:

**Test Scenarios**:
- Logs email when no SMTP credentials configured (graceful fallback)
- Handles empty from field by using username as fallback
- Uses explicit from field when set
- Handles missing password (falls back to logging)
- Tests SMTP connection attempts (expects errors without real server)

#### D. Environment Helper (`TestGetEnvWithDefault`)
Tests the environment variable helper function:

**Test Scenarios**:
- Returns default when environment variable not set
- Returns environment value when set
- Returns default when environment variable is empty
- Handles empty default value

#### E. Integration Tests (`TestEmailService_Integration`)
Tests complete email workflow:

**Test Scenarios**:
- Complete email workflow without real SMTP (logs only)
- Email service with all configurations (skipped test for real SMTP)

#### F. Edge Cases (`TestEmailService_EdgeCases`)
Tests boundary conditions and special cases:

**Test Scenarios**:
- Handles special characters in email content
- Handles empty values gracefully
- Handles very long tokens (1000 character test)

#### G. Performance Tests (`BenchmarkNewEmailService`, `BenchmarkSendPasswordResetEmail`)
Benchmark tests for performance monitoring:
- Email service creation performance
- Password reset email sending performance

### 2. Token Generation Tests (`token_test.go`)

#### A. Token Format Validation (`TestGeneratePasswordResetToken`)
Tests the password reset token generation format:

**Token Structure**: `uuid-hex-timestamp`
- UUID part: Standard UUID format (36 characters with hyphens)
- Hex part: 32 hexadecimal characters
- Timestamp part: Unix timestamp

**Test Scenarios**:
- Generates valid token format
- Generates unique tokens (1000 token uniqueness test)
- Generates tokens with current timestamp
- Token components are properly formatted
- Multiple calls produce different tokens with different timestamps

#### B. Edge Cases (`TestGeneratePasswordResetToken_EdgeCases`)
Tests boundary conditions:

**Test Scenarios**:
- Handles rapid successive calls (100 tokens)
- Token length is consistent
- Token contains no spaces or invalid characters
- Only contains valid hex characters and hyphens

#### C. Security Tests (`TestGeneratePasswordResetToken_Security`)
Tests security aspects of token generation:

**Test Scenarios**:
- Tokens have sufficient entropy (1000 token collision test)
- Tokens are not predictable (consecutive token comparison)
- Tokens have appropriate length for security (64+ entropy characters)

#### D. Performance Tests (`TestGeneratePasswordResetToken_Performance`)
Tests generation performance:
- Generation speed test (100 tokens in <100ms)
- Benchmark tests (sequential and parallel)

## Running the Tests

### Run All Utils Tests
```bash
go test ./internal/utils
```

### Run Specific Test Categories
```bash
# Email service tests only
go test ./internal/utils -run TestNewEmailService
go test ./internal/utils -run TestSendPasswordResetEmail

# Token generation tests only
go test ./internal/utils -run TestGeneratePasswordResetToken

# Edge case tests
go test ./internal/utils -run EdgeCases

# Security tests
go test ./internal/utils -run Security
```

### Run Tests with Verbose Output
```bash
go test ./internal/utils -v
```

### Run Tests with Coverage
```bash
go test ./internal/utils -cover
```

### Run Benchmark Tests
```bash
go test ./internal/utils -bench=.
```

## Test Patterns and Utilities Used

### From `testutils` Package
- `SetTestEnv(envVars map[string]string) func()` - Sets test environment variables with cleanup

### External Libraries
- `testify/assert` - Assertions for test validation
- `testify/require` - Required assertions that stop execution on failure
- `google/uuid` - UUID generation and validation
- `regexp` - Regular expression matching for token format validation

### Environment Variable Testing Pattern
```go
func TestSomeFunction(t *testing.T) {
    envVars := map[string]string{
        "SMTP_HOST": "test.smtp.com",
        "SMTP_PORT": "587",
    }
    cleanup := testutils.SetTestEnv(envVars)
    defer cleanup()

    // Test with environment variables set
    service := NewEmailService()
    // Assertions...
}
```

### Token Validation Pattern
```go
// Token format validation
parts := strings.Split(token, "-")
assert.GreaterOrEqual(t, len(parts), 7) // UUID(5) + hex(1) + timestamp(1)

// UUID validation
uuidPart := strings.Join(parts[:5], "-")
uuidPattern := `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`
matched, err := regexp.MatchString(uuidPattern, uuidPart)
```

## Example Test Execution

```bash
$ go test ./internal/utils -v
=== RUN   TestNewEmailService
=== RUN   TestNewEmailService/uses_default_values_when_no_environment_variables_are_set
=== RUN   TestNewEmailService/uses_environment_variables_when_set
=== RUN   TestNewEmailService/handles_partial_environment_variables
=== RUN   TestNewEmailService/handles_invalid_port_number_gracefully
--- PASS: TestNewEmailService (0.01s)
=== RUN   TestSendPasswordResetEmail
=== RUN   TestSendPasswordResetEmail/sends_password_reset_email_without_SMTP_credentials
=== RUN   TestSendPasswordResetEmail/constructs_correct_reset_URL_with_different_base_URL
--- PASS: TestSendPasswordResetEmail (0.02s)
=== RUN   TestGeneratePasswordResetToken
=== RUN   TestGeneratePasswordResetToken/generates_valid_token_format
=== RUN   TestGeneratePasswordResetToken/generates_unique_tokens
--- PASS: TestGeneratePasswordResetToken (0.15s)
PASS
ok      github.com/miftahulmahfuzh/lunch-delivery/internal/utils    0.180s
```

## Test Coverage

### Email Service (`email_test.go`)
- ✅ Email service initialization with environment variables
- ✅ SMTP configuration loading (5 environment variables)
- ✅ Password reset email sending functionality
- ✅ URL construction for reset links
- ✅ Email content formatting
- ✅ SMTP credential handling (with/without credentials)
- ✅ Error handling for SMTP connection failures
- ✅ Environment variable fallback logic
- ✅ Edge cases with special characters and long content
- ✅ Integration workflow testing
- ✅ Performance benchmarking

### Token Generation (`token_test.go`)
- ✅ Token format validation (UUID-hex-timestamp)
- ✅ Token uniqueness (1000+ token tests)
- ✅ Timestamp accuracy and consistency
- ✅ Security entropy validation
- ✅ Predictability resistance
- ✅ Character set validation (hex + hyphens only)
- ✅ Performance testing (generation speed)
- ✅ Edge cases (rapid generation, length consistency)
- ✅ Parallel generation benchmarking

## Maintenance Guidelines

### When to Update Tests
- When adding new environment variables for email configuration
- When modifying token generation algorithm
- When changing email template format
- When adding new utility functions
- When updating security requirements
- When modifying SMTP error handling

### Best Practices
- Always test both success and error scenarios
- Use environment variable mocking for configuration tests
- Test token uniqueness and security properties
- Include performance benchmarks for critical functions
- Test edge cases with special characters and empty values
- Verify proper cleanup of environment variables

### Common Issues
- **Environment Variables**: Always use `testutils.SetTestEnv()` with proper cleanup
- **SMTP Testing**: Test graceful fallback when no credentials are provided
- **Token Security**: Ensure sufficient entropy and uniqueness testing
- **Regular Expressions**: Use proper escaping for pattern matching
- **Time-Sensitive Tests**: Account for timestamp generation timing
- **Parallel Tests**: Ensure thread safety in token generation

## Integration Points

These tests support:
- Authentication service tests (through token generation)
- Email notification system tests
- Password reset functionality tests
- Environment configuration validation
- Security audit compliance
- Performance monitoring and optimization

## Security Considerations

### Token Generation Security
- **Entropy**: Tokens combine UUID + 16 random bytes + timestamp (64+ characters)
- **Unpredictability**: Uses crypto/rand for secure random number generation
- **Uniqueness**: UUID + random bytes ensure collision resistance
- **Temporal Component**: Timestamp allows for expiration validation

### Email Security
- **Credential Protection**: Graceful fallback when SMTP credentials unavailable
- **Content Validation**: Safe handling of user input in email content
- **URL Construction**: Proper token embedding in reset URLs
- **Logging**: Safe logging without exposing sensitive credentials