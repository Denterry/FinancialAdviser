package webapi

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Denterry/FinancialAdviser/Backend/x-service/config"
	"github.com/Denterry/FinancialAdviser/Backend/x-service/internal/entity"
	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/google/uuid"
)

const (
	twitterSearchEndpoint = "https://api.twitter.com/2/tweets/search/recent"
)

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

// TwitterAPI implements repo.SocialFetcher using go-twitter/v2
type TwitterAPI struct {
	client *twitter.Client
}

// NewTwitterAPI constructs a TwitterAPI using the given config
func NewTwitterAPI(cfg config.XProvider) (*TwitterAPI, error) {
	if cfg.XAPI.BearerToken == "" {
		return nil, fmt.Errorf("X_API_BEARER_TOKEN must be set for api mode")
	}

	// src := oauth2.StaticTokenSource(&oauth2.Token{
	// 	AccessToken: cfg.XAPI.BearerToken,
	// })
	// httpClient := oauth2.NewClient(context.Background(), src)

	client := &twitter.Client{
		Authorizer: authorize{
			Token: cfg.XAPI.BearerToken,
		},
		Client: http.DefaultClient,
		Host:   cfg.XAPI.BaseURL,
	}

	return &TwitterAPI{client: client}, nil
}

// SearchTweets performs a recent search via Twitter API
func (api *TwitterAPI) SearchTweets(ctx context.Context, query string, max int) ([]*entity.Tweet, error) {
	opts := twitter.TweetRecentSearchOpts{
		MaxResults: max,
		Expansions: []twitter.Expansion{
			twitter.ExpansionAuthorID,
		},
		TweetFields: []twitter.TweetField{
			twitter.TweetFieldCreatedAt,
			twitter.TweetFieldLanguage,
			twitter.TweetFieldPublicMetrics,
			twitter.TweetFieldEntities,
		},
		UserFields: []twitter.UserField{
			twitter.UserFieldUserName,
		},
	}

	resp, err := api.client.TweetRecentSearch(ctx, query, opts)
	if err != nil {
		return nil, fmt.Errorf("TweetRecentSearch error: %w", err)
	}
	if resp.RateLimit != nil {
		fmt.Printf("Rate limit: %d remaining, resets at %v\n", resp.RateLimit.Remaining, resp.RateLimit.Reset)
	}

	userMap := make(map[string]string)
	for _, user := range resp.Raw.Includes.Users {
		userMap[user.ID] = user.UserName
	}

	now := time.Now().UTC()
	var results []*entity.Tweet

	for _, t := range resp.Raw.Tweets {
		var urls []string
		if t.Entities != nil {
			for _, u := range t.Entities.URLs {
				if u.ExpandedURL != "" {
					urls = append(urls, u.ExpandedURL)
				}
			}
		}

		createdAt, err := time.Parse(time.RFC3339, t.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("time.Parse(): %w", err)
		}

		userName := userMap[t.AuthorID]
		if userName == "" {
			userName = "unknown"
		}

		symbols := extractSymbols(t.Text, nil)
		isFinancial := len(symbols) > 0

		results = append(results, &entity.Tweet{
			ID:          uuid.MustParse(t.ID),
			Text:        t.Text,
			Lang:        strings.ToLower(t.Language),
			AuthorID:    t.AuthorID,
			UserName:    userName,
			CreatedAt:   createdAt,
			FetchedAt:   now,
			UpdatedAt:   now,
			Likes:       t.PublicMetrics.Likes,
			Replies:     t.PublicMetrics.Replies,
			Retweets:    t.PublicMetrics.Retweets,
			Views:       0,
			URLs:        urls,
			Photos:      nil,
			Videos:      nil,
			IsFinancial: isFinancial,
			Symbols:     nil,
		})
	}

	return results, nil
}
