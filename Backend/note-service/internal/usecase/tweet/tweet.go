package tweet

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Denterry/FinancialAdviser/Backend/x-service/internal/entity"
	"github.com/Denterry/FinancialAdviser/Backend/x-service/internal/repo"
	"github.com/google/uuid"
)

// UseCase represents the Tweet use case
type UseCase struct {
	tweetRepo repo.TweetRepository
	fetcher   repo.SocialFetcher
}

// New creates a new Tweet use case
func New(
	tweetRepo repo.TweetRepository,
	fetcher repo.SocialFetcher,
) *UseCase {
	return &UseCase{
		tweetRepo: tweetRepo,
		fetcher:   fetcher,
	}
}

// Ingest searches via the configured fetcher, persists each new tweet,
// and returns the slice of tweets that were successfully inserted
func (uc *UseCase) Ingest(ctx context.Context, query string, maxResults int) ([]*entity.Tweet, error) {
	// 1) fetch from scraper or API
	fresh, err := uc.fetcher.SearchTweets(ctx, query, maxResults)
	if err != nil {
		return nil, fmt.Errorf("uc.fetcher.SearchTweets(): %w", err)
	}

	saved := make([]*entity.Tweet, 0, len(fresh))
	now := time.Now().UTC()

	// 2) persist each one
	for _, t := range fresh {
		t.FetchedAt = now
		t.UpdatedAt = now

		if err := uc.tweetRepo.Create(ctx, t); err != nil {
			if errors.Is(err, repo.ErrDuplicateTweet) {
				continue // skip duplicate
			}
			return nil, fmt.Errorf("uc.tweetRepo.Create(): %w", err)
		}
		saved = append(saved, t)
	}

	return saved, nil
}

// GetListLatest returns the N most recent tweets ordered by created_at desc
func (uc *UseCase) GetListLatest(ctx context.Context, limit int32) ([]*entity.Tweet, error) {
	filter := repo.TweetFilter{
		Limit:  limit,
		Offset: 0,
	}

	out, err := uc.tweetRepo.List(ctx, filter) // default order by created_at desc
	if err != nil {
		return nil, fmt.Errorf("uc.tweetRepo.List(): %w", err)
	}

	return out, nil
}

// GetByID looks up a tweet by its UUID
func (uc *UseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Tweet, error) {
	t, err := uc.tweetRepo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("uc.tweetRepo.Get(): %w", err)
	}

	return t, nil
}
