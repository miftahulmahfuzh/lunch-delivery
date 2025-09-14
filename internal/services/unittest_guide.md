# Services Package Unit Tests Guide

## Overview

This guide explains the unit tests for the `internal/services` package, which contains business logic services, particularly the AI-powered nutritionist service that provides meal recommendations.

## Test Structure

### Files Tested
- `nutritionist.go` - AI nutritionist service with complex meal selection logic

### Test Files
- `nutritionist_test.go` - Comprehensive tests for nutritionist service functionality

## Test Categories

### 1. Service Creation Tests (`TestNewNutritionistService`)

**Purpose**: Test the service constructor and initialization.

**Scenarios Tested**:
- Creates service with valid LLM configuration
- Fails with invalid/missing API configuration
- Proper dependency injection (repository, LLM client)

**Key Components**:
- LLM client initialization
- Configuration validation
- Repository dependency setup

### 2. Core Selection Logic Tests (`TestNutritionistService_GetNutritionistSelection`)

**Purpose**: Test the main business logic for AI-powered meal selection.

#### A. Input Validation
- **Empty Menu Items**: Returns error when no menu items available
- **All Stock Empty**: Returns error when all items are out of stock for user
- **User-Specific Stock Filtering**: Properly filters out stock empty items per user

#### B. Caching Logic
- **Cache Hit**: Returns cached selection when available and valid
- **Cache Miss**: Calls LLM when no cached selection exists
- **Cache Invalidation**: Invalidates cache when reset flag is set
- **Cache Validation**: Validates cached selection against current menu items

#### C. Stock Management Integration
- **User-Specific Stock**: Filters menu items based on user's stock empty items
- **Stock Empty Handling**: Gracefully handles items that become unavailable
- **Dynamic Filtering**: Adjusts available items per user context

#### D. LLM Integration
- **Successful LLM Call**: Processes valid LLM responses
- **LLM Error Handling**: Handles LLM service failures gracefully
- **Response Parsing**: Correctly parses JSON responses from LLM
- **Fallback Logic**: Provides fallback when LLM parsing fails

#### E. Menu Reset Logic
- **Reset Flag Detection**: Detects when admin has reset the menu
- **Cache Invalidation**: Clears cache when reset flag is detected
- **Flag Clearing**: Resets the flag after processing

### 3. User Tracking Tests (`TestNutritionistService_TrackUserSelection`)

**Purpose**: Test user selection tracking functionality.

**Scenarios Tested**:
- Successfully tracks user selections
- Handles tracking errors gracefully
- Links selections to specific orders

### 4. Notification Support Tests (`TestNutritionistService_GetUsersNeedingNotification`)

**Purpose**: Test notification system integration.

**Scenarios Tested**:
- Retrieves users who need notifications after menu changes
- Handles empty result sets
- Filters by unpaid orders correctly

### 5. Helper Method Tests (`TestNutritionistService_HelperMethods`)

**Purpose**: Test utility methods and internal logic.

#### A. Menu Validation (`menuItemsMatch`)
- **Exact Match**: Returns true when cached items match current menu
- **Different Lengths**: Returns false for different item counts
- **Different Items**: Returns false when item IDs don't match

#### B. Index Validation (`validateIndices`)
- **Valid Ranges**: Accepts indices within bounds
- **Out of Bounds**: Rejects indices outside valid range
- **Negative Indices**: Rejects negative index values
- **Empty Arrays**: Handles empty index arrays
- **Too Many Items**: Limits maximum selection count

#### C. Menu Description (`buildMenuDescription`)
- **Format Consistency**: Creates properly formatted menu descriptions
- **Index Mapping**: Maps items to correct indices
- **Price Formatting**: Includes proper price formatting

#### D. Text Processing (`cleanMarkdownCodeBlocks`)
- **JSON Code Blocks**: Removes ```json``` markdown formatting
- **Generic Code Blocks**: Removes ``` markdown formatting
- **Clean Text**: Passes through already clean text
- **Whitespace Handling**: Properly trims whitespace

#### E. Number Extraction (`extractNumbers`)
- **Array Format**: Extracts numbers from array-like text
- **Mixed Text**: Finds numbers in mixed content
- **Multiple Numbers**: Handles multiple number extraction

#### F. Deduplication (`uniqueIndices`)
- **Duplicate Removal**: Removes duplicate indices
- **Order Preservation**: Maintains original order of first occurrence

## Mock Architecture

### Repository Mock (`mocks.RepositoryMock`)
Mocks all database operations:
- Stock empty item retrieval
- Menu reset flag management
- Cached selection storage/retrieval
- User selection tracking

### LLM Client Mock (`mocks.LLMClientMock`)
Mocks AI service interactions:
- Content generation requests
- Response formatting
- Error simulation

### Mock Response Helpers
- `MockNutritionistJSONResponse()` - Valid AI response
- `MockInvalidJSONResponse()` - Malformed response
- `MockPartialJSONResponse()` - Incomplete response

## Running the Tests

### Run All Service Tests
```bash
go test ./internal/services
```

### Run Specific Test Categories
```bash
# Service creation tests
go test ./internal/services -run TestNewNutritionistService

# Core selection logic
go test ./internal/services -run TestNutritionistService_GetNutritionistSelection

# Helper methods
go test ./internal/services -run TestNutritionistService_HelperMethods
```

### Run Tests with Verbose Output
```bash
go test ./internal/services -v
```

### Run Tests with Coverage
```bash
go test ./internal/services -cover
```

## Test Patterns

### 1. Mock Service Creation
```go
func createMockService() (*NutritionistService, *mocks.RepositoryMock, *mocks.LLMClientMock) {
    mockRepo := &mocks.RepositoryMock{}
    mockLLM := &mocks.LLMClientMock{}

    service := &NutritionistService{
        llmClient: mockLLM,
        repo:      mockRepo,
    }

    return service, mockRepo, mockLLM
}
```

### 2. Repository Mock Setup
```go
// Mock stock empty items
mockRepo.On("GetStockEmptyItemsForUser", 1, testDate).Return([]int{}, nil)

// Mock reset flag check
mockRepo.On("GetDailyMenuResetFlag", testDate).Return(false, nil)

// Mock cached selection
mockRepo.On("GetNutritionistSelectionByDate", testDate).Return(cachedSelection, nil)
```

### 3. LLM Mock Setup
```go
// Mock successful LLM response
llmResponse := mocks.CreateMockLLMResponse(mocks.MockNutritionistJSONResponse())
mockLLM.On("GenerateContent", mock.Anything, mock.Anything, mock.Anything).Return(llmResponse, nil)
```

## Test Data Management

### Menu Items
Uses `testutils.MockMenuItems()` for consistent test data:
- Nasi Gudeg (25000)
- Ayam Bakar (30000)
- Sayur Lodeh (15000)
- Tempe Goreng (10000)
- Es Teh Manis (5000)

### Test Dates
Uses `testutils.TestDate()` for consistent date handling.

### Mock Responses
Predefined JSON responses for various test scenarios.

## Example Test Execution

```bash
$ go test ./internal/services -v
=== RUN   TestNewNutritionistService
=== RUN   TestNewNutritionistService/creates_service_with_valid_configuration
=== RUN   TestNewNutritionistService/fails_with_invalid_configuration
--- PASS: TestNewNutritionistService (0.01s)
=== RUN   TestNutritionistService_GetNutritionistSelection
=== RUN   TestNutritionistService_GetNutritionistSelection/returns_error_when_no_menu_items_available
=== RUN   TestNutritionistService_GetNutritionistSelection/returns_cached_selection_when_available_and_valid
--- PASS: TestNutritionistService_GetNutritionistSelection (0.05s)
PASS
ok      github.com/miftahulmahfuzh/lunch-delivery/internal/services    0.078s
```

## Complex Test Scenarios

### 1. Cache Invalidation Flow
```go
t.Run("invalidates cache when reset flag is set", func(t *testing.T) {
    // Setup: reset flag is true
    mockRepo.On("GetDailyMenuResetFlag", testDate).Return(true, nil)

    // Expect: cache deletion and flag clearing
    mockRepo.On("DeleteNutritionistSelection", testDate).Return(nil)
    mockRepo.On("SetDailyMenuResetFlag", testDate, false).Return(nil)

    // Then: new LLM call and cache update
    // ... additional expectations
})
```

### 2. User-Specific Stock Filtering
```go
t.Run("filters out stock empty items for user", func(t *testing.T) {
    // Setup: first two menu items are stock empty for user
    mockRepo.On("GetStockEmptyItemsForUser", 1, testDate).Return([]int{1, 2}, nil)

    // Expect: LLM called with only available items (3, 4, 5)
    // Result: indices mapped back to original menu
})
```

### 3. LLM Response Processing
```go
t.Run("calls LLM when cache miss", func(t *testing.T) {
    // Setup: no cached selection
    mockRepo.On("GetNutritionistSelectionByDate", testDate).Return(nil, nil)

    // Expect: LLM call with proper menu description
    mockLLM.On("GenerateContent", mock.Anything, mock.Anything, mock.Anything).Return(validResponse, nil)

    // Expect: cache save with processed results
    mockRepo.On("CreateNutritionistSelection", ...).Return(nil, nil)
})
```

## Test Coverage

### Core Functionality
- ✅ Service initialization and configuration
- ✅ AI-powered meal selection logic
- ✅ Cache management and invalidation
- ✅ User-specific stock filtering
- ✅ LLM integration and response processing
- ✅ Error handling and fallback logic

### Helper Methods
- ✅ Menu validation and comparison
- ✅ Index validation and bounds checking
- ✅ Text processing and cleanup
- ✅ Number extraction and parsing
- ✅ Deduplication algorithms

### Integration Points
- ✅ Repository integration
- ✅ LLM client integration
- ✅ Configuration management
- ✅ User notification support

## Maintenance Guidelines

### When to Update Tests
- When adding new meal selection algorithms
- When modifying cache invalidation logic
- When changing LLM integration patterns
- When adding new user preference features
- When updating response parsing logic

### Best Practices
- Always mock external dependencies (LLM, database)
- Test both success and error scenarios
- Verify all mock expectations are met
- Use consistent test data from `testutils`
- Test complex business logic thoroughly
- Include edge cases and boundary conditions

### Common Issues
- **Mock Expectations**: Ensure all repository and LLM calls are mocked
- **JSON Parsing**: Test various LLM response formats
- **Context Handling**: Properly handle context cancellation
- **Memory Management**: Clean up mock objects after tests
- **Time Sensitivity**: Use consistent test dates

## Integration Points

These service tests support:
- Handler integration testing
- End-to-end API testing
- Performance testing of AI recommendations
- User experience validation
- Business logic verification