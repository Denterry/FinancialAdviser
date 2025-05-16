package enricher

import (
	"context"
	"fmt"
	"time"

	"github.com/todorov_want/Desktop/HSE/FinancialAdviser/Backend/brain-service/internal/entity"
	"github.com/todorov_want/Desktop/HSE/FinancialAdviser/Backend/brain-service/internal/usecase"
)

// ContextEnricher enriches conversation context with relevant information
type ContextEnricher struct {
	tweetUseCase usecase.TweetUseCase
	maxTweets    int32
}

// NewContextEnricher creates a new context enricher
func NewContextEnricher(tweetUseCase usecase.TweetUseCase) *ContextEnricher {
	return &ContextEnricher{
		tweetUseCase: tweetUseCase,
		maxTweets:    10,
	}
}

// EnrichContext enriches the conversation with relevant context
func (e *ContextEnricher) EnrichContext(ctx context.Context, conv *entity.Conversation) error {
	// Get latest tweets for context
	tweets, err := e.tweetUseCase.GetListLatest(ctx, e.maxTweets)
	if err != nil {
		return fmt.Errorf("failed to get latest tweets: %w", err)
	}

	// Create a system message with enriched context
	contextMsg := entity.Message{
		ID:        uuid.New().String(),
		Role:      entity.SystemRole,
		Content:   e.formatTweetsContext(tweets),
		Status:    entity.MessageStatusComplete,
		CreatedAt: time.Now().Unix(),
	}

	// Add context message at the beginning of conversation
	conv.Messages = append([]entity.Message{contextMsg}, conv.Messages...)

	return nil
}

// formatTweetsContext formats tweets into a context message
func (e *ContextEnricher) formatTweetsContext(tweets []*entity.Tweet) string {
	var context string
	context = "Here are the latest relevant tweets for context:\n\n"

	for _, tweet := range tweets {
		context += fmt.Sprintf("Tweet from %s:\n%s\nSentiment: %f\n\n",
			tweet.AuthorID,
			tweet.Text,
			tweet.SentimentScore)
	}

	return context
}