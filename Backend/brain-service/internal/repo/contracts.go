// Package repo implements application outer layer logic. Each logic group in own file.
package repo

import (
	"context"
	"time"

	"github.com/Denterry/FinancialAdviser/Backend/x-service/internal/entity"
	"github.com/google/uuid"
)

//go:generate mockgen -source=contracts.go -destination=../usecase/mocks_repo_test.go -package=usecase_test

type (
	TweetRepository interface {
		// Create creates a new tweet
		Create(context.Context, *entity.Tweet) error
		// Get fetches tweet by ID
		Get(context.Context, uuid.UUID) (*entity.Tweet, error)
		// Update updates a tweet
		Update(context.Context, *entity.Tweet) error
		// Delete removes a tweet
		Delete(context.Context, uuid.UUID) error
		// List returns a list of tweets
		List(context.Context, TweetFilter) ([]*entity.Tweet, error)
		// ListBySymbol returns a list of tweets by symbol
		ListBySymbol(context.Context, string, int32, int32) ([]*entity.Tweet, error)
		// ListBySentiment returns a list of tweets by sentiment
		ListBySentiment(context.Context, string, int32, int32) ([]*entity.Tweet, error)
	}

	TweetProvider interface {
		// Search runs the query and returns up to maxResults posts
		Search(context.Context, string, int) ([]*entity.Tweet, error)
	}

	// TweetFilter represents filtering options for tweet queries
	TweetFilter struct {
		AuthorID       string
		IsFinancial    *bool
		SentimentLabel string
		Symbols        []string
		StartTime      *time.Time
		EndTime        *time.Time
		Limit, Offset  int32
	}
)

type (
	SocialFetcher interface {
		// SearchTweets runs the query and returns up to maxResults posts
		SearchTweets(ctx context.Context, query string, maxResults int) ([]*entity.Tweet, error)
	}
)
