package provider

import (
	"context"
	"fmt"
	"time"
	"net/http"
	"encoding/json"
	"bytes"

	"github.com/google/uuid"
	"github.com/todorov_want/Desktop/HSE/FinancialAdviser/Backend/brain-service/internal/config"
	"github.com/todorov_want/Desktop/HSE/FinancialAdviser/Backend/brain-service/internal/entity"
)

type LocalProvider struct {
	client      *http.Client
	config      config.LLMConfig
	baseURL     string
	maxAttempts int
}

type localRequest struct {
	Messages    []localMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens"`
	Temperature float32       `json:"temperature"`
}

type localMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type localResponse struct {
	Content string `json:"content"`
	Error   string `json:"error,omitempty"`
}

func NewLocalProvider(cfg config.LLMConfig) *LocalProvider {
	return &LocalProvider{
		client:      &http.Client{Timeout: 30 * time.Second},
		config:      cfg,
		baseURL:     "http://localhost:8000", // Default local model server
		maxAttempts: 3,
	}
}

func (p *LocalProvider) GenerateResponse(ctx context.Context, messages []entity.Message) (*entity.Message, error) {
	// Convert messages to local format
	localMessages := make([]localMessage, len(messages))
	for i, msg := range messages {
		localMessages[i] = localMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
		}
	}

	// Prepare request
	req := localRequest{
		Messages:    localMessages,
		MaxTokens:   p.config.MaxTokens,
		Temperature: p.config.Temperature,
	}

	// Marshal request body
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/generate", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// Make request
	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var localResp localResponse
	if err := json.NewDecoder(resp.Body).Decode(&localResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check for error
	if localResp.Error != "" {
		return nil, fmt.Errorf("local model error: %s", localResp.Error)
	}

	// Return response in our format
	return &entity.Message{
		ID:        uuid.New().String(),
		Role:      entity.AssistantRole,
		Content:   localResp.Content,
		Status:    entity.MessageStatusComplete,
		CreatedAt: time.Now().Unix(),
	}, nil
}

func (p *LocalProvider) GetModelInfo() entity.ModelInfo {
	return entity.ModelInfo{
		Name:           p.config.Model,
		Provider:       "local",
		MaxTokens:      p.config.MaxTokens,
		Temperature:    p.config.Temperature,
		ContextWindow:  4096,  // Adjust based on your local model
		TokensPerMin:   6000,  // Adjust based on your hardware
		RequestsPerMin: 60,    // Adjust based on your hardware
	}
}