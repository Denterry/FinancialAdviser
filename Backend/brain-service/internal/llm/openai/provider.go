package openai

import (
	"context"
	"fmt"

	"github.com/Denterry/FinancialAdviser/Backend/brain-service/internal/entity"
	"github.com/sashabaranov/go-openai"
)

type Config struct {
	APIKey      string
	Model       string
	MaxTokens   int
	Temperature float32
}

type Provider struct {
	client *openai.Client
	config Config
}

func NewProvider(config Config) *Provider {
	return &Provider{
		client: openai.NewClient(config.APIKey),
		config: config,
	}
}

func (p *Provider) GenerateResponse(ctx context.Context, messages []entity.Message) (string, error) {
	// Convert our messages to OpenAI format
	openaiMessages := make([]openai.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		openaiMessages[i] = openai.ChatCompletionMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
		}
	}

	// Create completion request
	req := openai.ChatCompletionRequest{
		Model:       p.config.Model,
		Messages:    openaiMessages,
		MaxTokens:   p.config.MaxTokens,
		Temperature: p.config.Temperature,
	}

	// Get completion
	resp, err := p.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to create chat completion: %v", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no completion choices returned")
	}

	return resp.Choices[0].Message.Content, nil
}

func (p *Provider) GetModelInfo() map[string]interface{} {
	return map[string]interface{}{
		"provider": "openai",
		"model":    p.config.Model,
		"settings": map[string]interface{}{
			"max_tokens":   p.config.MaxTokens,
			"temperature":  p.config.Temperature,
		},
	}
}