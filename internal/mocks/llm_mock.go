package mocks

import (
	"github.com/miftahulmahfuzh/lunch-delivery/internal/interfaces"
	"github.com/stretchr/testify/mock"
)

// LLMClientMock is a mock implementation of the LLMClientInterface
type LLMClientMock struct {
	mock.Mock
}

// Compile-time check to ensure LLMClientMock implements LLMClientInterface
var _ interfaces.LLMClientInterface = (*LLMClientMock)(nil)

func (m *LLMClientMock) GenerateContent(systemPrompt, userPrompt, temperature string) (string, error) {
	args := m.Called(systemPrompt, userPrompt, temperature)
	return args.String(0), args.Error(1)
}

// Helper method to create mock responses
func CreateMockLLMResponse(content string) string {
	return content
}

// Common mock responses for testing
func MockNutritionistJSONResponse() string {
	return `{
		"selected_menu_items": [0, 2, 3],
		"reasoning": "Selected items provide balanced nutrition with protein from chicken, vegetables from salad, and carbohydrates from rice.",
		"nutritional_summary": {
			"protein": "high",
			"vegetables": "high",
			"carbohydrates": "moderate",
			"overall_rating": "excellent"
		}
	}`
}

func MockInvalidJSONResponse() string {
	return `{invalid json response`
}

func MockPartialJSONResponse() string {
	return `{
		"selected_menu_items": [0, 1],
		"reasoning": "Good combination"
	}`
}