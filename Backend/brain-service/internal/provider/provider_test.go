package provider

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenAIProvider(t *testing.T) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		t.Skip("OPENAI_API_KEY not set")
	}

	cfg := &Config{
		APIKey:      apiKey,
		Model:       "gpt-3.5-turbo",
		MaxTokens:   1000,
		Temperature: 0.7,
	}

	provider := NewOpenAIProvider(cfg)

	t.Run("generate response", func(t *testing.T) {
		messages := []Message{
			{
				Role:    MessageRoleUser,
				Content: "What is 2+2?",
			},
		}

		resp, err := provider.GenerateResponse(context.Background(), messages)
		require.NoError(t, err)
		assert.NotEmpty(t, resp)
	})

	t.Run("get model info", func(t *testing.T) {
		info := provider.GetModelInfo()
		assert.Equal(t, "gpt-3.5-turbo", info.Name)
		assert.True(t, info.MaxTokens > 0)
	})
}

func TestAnthropicProvider(t *testing.T) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("ANTHROPIC_API_KEY not set")
	}

	cfg := &Config{
		APIKey:      apiKey,
		Model:       "claude-2",
		MaxTokens:   1000,
		Temperature: 0.7,
	}

	provider := NewAnthropicProvider(cfg)

	t.Run("generate response", func(t *testing.T) {
		messages := []Message{
			{
				Role:    MessageRoleUser,
				Content: "What is 2+2?",
			},
		}

		resp, err := provider.GenerateResponse(context.Background(), messages)
		require.NoError(t, err)
		assert.NotEmpty(t, resp)
	})

	t.Run("get model info", func(t *testing.T) {
		info := provider.GetModelInfo()
		assert.Equal(t, "claude-2", info.Name)
		assert.True(t, info.MaxTokens > 0)
	})
}

func TestLocalProvider(t *testing.T) {
	cfg := &Config{
		Endpoint:    "http://localhost:8080",
		Model:       "local-model",
		MaxTokens:   1000,
		Temperature: 0.7,
	}

	provider := NewLocalProvider(cfg)

	t.Run("generate response", func(t *testing.T) {
		messages := []Message{
			{
				Role:    MessageRoleUser,
				Content: "What is 2+2?",
			},
		}

		resp, err := provider.GenerateResponse(context.Background(), messages)
		require.NoError(t, err)
		assert.NotEmpty(t, resp)
	})

	t.Run("get model info", func(t *testing.T) {
		info := provider.GetModelInfo()
		assert.Equal(t, "local-model", info.Name)
		assert.True(t, info.MaxTokens > 0)
	})
}

func TestProviderFactory(t *testing.T) {
	t.Run("create openai provider", func(t *testing.T) {
		cfg := &Config{
			Provider:    ModelProviderOpenAI,
			APIKey:      "test-key",
			Model:       "gpt-3.5-turbo",
			MaxTokens:   1000,
			Temperature: 0.7,
		}

		factory := NewFactory()
		provider, err := factory.CreateProvider(cfg)
		require.NoError(t, err)
		assert.IsType(t, &OpenAIProvider{}, provider)
	})

	t.Run("create anthropic provider", func(t *testing.T) {
		cfg := &Config{
			Provider:    ModelProviderAnthropic,
			APIKey:      "test-key",
			Model:       "claude-2",
			MaxTokens:   1000,
			Temperature: 0.7,
		}

		factory := NewFactory()
		provider, err := factory.CreateProvider(cfg)
		require.NoError(t, err)
		assert.IsType(t, &AnthropicProvider{}, provider)
	})

	t.Run("create local provider", func(t *testing.T) {
		cfg := &Config{
			Provider:    ModelProviderLocal,
			Endpoint:    "http://localhost:8080",
			Model:       "local-model",
			MaxTokens:   1000,
			Temperature: 0.7,
		}

		factory := NewFactory()
		provider, err := factory.CreateProvider(cfg)
		require.NoError(t, err)
		assert.IsType(t, &LocalProvider{}, provider)
	})

	t.Run("invalid provider", func(t *testing.T) {
		cfg := &Config{
			Provider: "invalid",
		}

		factory := NewFactory()
		_, err := factory.CreateProvider(cfg)
		assert.Error(t, err)
	})
}