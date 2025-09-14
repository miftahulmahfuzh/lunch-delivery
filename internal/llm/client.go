// internal/llm/client.go
package llm

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/miftahulmahfuzh/lunch-delivery/internal/config"
	"github.com/miftahulmahfuzh/lunch-delivery/internal/interfaces"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

// Compile-time check to ensure Client implements LLMClientInterface
var _ interfaces.LLMClientInterface = (*Client)(nil)

type Client struct {
	llm llms.Model
}

func NewClient(cfg *config.Config) (*Client, error) {
	var llm llms.Model
	var err error

	if cfg.DeepseekTencentAPIKey == "" {
		return nil, fmt.Errorf("DeepseekTencentAPIKey is required")
	}

	httpClient := &http.Client{
		Timeout: cfg.LLMRequestTimeout,
	}

	opts := []openai.Option{
		openai.WithModel(cfg.DeepseekTencentModel),
		openai.WithToken(cfg.DeepseekTencentAPIKey),
		openai.WithHTTPClient(httpClient),
		openai.WithBaseURL(cfg.DeepseekTencentBaseURL),
	}

	llm, err = openai.New(opts...)
	if err != nil {
		return nil, err
	}

	return &Client{llm: llm}, nil
}

func (c *Client) GenerateContentRaw(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error) {
	return c.llm.GenerateContent(ctx, messages, options...)
}

// GenerateContent implements the LLMClientInterface
func (c *Client) GenerateContent(systemPrompt, userPrompt, temperature string) (string, error) {
	ctx := context.Background()

	// Parse temperature
	temp := 0.7 // default
	if temperature != "" {
		if t, err := strconv.ParseFloat(temperature, 64); err == nil {
			temp = t
		}
	}

	messages := []llms.MessageContent{
		{
			Role: llms.ChatMessageTypeSystem,
			Parts: []llms.ContentPart{
				llms.TextPart(systemPrompt),
			},
		},
		{
			Role: llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{
				llms.TextPart(userPrompt),
			},
		},
	}

	options := []llms.CallOption{
		llms.WithTemperature(temp),
	}

	response, err := c.llm.GenerateContent(ctx, messages, options...)
	if err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", nil
	}

	return response.Choices[0].Content, nil
}