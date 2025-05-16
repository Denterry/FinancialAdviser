package usecase

import (
	"context"

	"github.com/Denterry/FinancialAdviser/Backend/x-service/internal/entity"
	"github.com/Denterry/FinancialAdviser/Backend/x-service/internal/repo"
	"github.com/google/uuid"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_usecase_test.go -package=usecase_test

type (
	TweetUseCase interface {
		// Ingest - fetches fresh Tweets that match the query, stores them,
		// and returns the slice that were persisted this round
		Ingest(context.Context, string, int) ([]*entity.Tweet, error)

		// GetListLatest - returns newest stored tweets
		GetListLatest(context.Context, int32) ([]*entity.Tweet, error)

		// GetByID - returns a single stored tweet by ID
		GetByID(ctx context.Context, id uuid.UUID) (*entity.Tweet, error)
	}
)

type (
	AdminTweetUseCase interface {
		// Create - creates a new tweet
		Create(ctx context.Context, text, authorID string) (*entity.Tweet, error)

		// Get - retrieves a tweet by its ID
		Get(ctx context.Context, id uuid.UUID) (*entity.Tweet, error)

		// List - retrieves all tweets matching the given filter
		List(ctx context.Context, f repo.TweetFilter) ([]*entity.Tweet, error)

		// GetTweetsBySymbol - retrieves all tweets matching the given symbol
		GetTweetsBySymbol(ctx context.Context, symbol string, limit, offset int32) ([]*entity.Tweet, error)

		// GetTweetsBySentiment - retrieves all tweets matching the given sentiment
		GetTweetsBySentiment(ctx context.Context, sentiment string, limit, offset int32) ([]*entity.Tweet, error)

		// Update - updates a tweet
		Update(ctx context.Context, t *entity.Tweet) error

		// Delete - deletes a tweet
		Delete(ctx context.Context, id uuid.UUID) error

		// AddSymbol - adds a symbol to a tweet
		AddSymbol(ctx context.Context, id uuid.UUID, symbol string) error

		// UpdateSentiment - updates the sentiment of a tweet
		UpdateSentiment(ctx context.Context, id uuid.UUID, score float64, label string) error

		// UpdateEngagement - updates the engagement of a tweet
		UpdateEngagement(ctx context.Context, id uuid.UUID, likes, replies, retweets, views int) error
	}
)
