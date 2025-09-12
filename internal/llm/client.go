// internal/llm/client.go
package llm

import (
	"context"
	"net/http"

	"github.com/miftahulmahfuzh/lunch-delivery/internal/config"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type Client struct {
	llm llms.Model
}

func NewClient(cfg *config.Config) (*Client, error) {
	var llm llms.Model
	var err error

	if cfg.DeepseekTencentAPIKey == "" {
		return nil, err
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

func (c *Client) GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error) {
	return c.llm.GenerateContent(ctx, messages, options...)
}