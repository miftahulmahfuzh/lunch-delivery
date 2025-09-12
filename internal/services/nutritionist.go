// internal/services/nutritionist.go
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/miftahulmahfuzh/lunch-delivery/internal/config"
	"github.com/miftahulmahfuzh/lunch-delivery/internal/llm"
	"github.com/miftahulmahfuzh/lunch-delivery/internal/models"
	"github.com/rs/zerolog/log"
	"github.com/tmc/langchaingo/llms"
)

type NutritionistService struct {
	llmClient *llm.Client
	repo      *models.Repository
}

type NutritionistResponse struct {
	SelectedIndices     []int               `json:"selected_menu_items"`
	Reasoning          string              `json:"reasoning"`
	NutritionalSummary NutritionalSummary  `json:"nutritional_summary"`
}

type NutritionalSummary struct {
	Protein       string `json:"protein"`
	Vegetables    string `json:"vegetables"`
	Carbohydrates string `json:"carbohydrates"`
	OverallRating string `json:"overall_rating"`
}

func NewNutritionistService(cfg *config.Config, repo *models.Repository) (*NutritionistService, error) {
	llmClient, err := llm.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create LLM client: %w", err)
	}

	return &NutritionistService{
		llmClient: llmClient,
		repo:      repo,
	}, nil
}

func (s *NutritionistService) GetNutritionistSelection(ctx context.Context, date time.Time, menuItems []models.MenuItem) (*NutritionistResponse, error) {
	// Use all menu items since stock management is now user-specific
	// Each user will see their own stock constraints when placing orders
	availableMenuItems := menuItems
	if len(availableMenuItems) == 0 {
		return nil, fmt.Errorf("no menu items available")
	}

	// 3. Check if admin has set reset flag
	resetFlag, err := s.repo.GetDailyMenuResetFlag(date)
	if err == nil && resetFlag {
		log.Info().Msg("Reset flag detected - invalidating cache and clearing flag")
		s.repo.DeleteNutritionistSelection(date)
		s.repo.SetDailyMenuResetFlag(date, false) // Clear the flag
	}

	// 4. Check cache first by date
	cached, err := s.repo.GetNutritionistSelectionByDate(date)
	if err == nil && cached != nil {
		// Cache hit - validate menu items match (no need to check stock conflicts with user-specific model)
		if s.menuItemsMatch(cached.MenuItemIDs, availableMenuItems) {
			log.Info().Msg("Cache hit - returning cached nutritionist selection")
			return s.convertCachedToResponse(cached), nil
		}
		// Menu changed, invalidate cache
		log.Info().Msg("Menu items changed - invalidating cache")
		s.repo.DeleteNutritionistSelection(date)
	}

	log.Info().Msg("Cache miss or menu changed - calling LLM for nutritionist selection")
	
	// 5. Cache miss or menu changed - call LLM with available items
	response, err := s.callLLMForSelection(ctx, availableMenuItems)
	if err != nil {
		return nil, fmt.Errorf("LLM call failed: %w", err)
	}

	// 6. Map the response indices back to original menu items and save to cache
	mappedResponse := s.mapIndicesToOriginalMenu(response, availableMenuItems, menuItems)
	if err := s.saveToCacheIfValid(date, menuItems, mappedResponse); err != nil {
		log.Error().Err(err).Msg("Failed to save to cache")
		// Don't fail the request if caching fails
	}

	return response, nil
}

// Method to track that user used nutritionist selection
func (s *NutritionistService) TrackUserSelection(employeeID int, date time.Time, orderID *int) error {
	return s.repo.CreateNutritionistUserSelection(employeeID, date, orderID)
}

// Method to get users who need notification after menu reset
func (s *NutritionistService) GetUsersNeedingNotification(date time.Time) ([]models.NutritionistUserSelection, error) {
	return s.repo.GetNutritionistUsersByDateAndUnpaid(date)
}


// Map indices from available items back to original menu items
func (s *NutritionistService) mapIndicesToOriginalMenu(response *NutritionistResponse, availableItems []models.MenuItem, originalItems []models.MenuItem) *NutritionistResponse {
	// Create a map from available item ID to original menu index
	availableIDToOriginalIndex := make(map[int]int)
	for originalIndex, originalItem := range originalItems {
		availableIDToOriginalIndex[originalItem.ID] = originalIndex
	}

	// Map the selected indices
	var mappedIndices []int
	for _, availableIndex := range response.SelectedIndices {
		if availableIndex < len(availableItems) {
			availableItem := availableItems[availableIndex]
			if originalIndex, exists := availableIDToOriginalIndex[availableItem.ID]; exists {
				mappedIndices = append(mappedIndices, originalIndex)
			}
		}
	}

	return &NutritionistResponse{
		SelectedIndices:    mappedIndices,
		Reasoning:          response.Reasoning,
		NutritionalSummary: response.NutritionalSummary,
	}
}

func (s *NutritionistService) menuItemsMatch(cachedIDs []int64, menuItems []models.MenuItem) bool {
	if len(cachedIDs) != len(menuItems) {
		return false
	}

	// Convert menu items to ID set for comparison
	menuIDs := make(map[int64]bool)
	for _, item := range menuItems {
		menuIDs[int64(item.ID)] = true
	}

	// Check if all cached IDs exist in current menu
	for _, cachedID := range cachedIDs {
		if !menuIDs[cachedID] {
			return false
		}
	}

	return true
}

func (s *NutritionistService) convertCachedToResponse(cached *models.NutritionistSelection) *NutritionistResponse {
	// Parse the nutritional summary JSON
	var summary NutritionalSummary
	if err := json.Unmarshal([]byte(cached.NutritionalSummary), &summary); err != nil {
		log.Warn().Err(err).Msg("Failed to parse cached nutritional summary")
		summary = NutritionalSummary{
			Protein:       "unknown",
			Vegetables:    "unknown", 
			Carbohydrates: "unknown",
			OverallRating: "balanced",
		}
	}

	// Convert pq.Int32Array to []int
	selectedIndices := make([]int, len(cached.SelectedIndices))
	for i, idx := range cached.SelectedIndices {
		selectedIndices[i] = int(idx)
	}

	return &NutritionistResponse{
		SelectedIndices:    selectedIndices,
		Reasoning:          cached.Reasoning,
		NutritionalSummary: summary,
	}
}

func (s *NutritionistService) callLLMForSelection(ctx context.Context, menuItems []models.MenuItem) (*NutritionistResponse, error) {
	// Build menu items description for LLM
	menuDescription := s.buildMenuDescription(menuItems)
	
	systemPrompt := `You are a highly experienced nutritionist. Your task is to select the most healthy and balanced meal combination from the available menu items.

CRITICAL REQUIREMENTS:
1. You MUST respond with ONLY a valid JSON object in this exact format:
{
  "selected_menu_items": [0, 2, 4],
  "reasoning": "Brief explanation of why these items provide balanced nutrition",
  "nutritional_summary": {
    "protein": "high|moderate|low",
    "vegetables": "high|moderate|low|none", 
    "carbohydrates": "high|moderate|low",
    "overall_rating": "excellent|good|balanced|adequate"
  }
}

2. The "selected_menu_items" array MUST contain INDICES (0-based) of menu items, not IDs
3. Select 2-4 items that provide the most balanced nutrition
4. Prioritize: protein sources, vegetables, whole grains, balanced portions
5. Avoid: excessive fried foods, too much sugar, unbalanced combinations

Available menu items (with their indices):`

	userPrompt := menuDescription

	messages := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemPrompt),
		llms.TextParts(llms.ChatMessageTypeHuman, userPrompt),
	}

	response, err := s.llmClient.GenerateContent(ctx, messages)
	if err != nil {
		return nil, err
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("LLM returned no response")
	}

	content := response.Choices[0].Content
	log.Info().Str("llm_response", content).Msg("LLM raw response")

	// Parse the JSON response
	return s.parseStructuredResponse(content, len(menuItems))
}

func (s *NutritionistService) buildMenuDescription(menuItems []models.MenuItem) string {
	var builder strings.Builder
	for i, item := range menuItems {
		builder.WriteString(fmt.Sprintf("\nIndex %d: %s (Rp %d)", i, item.Name, item.Price))
	}
	return builder.String()
}

func (s *NutritionistService) parseStructuredResponse(content string, maxIndex int) (*NutritionistResponse, error) {
	// Try to parse as JSON first
	var response NutritionistResponse
	if err := json.Unmarshal([]byte(content), &response); err == nil {
		// Validate indices
		if s.validateIndices(response.SelectedIndices, maxIndex) {
			return &response, nil
		}
		log.Warn().Interface("indices", response.SelectedIndices).Msg("Invalid indices in JSON response")
	}

	// Fallback: try to extract indices using regex/parsing
	log.Warn().Msg("JSON parsing failed, attempting fallback parsing")
	return s.fallbackParseResponse(content, maxIndex)
}

func (s *NutritionistService) fallbackParseResponse(content string, maxIndex int) (*NutritionistResponse, error) {
	// Look for array-like patterns in the content
	// This is a simple fallback - you can make it more sophisticated
	var indices []int
	
	// Try to find numbers in brackets or array-like format
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "selected") || strings.Contains(line, "indices") || strings.Contains(line, "[") {
			// Extract numbers from this line
			nums := s.extractNumbers(line)
			for _, num := range nums {
				if num >= 0 && num < maxIndex {
					indices = append(indices, num)
				}
			}
		}
	}

	if len(indices) == 0 {
		return nil, fmt.Errorf("could not extract valid indices from LLM response")
	}

	// Remove duplicates and limit to reasonable count
	indices = s.uniqueIndices(indices)
	if len(indices) > 4 {
		indices = indices[:4] // Limit to max 4 items
	}

	return &NutritionistResponse{
		SelectedIndices: indices,
		Reasoning:       "AI-selected balanced combination",
		NutritionalSummary: NutritionalSummary{
			Protein:       "balanced",
			Vegetables:    "adequate",
			Carbohydrates: "balanced",
			OverallRating: "good",
		},
	}, nil
}

func (s *NutritionistService) extractNumbers(text string) []int {
	var numbers []int
	words := strings.Fields(text)
	for _, word := range words {
		// Clean the word of common punctuation
		word = strings.Trim(word, "[](),")
		if num, err := strconv.Atoi(word); err == nil {
			numbers = append(numbers, num)
		}
	}
	return numbers
}

func (s *NutritionistService) uniqueIndices(indices []int) []int {
	seen := make(map[int]bool)
	var unique []int
	for _, idx := range indices {
		if !seen[idx] {
			seen[idx] = true
			unique = append(unique, idx)
		}
	}
	return unique
}

func (s *NutritionistService) validateIndices(indices []int, maxIndex int) bool {
	if len(indices) == 0 || len(indices) > 6 {
		return false
	}
	for _, idx := range indices {
		if idx < 0 || idx >= maxIndex {
			return false
		}
	}
	return true
}

func (s *NutritionistService) saveToCacheIfValid(date time.Time, menuItems []models.MenuItem, response *NutritionistResponse) error {
	// Convert menu items to IDs
	var menuItemIDs []int64
	for _, item := range menuItems {
		menuItemIDs = append(menuItemIDs, int64(item.ID))
	}

	// Convert indices to int32 for database
	var selectedIndices []int32
	for _, idx := range response.SelectedIndices {
		selectedIndices = append(selectedIndices, int32(idx))
	}

	// Convert nutritional summary to JSON
	summaryJSON, err := json.Marshal(response.NutritionalSummary)
	if err != nil {
		return fmt.Errorf("failed to marshal nutritional summary: %w", err)
	}

	_, err = s.repo.CreateNutritionistSelection(date, menuItemIDs, selectedIndices, response.Reasoning, string(summaryJSON))
	return err
}