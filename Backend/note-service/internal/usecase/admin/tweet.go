package admin

import (
	"context"
	"fmt"
	"time"

	"github.com/Denterry/FinancialAdviser/Backend/x-service/internal/entity"
	"github.com/Denterry/FinancialAdviser/Backend/x-service/internal/repo"
	"github.com/google/uuid"
)

// UseCase implements AdminTweet
type UseCase struct {
	repo repo.TweetRepository
}

// New creates a new AdminTweet use case
func New(repo repo.TweetRepository) *UseCase {
	return &UseCase{repo: repo}
}

// Create creates a new tweet
func (uc *UseCase) Create(ctx context.Context, text, authorID string) (*entity.Tweet, error) {
	t, err := entity.NewTweet(text, authorID)
	if err != nil {
		return nil, fmt.Errorf("entity.NewTweet(): %w", err)
	}

	if err := uc.repo.Create(ctx, t); err != nil {
		return nil, fmt.Errorf("uc.repo.Create(): %w", err)
	}

	return t, nil
}

// Get retrieves a tweet by its ID
func (uc *UseCase) Get(ctx context.Context, id uuid.UUID) (*entity.Tweet, error) {
	t, err := uc.repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("repo.Get(): %w", err)
	}

	return t, nil
}

// List returns all tweets matching the given filter
func (uc *UseCase) List(ctx context.Context, f repo.TweetFilter) ([]*entity.Tweet, error) {
	list, err := uc.repo.List(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("repo.List(): %w", err)
	}

	return list, nil
}

// GetTweetsBySymbol returns all tweets matching the given symbol
func (uc *UseCase) GetTweetsBySymbol(ctx context.Context, symbol string, limit, offset int32) ([]*entity.Tweet, error) {
	list, err := uc.repo.ListBySymbol(ctx, symbol, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("repo.ListBySymbol(): %w", err)
	}

	return list, nil
}

// GetTweetsBySentiment returns all tweets matching the given sentiment
func (uc *UseCase) GetTweetsBySentiment(ctx context.Context, sentiment string, limit, offset int32) ([]*entity.Tweet, error) {
	list, err := uc.repo.ListBySentiment(ctx, sentiment, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("repo.ListBySentiment(): %w", err)
	}

	return list, nil
}

// Update updates a tweet
func (uc *UseCase) Update(ctx context.Context, t *entity.Tweet) error {
	if err := uc.repo.Update(ctx, t); err != nil {
		return fmt.Errorf("uc.repo.Update(): %w", err)
	}

	return nil
}

// Delete deletes a tweet
func (uc *UseCase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("uc.repo.Delete(): %w", err)
	}

	return nil
}

// AddSymbol adds a symbol to a tweet
func (uc *UseCase) AddSymbol(ctx context.Context, id uuid.UUID, symbol string) error {
	t, err := uc.repo.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("uc.repo.Get(): %w", err)
	}

	t.AddSymbol(symbol, time.Now().UTC())
	return uc.repo.Update(ctx, t)
}

// UpdateSentiment updates the sentiment of a tweet
func (uc *UseCase) UpdateSentiment(ctx context.Context, id uuid.UUID, score float64, label string) error {
	t, err := uc.repo.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("uc.repo.Get(): %w", err)
	}

	t.UpdateSentiment(score, label, time.Now().UTC())
	return uc.repo.Update(ctx, t)
}

// UpdateEngagement updates the engagement of a tweet
func (uc *UseCase) UpdateEngagement(ctx context.Context, id uuid.UUID, likes, replies, retweets, views int) error {
	t, err := uc.repo.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("uc.repo.Get(): %w", err)
	}

	t.UpdateEngagement(likes, replies, retweets, views, time.Now().UTC())
	return uc.repo.Update(ctx, t)
}
