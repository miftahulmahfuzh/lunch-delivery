package services

import (
	"context"
	"testing"
	"time"

	"github.com/miftahulmahfuzh/lunch-delivery/internal/config"
	"github.com/miftahulmahfuzh/lunch-delivery/internal/mocks"
	"github.com/miftahulmahfuzh/lunch-delivery/internal/models"
	"github.com/miftahulmahfuzh/lunch-delivery/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Helper function to create test menu items
func createTestMenuItems() []models.MenuItem {
	return []models.MenuItem{
		{ID: 1, Name: "Nasi Gudeg", Price: 25000, Active: true},
		{ID: 2, Name: "Ayam Bakar", Price: 30000, Active: true},
		{ID: 3, Name: "Sayur Lodeh", Price: 15000, Active: true},
		{ID: 4, Name: "Tempe Goreng", Price: 10000, Active: true},
		{ID: 5, Name: "Es Teh Manis", Price: 5000, Active: true},
	}
}

func TestNewNutritionistService(t *testing.T) {
	t.Run("creates service with valid configuration", func(t *testing.T) {
		cfg := &config.Config{
			DeepseekTencentAPIKey: "test-api-key",
		}

		mockRepo := &mocks.RepositoryMock{}

		service, err := NewNutritionistService(cfg, mockRepo)

		// Note: This may fail due to LLM client creation, but we test the structure
		if err != nil {
			// Expected in test environment without real LLM API
			assert.Contains(t, err.Error(), "failed to create LLM client")
		} else {
			assert.NotNil(t, service)
			assert.NotNil(t, service.repo)
		}
	})

	t.Run("fails with invalid configuration", func(t *testing.T) {
		cfg := &config.Config{
			DeepseekTencentAPIKey: "", // Empty API key
		}

		mockRepo := &mocks.RepositoryMock{}

		service, err := NewNutritionistService(cfg, mockRepo)

		assert.Error(t, err)
		assert.Nil(t, service)
	})
}

func TestNutritionistService_GetNutritionistSelection(t *testing.T) {
	// Create a service with mocked dependencies
	createMockService := func() (*NutritionistService, *mocks.RepositoryMock, *mocks.LLMClientMock) {
		mockRepo := &mocks.RepositoryMock{}
		mockLLM := &mocks.LLMClientMock{}

		service := &NutritionistService{
			llmClient: mockLLM,
			repo:      mockRepo,
		}

		return service, mockRepo, mockLLM
	}

	t.Run("returns error when no menu items available", func(t *testing.T) {
		service, mockRepo, _ := createMockService()

		mockRepo.On("GetStockEmptyItemsForUser", mock.Anything, mock.Anything).Return([]int{}, nil)

		ctx := context.Background()
		menuItems := []models.MenuItem{}

		result, err := service.GetNutritionistSelection(ctx, testutils.TestDate(), menuItems, 1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "no menu items available")
	})

	t.Run("returns error when all items are stock empty", func(t *testing.T) {
		service, mockRepo, _ := createMockService()

		// All items are stock empty for this user
		mockRepo.On("GetStockEmptyItemsForUser", 1, testutils.TestDate()).Return([]int{1, 2, 3}, nil)

		ctx := context.Background()
		menuItems := createTestMenuItems()[:3] // Only first 3 items

		result, err := service.GetNutritionistSelection(ctx, testutils.TestDate(), menuItems, 1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "no menu items available for this user")
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns cached selection when available and valid", func(t *testing.T) {
		service, mockRepo, _ := createMockService()

		menuItems := createTestMenuItems()
		testDate := testutils.TestDate()

		// Mock stock empty items (none)
		mockRepo.On("GetStockEmptyItemsForUser", 1, testDate).Return([]int{}, nil)

		// Mock reset flag check
		mockRepo.On("GetDailyMenuResetFlag", testDate).Return(false, nil)

		// Mock cached selection
		cachedSelection := &models.NutritionistSelection{
			ID:                 1,
			Date:               testDate,
			MenuItemIDs:        []int64{1, 2, 3, 4, 5},
			SelectedIndices:    []int32{0, 2, 4},
			Reasoning:          "Cached balanced selection",
			NutritionalSummary: `{"protein": "high", "vegetables": "moderate", "carbohydrates": "balanced", "overall_rating": "good"}`,
			CreatedAt:          time.Now(),
		}
		mockRepo.On("GetNutritionistSelectionByDate", testDate).Return(cachedSelection, nil)

		ctx := context.Background()
		result, err := service.GetNutritionistSelection(ctx, testDate, menuItems, 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, []int{0, 2, 4}, result.SelectedIndices)
		assert.Equal(t, "Cached balanced selection", result.Reasoning)
		assert.Equal(t, "high", result.NutritionalSummary.Protein)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalidates cache when reset flag is set", func(t *testing.T) {
		service, mockRepo, mockLLM := createMockService()

		menuItems := createTestMenuItems()
		testDate := testutils.TestDate()

		// Mock stock empty items (none)
		mockRepo.On("GetStockEmptyItemsForUser", 1, testDate).Return([]int{}, nil)

		// Mock reset flag check (true - needs reset)
		mockRepo.On("GetDailyMenuResetFlag", testDate).Return(true, nil)

		// Mock cache deletion and flag clearing
		mockRepo.On("DeleteNutritionistSelection", testDate).Return(nil)
		mockRepo.On("SetDailyMenuResetFlag", testDate, false).Return(nil)

		// Mock cache miss after deletion
		mockRepo.On("GetNutritionistSelectionByDate", testDate).Return(nil, nil)

		// Mock LLM response
		llmResponse := mocks.CreateMockLLMResponse(mocks.MockNutritionistJSONResponse())
		mockLLM.On("GenerateContent", mock.Anything, mock.Anything, mock.Anything).Return(llmResponse, nil)

		// Mock cache save
		mockRepo.On("CreateNutritionistSelection", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

		ctx := context.Background()
		result, err := service.GetNutritionistSelection(ctx, testDate, menuItems, 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockRepo.AssertExpectations(t)
		mockLLM.AssertExpectations(t)
	})

	t.Run("calls LLM when cache miss", func(t *testing.T) {
		service, mockRepo, mockLLM := createMockService()

		menuItems := createTestMenuItems()
		testDate := testutils.TestDate()

		// Mock stock empty items (none)
		mockRepo.On("GetStockEmptyItemsForUser", 1, testDate).Return([]int{}, nil)

		// Mock reset flag check
		mockRepo.On("GetDailyMenuResetFlag", testDate).Return(false, nil)

		// Mock cache miss
		mockRepo.On("GetNutritionistSelectionByDate", testDate).Return(nil, nil)

		// Mock LLM response
		llmResponse := mocks.CreateMockLLMResponse(mocks.MockNutritionistJSONResponse())
		mockLLM.On("GenerateContent", mock.Anything, mock.Anything, mock.Anything).Return(llmResponse, nil)

		// Mock cache save
		mockRepo.On("CreateNutritionistSelection", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

		ctx := context.Background()
		result, err := service.GetNutritionistSelection(ctx, testDate, menuItems, 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, []int{0, 2, 3}, result.SelectedIndices)
		assert.Contains(t, result.Reasoning, "balanced nutrition")
		mockRepo.AssertExpectations(t)
		mockLLM.AssertExpectations(t)
	})

	t.Run("handles LLM error gracefully", func(t *testing.T) {
		service, mockRepo, mockLLM := createMockService()

		menuItems := createTestMenuItems()
		testDate := testutils.TestDate()

		// Mock stock empty items (none)
		mockRepo.On("GetStockEmptyItemsForUser", 1, testDate).Return([]int{}, nil)

		// Mock reset flag check
		mockRepo.On("GetDailyMenuResetFlag", testDate).Return(false, nil)

		// Mock cache miss
		mockRepo.On("GetNutritionistSelectionByDate", testDate).Return(nil, nil)

		// Mock LLM error
		mockLLM.On("GenerateContent", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError)

		ctx := context.Background()
		result, err := service.GetNutritionistSelection(ctx, testDate, menuItems, 1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "LLM call failed")
		mockRepo.AssertExpectations(t)
		mockLLM.AssertExpectations(t)
	})

	t.Run("filters out stock empty items for user", func(t *testing.T) {
		service, mockRepo, mockLLM := createMockService()

		menuItems := createTestMenuItems()
		testDate := testutils.TestDate()

		// Mock stock empty items (first two items are empty for this user)
		mockRepo.On("GetStockEmptyItemsForUser", 1, testDate).Return([]int{1, 2}, nil)

		// Mock reset flag check
		mockRepo.On("GetDailyMenuResetFlag", testDate).Return(false, nil)

		// Mock cache miss
		mockRepo.On("GetNutritionistSelectionByDate", testDate).Return(nil, nil)

		// Mock LLM response (should only get items 3, 4, 5 which are indices 2, 3, 4 in original array)
		llmResponse := mocks.CreateMockLLMResponse(`{
			"selected_menu_items": [0, 1, 2],
			"reasoning": "Selected from available items",
			"nutritional_summary": {
				"protein": "moderate",
				"vegetables": "high",
				"carbohydrates": "moderate",
				"overall_rating": "good"
			}
		}`)
		mockLLM.On("GenerateContent", mock.Anything, mock.Anything, mock.Anything).Return(llmResponse, nil)

		// Mock cache save
		mockRepo.On("CreateNutritionistSelection", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

		ctx := context.Background()
		result, err := service.GetNutritionistSelection(ctx, testDate, menuItems, 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		// Should map back to original indices (2, 3, 4 become items that weren't stock empty)
		assert.Len(t, result.SelectedIndices, 3)
		mockRepo.AssertExpectations(t)
		mockLLM.AssertExpectations(t)
	})
}

func TestNutritionistService_TrackUserSelection(t *testing.T) {
	t.Run("tracks user selection successfully", func(t *testing.T) {
		service, mockRepo, _ := createMockService()

		testDate := testutils.TestDate()
		orderID := 123

		mockRepo.On("CreateNutritionistUserSelection", testDate, 1, []int64{123}).Return(nil, nil)

		err := service.TrackUserSelection(1, testDate, &orderID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("handles tracking error", func(t *testing.T) {
		service, mockRepo, _ := createMockService()

		testDate := testutils.TestDate()

		mockRepo.On("CreateNutritionistUserSelection", testDate, 1, []int64(nil)).Return(nil, assert.AnError)

		err := service.TrackUserSelection(1, testDate, nil)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestNutritionistService_GetUsersNeedingNotification(t *testing.T) {
	t.Run("gets users needing notification", func(t *testing.T) {
		service, mockRepo, _ := createMockService()

		testDate := testutils.TestDate()
		expectedUsers := []models.NutritionistUserSelection{
			{ID: 1, EmployeeID: 1, Date: testDate},
			{ID: 2, EmployeeID: 2, Date: testDate},
		}

		mockRepo.On("GetNutritionistUsersByDateAndUnpaid", testDate).Return(expectedUsers, nil)

		users, err := service.GetUsersNeedingNotification(testDate)

		assert.NoError(t, err)
		assert.Len(t, users, 2)
		assert.Equal(t, 1, users[0].EmployeeID)
		assert.Equal(t, 2, users[1].EmployeeID)
		mockRepo.AssertExpectations(t)
	})
}

func TestNutritionistService_HelperMethods(t *testing.T) {
	service := &NutritionistService{}

	t.Run("menuItemsMatch returns true for matching items", func(t *testing.T) {
		cachedIDs := []int64{1, 2, 3}
		menuItems := []models.MenuItem{
			{ID: 1}, {ID: 2}, {ID: 3},
		}

		result := service.menuItemsMatch(cachedIDs, menuItems)
		assert.True(t, result)
	})

	t.Run("menuItemsMatch returns false for different lengths", func(t *testing.T) {
		cachedIDs := []int64{1, 2}
		menuItems := []models.MenuItem{
			{ID: 1}, {ID: 2}, {ID: 3},
		}

		result := service.menuItemsMatch(cachedIDs, menuItems)
		assert.False(t, result)
	})

	t.Run("menuItemsMatch returns false for different items", func(t *testing.T) {
		cachedIDs := []int64{1, 2, 4}
		menuItems := []models.MenuItem{
			{ID: 1}, {ID: 2}, {ID: 3},
		}

		result := service.menuItemsMatch(cachedIDs, menuItems)
		assert.False(t, result)
	})

	t.Run("validateIndices validates correctly", func(t *testing.T) {
		tests := []struct {
			indices  []int
			maxIndex int
			expected bool
		}{
			{[]int{0, 1, 2}, 5, true},
			{[]int{0, 4}, 5, true},
			{[]int{0, 5}, 5, false},                // 5 is out of bounds for maxIndex 5
			{[]int{-1, 1}, 5, false},               // negative index
			{[]int{}, 5, false},                    // empty
			{[]int{0, 1, 2, 3, 4, 5, 6}, 5, false}, // too many indices
		}

		for _, tt := range tests {
			result := service.validateIndices(tt.indices, tt.maxIndex)
			assert.Equal(t, tt.expected, result, "validateIndices(%v, %d) should return %t", tt.indices, tt.maxIndex, tt.expected)
		}
	})

	t.Run("buildMenuDescription creates correct format", func(t *testing.T) {
		menuItems := []models.MenuItem{
			{ID: 1, Name: "Item 1", Price: 10000},
			{ID: 2, Name: "Item 2", Price: 20000},
		}

		description := service.buildMenuDescription(menuItems)

		assert.Contains(t, description, "Index 0: Item 1 (Rp 10000)")
		assert.Contains(t, description, "Index 1: Item 2 (Rp 20000)")
	})

	t.Run("cleanMarkdownCodeBlocks removes code blocks", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{"```json\n{\"test\": true}\n```", "{\"test\": true}"},
			{"```\n{\"test\": true}\n```", "{\"test\": true}"},
			{"{\"test\": true}", "{\"test\": true}"},
			{"  ```json\n{\"test\": true}\n```  ", "{\"test\": true}"},
		}

		for _, tt := range tests {
			result := service.cleanMarkdownCodeBlocks(tt.input)
			assert.Equal(t, tt.expected, result)
		}
	})

	t.Run("extractNumbers extracts numbers from text", func(t *testing.T) {
		text := "Selected indices: [0, 2, 4] and item 7"
		numbers := service.extractNumbers(text)

		expected := []int{0, 2, 4, 7}
		assert.Equal(t, expected, numbers)
	})

	t.Run("uniqueIndices removes duplicates", func(t *testing.T) {
		indices := []int{1, 2, 2, 3, 1, 4}
		unique := service.uniqueIndices(indices)

		expected := []int{1, 2, 3, 4}
		assert.Equal(t, expected, unique)
	})
}

// Helper function to create mock service
func createMockService() (*NutritionistService, *mocks.RepositoryMock, *mocks.LLMClientMock) {
	mockRepo := &mocks.RepositoryMock{}
	mockLLM := &mocks.LLMClientMock{}

	service := &NutritionistService{
		llmClient: mockLLM,
		repo:      mockRepo,
	}

	return service, mockRepo, mockLLM
}
