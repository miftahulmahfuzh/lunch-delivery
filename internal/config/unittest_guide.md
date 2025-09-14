# Config Package Unit Tests Guide

## Overview

This guide explains the unit tests for the `internal/config` package, which handles application configuration loading from environment variables and `.env` files.

## Test Structure

### Files Tested
- `config.go` - Configuration loading and environment variable handling

### Test File
- `config_test.go` - Comprehensive tests for configuration functionality

## Test Categories

### 1. Configuration Loading Tests (`TestLoad`)

**Purpose**: Test the `Load()` function that creates configuration from environment variables.

**Scenarios Tested**:
- Default values when no environment variables are set
- Using environment variables when they are set
- Partial environment variable configuration
- Environment variable precedence over defaults

**Key Test Cases**:
- `default values when no environment variables are set`
- `uses environment variables when set`
- `partial environment variables with defaults`

### 2. Environment Variable Helper Tests (`TestGetEnv`)

**Purpose**: Test the `getEnv()` helper function for environment variable retrieval.

**Scenarios Tested**:
- Returns default value when environment variable is not set
- Returns environment value when set
- Handles empty environment variables
- Handles empty default values

### 3. .env File Loading Tests (`TestLoadWithDotEnvFile`)

**Purpose**: Test configuration loading with `.env` file support.

**Scenarios Tested**:
- Loads configuration from `.env` file when it exists
- Handles missing `.env` file gracefully
- Environment variables override `.env` file values

### 4. Configuration Validation Tests (`TestConfig_ValidationScenarios`)

**Purpose**: Test various configuration scenarios for completeness.

**Scenarios Tested**:
- Configuration with all fields populated
- Configuration with minimal required fields
- Validates field types and constraints

### 5. Performance Tests

**Benchmark Tests**:
- `BenchmarkLoad` - Measures configuration loading performance
- `BenchmarkGetEnv` - Measures environment variable retrieval performance

## Running the Tests

### Run All Config Tests
```bash
go test ./internal/config
```

### Run Specific Test
```bash
go test ./internal/config -run TestLoad
```

### Run Tests with Verbose Output
```bash
go test ./internal/config -v
```

### Run Benchmark Tests
```bash
go test ./internal/config -bench=.
```

### Run Tests with Coverage
```bash
go test ./internal/config -cover
```

## Test Utilities Used

- **Environment Management**: Uses `testutils.SetTestEnv()` to safely set and restore environment variables
- **Temporary Directories**: Uses `testutils.CreateTempDir()` for `.env` file testing
- **Mock Data**: Uses `testutils.MockEnvironment()` and `testutils.MockConfig()` for consistent test data

## Example Test Execution

```bash
# Example output when running tests
$ go test ./internal/config -v
=== RUN   TestLoad
=== RUN   TestLoad/default_values_when_no_environment_variables_are_set
=== RUN   TestLoad/uses_environment_variables_when_set
=== RUN   TestLoad/partial_environment_variables_with_defaults
--- PASS: TestLoad (0.01s)
    --- PASS: TestLoad/default_values_when_no_environment_variables_are_set (0.00s)
    --- PASS: TestLoad/uses_environment_variables_when_set (0.00s)
    --- PASS: TestLoad/partial_environment_variables_with_defaults (0.00s)
=== RUN   TestGetEnv
--- PASS: TestGetEnv (0.00s)
PASS
ok      github.com/miftahulmahfuzh/lunch-delivery/internal/config    0.012s
```

## Common Test Patterns

### 1. Environment Variable Testing
```go
// Setup environment
cleanup := testutils.SetTestEnv(tt.envVars)
defer cleanup()

// Test configuration loading
cfg, err := Load()

// Verify results
require.NoError(t, err)
assert.Equal(t, expected, cfg.DBHost)
```

### 2. File-based Configuration Testing
```go
// Create temporary directory and .env file
tempDir, cleanup := testutils.CreateTempDir(t)
defer cleanup()

// Change to temp directory
err = os.Chdir(tempDir)
require.NoError(t, err)

// Test with .env file
cfg, err := Load()
```

## Test Coverage

The tests cover:
- ✅ All public functions (`Load`, `getEnv`)
- ✅ Environment variable handling
- ✅ Default value logic
- ✅ .env file integration
- ✅ Error handling
- ✅ Edge cases (empty values, missing files)
- ✅ Performance benchmarks

## Maintenance Guidelines

### When to Update Tests
- When adding new configuration fields
- When changing default values
- When modifying environment variable names
- When adding new configuration sources

### Best Practices
- Always use `testutils.SetTestEnv()` for environment variable manipulation
- Clean up temporary files and directories
- Test both success and failure scenarios
- Verify all configuration fields are properly set
- Include performance benchmarks for configuration loading

### Common Issues
- **Env Variable Pollution**: Always clean up environment variables after tests
- **File System State**: Ensure temporary files are cleaned up
- **Default Values**: Keep test default values in sync with actual defaults
- **Time Zones**: Be careful with time-based configuration fields

## Integration with Other Tests

The config tests provide foundation for:
- Database connection tests (using DB configuration)
- LLM client tests (using API configuration)
- Service integration tests (using timeout configuration)
- Email service tests (using SMTP configuration)