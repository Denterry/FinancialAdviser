package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/todorov_want/Desktop/HSE/FinancialAdviser/Backend/brain-service/internal/config"
	"github.com/todorov_want/Desktop/HSE/FinancialAdviser/Backend/brain-service/internal/entity"
	"github.com/anthropic-ai/anthropic-sdk-go"
)

type AnthropicProvider struct {
	client      *anthropic.Client
	config      config.LLMConfig
	maxAttempts int
}

func NewAnthropicProvider(cfg config.LLMConfig) *AnthropicProvider {
	client := anthropic.NewClient(cfg.APIKey)
	return &AnthropicProvider{
		client:      client,
		config:      cfg,
		maxAttempts: 3,
	}
}

func (p *AnthropicProvider) GenerateResponse(ctx context.Context, messages []entity.Message) (*entity.Message, error) {
	// Convert our messages to Anthropic format
	anthropicMessages := make([]anthropic.Message, len(messages))
	for i, msg := range messages {
		anthropicMessages[i] = anthropic.Message{
			Role:    string(msg.Role),
			Content: msg.Content,
	}
	}

	// Create completion request
	req := &anthropic.CreateMessageRequest{
		Model:       p.config.Model,
		Messages:    anthropicMessages,
		MaxTokens:   p.config.MaxTokens,
		Temperature: float32(p.config.Temperature),
	}

	// Make API call
	resp, err := p.client.CreateMessage(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate response: %w", err)
	}

	// Convert response to our format
	return &entity.Message{
		ID:        uuid.New().String(),
		Role:      entity.AssistantRole,
		Content:   resp.Content,
		Status:    entity.MessageStatusComplete,
		CreatedAt: time.Now().Unix(),
	}, nil
}

func (p *AnthropicProvider) GetModelInfo() entity.ModelInfo {
	return entity.ModelInfo{
		Name:           p.config.Model,
		Provider:       "anthropic",
		MaxTokens:      p.config.MaxTokens,
		Temperature:    p.config.Temperature,
		ContextWindow:  100000, // Claude has a large context window
		TokensPerMin:   45000,  // Adjust based on your rate limits
		RequestsPerMin: 1500,   // Adjust based on your rate limits
	}
}