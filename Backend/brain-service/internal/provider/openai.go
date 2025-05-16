package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/sashabaranov/go-openai"
	"github.com/todorov_want/Desktop/HSE/FinancialAdviser/Backend/brain-service/internal/config"
	"github.com/todorov_want/Desktop/HSE/FinancialAdviser/Backend/brain-service/internal/entity"
)

type OpenAIProvider struct {
	client      *openai.Client
	config      config.LLMConfig
	maxAttempts int
}

func NewOpenAIProvider(cfg config.LLMConfig) *OpenAIProvider {
	client := openai.NewClient(cfg.APIKey)
	return &OpenAIProvider{
		client:      client,
		config:      cfg,
		maxAttempts: 3,
	}
}

func (p *OpenAIProvider) GenerateResponse(ctx context.Context, messages []entity.Message) (*entity.Message, error) {
	openaiMessages := make([]openai.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		openaiMessages[i] = openai.ChatCompletionMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
	}
	}

	req := openai.ChatCompletionRequest{
		Model:       p.config.Model,
		Messages:    openaiMessages,
		MaxTokens:   p.config.MaxTokens,
		Temperature: float32(p.config.Temperature),
	}

	resp, err := p.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate response: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, errors.New("no response generated")
	}

	return &entity.Message{
		Role:    entity.MessageRole(resp.Choices[0].Message.Role),
		Content: resp.Choices[0].Message.Content,
		Status:  entity.MessageStatusCompleted,
	}, nil
}

func (p *OpenAIProvider) GetModelInfo() entity.ModelInfo {
	return entity.ModelInfo{
		Name:           p.config.Model,
		Provider:       "openai",
		MaxTokens:      p.config.MaxTokens,
		Temperature:    p.config.Temperature,
		ContextWindow:  8192, // This varies by model, adjust accordingly
		TokensPerMin:   90000, // Adjust based on your rate limits
		RequestsPerMin: 3500,  // Adjust based on your rate limits
	}
}