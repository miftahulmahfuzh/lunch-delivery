package interfaces

// LLMClientInterface defines the interface for LLM client operations
type LLMClientInterface interface {
	GenerateContent(systemPrompt, userPrompt, temperature string) (string, error)
}