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
	} else if cfg.XScraper.Username != "" {
		if err := s.Login(cfg.XScraper.Username, cfg.XScraper.Password); err != nil {
			return nil, fmt.Errorf("s.Login(): %w", err)
		}
	}

	return &TwitterScraper{s}, nil
}

// SearchTweets runs the query and returns up to maxTweets tweets
func (ts *TwitterScraper) SearchTweets(ctx context.Context, query string, maxTweets int) ([]*entity.Tweet, error) {
	result := make([]*entity.Tweet, 0, maxTweets)
	for tr := range ts.scraper.SearchTweets(ctx, query, maxTweets) {
		if tr == nil || tr.Error != nil {
			return nil, fmt.Errorf("ts.scraper.SearchTweets(): %w", tr.Error)
		}
		result = append(result, mapScrapedTweet(&tr.Tweet))
	}

	return result, nil
}
