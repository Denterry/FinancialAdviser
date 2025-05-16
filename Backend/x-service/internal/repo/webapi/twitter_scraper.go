package webapi

import (
	"context"
	"fmt"

	"github.com/Denterry/FinancialAdviser/Backend/x-service/config"
	"github.com/Denterry/FinancialAdviser/Backend/x-service/internal/entity"
	twitterscraper "github.com/n0madic/twitter-scraper"
)

// TwitterScraper implements scraper for Twitter
type TwitterScraper struct {
	scraper *twitterscraper.Scraper
}

// NewTwitterScraper returns a new TwitterScraper
func NewTwitterScraper(cfg config.XProvider) (*TwitterScraper, error) {
	s := twitterscraper.New().WithDelay(int64(cfg.XScraper.DelaySec))
	s.WithReplies(cfg.XScraper.IncludeReply)
	s.SetSearchMode(twitterscraper.SearchLatest)

	if cfg.XScraper.Proxy != "" {
		if err := s.SetProxy(cfg.XScraper.Proxy); err != nil {
			return nil, fmt.Errorf("s.SetProxy(): %w", err)
		}
	}

	if cfg.XScraper.UseAppLogin {
		if err := s.LoginOpenAccount(); err != nil {
			return nil, fmt.Errorf("s.LoginOpenAccount(): %w", err)
		}
	} else if cfg.XScraper.Username != "" && cfg.XScraper.Password != "" {
		if err := s.Login(cfg.XScraper.Username, cfg.XScraper.Password); err != nil {
			return nil, fmt.Errorf("s.Login(): %w", err)
		}
	} else {
		return nil, fmt.Errorf("no valid login method provided: set UseAppLogin=true or provide username/password")
	}

	if !s.IsLoggedIn() {
		return nil, fmt.Errorf("s.IsLoggedIn(): failed to login")
	}

	return &TwitterScraper{s}, nil
}

// SearchTweets runs the query and returns up to maxTweets tweets
func (ts *TwitterScraper) SearchTweets(ctx context.Context, query string, maxTweets int) ([]*entity.Tweet, error) {
	result := make([]*entity.Tweet, 0, maxTweets)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case tr, ok := <-ts.scraper.SearchTweets(ctx, query, maxTweets):
			if !ok {
				return result, nil
			}
			if tr == nil || tr.Error != nil {
				id := "<nil>"
				if tr != nil && tr.Tweet.ID != "" {
					id = tr.Tweet.ID
				}
				return nil, fmt.Errorf("tweet scrape failed (id=%s): %w", id, tr.Error)
			}
			result = append(result, mapScrapedTweet(&tr.Tweet))
		}
	}
}
