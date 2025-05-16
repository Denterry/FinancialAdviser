package provider

import (
	"fmt"

	"github.com/todorov_want/Desktop/HSE/FinancialAdviser/Backend/brain-service/internal/config"
	"github.com/todorov_want/Desktop/HSE/FinancialAdviser/Backend/brain-service/internal/entity"
)

// LLMProvider defines the interface for language model providers
type LLMProvider interface {
	// GenerateResponse generates a response for the given messages
	GenerateResponse(ctx context.Context, messages []entity.Message) (*entity.Message, error)
	// GetModelInfo returns information about the model's capabilities
	GetModelInfo() entity.ModelInfo
}

// Factory creates LLM providers based on configuration
type Factory struct {
	config config.LLMConfig
}

// NewFactory creates a new provider factory
func NewFactory(cfg config.LLMConfig) *Factory {
	return &Factory{config: cfg}
}

// CreateProvider creates a new LLM provider based on the configuration
func (f *Factory) CreateProvider() (LLMProvider, error) {
	switch entity.ModelProvider(f.config.Provider) {
	case entity.OpenAIProvider:
		return NewOpenAIProvider(f.config), nil
	case entity.AnthropicProvider:
		return nil, fmt.Errorf("anthropic provider not implemented yet")
	case entity.LocalProvider:
		return nil, fmt.Errorf("local provider not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported provider: %s", f.config.Provider)
	}
}