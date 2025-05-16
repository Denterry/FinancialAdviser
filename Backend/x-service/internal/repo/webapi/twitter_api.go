package webapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/Denterry/FinancialAdviser/Backend/x-service/config"
	"github.com/Denterry/FinancialAdviser/Backend/x-service/internal/entity"
	"github.com/google/uuid"
)

const (
	twitterSearchEndpoint = "https://api.twitter.com/2/tweets/search/recent"
)

// TwitterAPI implements the official Twitter v2 "Recent Search" endpoint
type TwitterAPI struct {
	client      *http.Client
	bearerToken string
}

// NewTwitterAPI constructs a v2 client using the bearer token in cfg.XProvider.BearerToken
func NewTwitterAPI(cfg config.XProvider) (*TwitterAPI, error) {
	if cfg.XAPI.BearerToken == "" {
		return nil, fmt.Errorf("X_TOKEN must be set for api mode")
	}
	return &TwitterAPI{
		client:      &http.Client{Timeout: 10 * time.Second},
		bearerToken: cfg.XAPI.BearerToken,
	}, nil
}

// SearchTweets calls the twitter API and decodes the JSON
func (api *TwitterAPI) SearchTweets(ctx context.Context, query string, max int) ([]*entity.Tweet, error) {
	u, err := url.Parse(twitterSearchEndpoint)
	if err != nil {
		return nil, fmt.Errorf("url.Parse(): %w", err)
	}

	q := u.Query()
	q.Set("query", query)
	q.Set("max_results", fmt.Sprintf("%d", max))

	// required fields for tweet data
	q.Set("tweet.fields", "author_id,created_at,lang,public_metrics,entities")

	// required for user information
	q.Set("expansions", "author_id")
	q.Set("user.fields", "username")

	// required for media information
	q.Set("media.fields", "url,preview_image_url,type")

	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext(): %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+api.bearerToken)

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("api.client.Do(): %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll(): %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("twitter api status %d: %s", resp.StatusCode, string(body))
	}

	var wrapper struct {
		Data []struct {
			ID            string    `json:"id"`
			Text          string    `json:"text"`
			AuthorID      string    `json:"author_id"`
			CreatedAt     time.Time `json:"created_at"`
			Lang          string    `json:"lang"`
			PublicMetrics struct {
				RetweetCount int `json:"retweet_count"`
				ReplyCount   int `json:"reply_count"`
				LikeCount    int `json:"like_count"`
				QuoteCount   int `json:"quote_count"`
			} `json:"public_metrics"`
			Entities struct {
				URLs []struct {
					URL         string `json:"url"`
					ExpandedURL string `json:"expanded_url"`
				} `json:"urls"`
			} `json:"entities"`
		} `json:"data"`
		Includes struct {
			Users []struct {
				ID       string `json:"id"`
				Username string `json:"username"`
			} `json:"users"`
		} `json:"includes"`
	}

	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("json.Unmarshal(): %w", err)
	}

	userMap := make(map[string]string)
	for _, user := range wrapper.Includes.Users {
		userMap[user.ID] = user.Username
	}

	now := time.Now().UTC()
	out := make([]*entity.Tweet, 0, len(wrapper.Data))
	for _, d := range wrapper.Data {
		// extract URLs from entities
		urls := make([]string, 0, len(d.Entities.URLs))
		for _, u := range d.Entities.URLs {
			urls = append(urls, u.ExpandedURL)
		}

		t := &entity.Tweet{
			ID:        uuid.MustParse(d.ID),
			Text:      d.Text,
			Lang:      d.Lang,
			AuthorID:  d.AuthorID,
			UserName:  userMap[d.AuthorID],
			CreatedAt: d.CreatedAt,
			FetchedAt: now,
			UpdatedAt: now,
			Likes:     d.PublicMetrics.LikeCount,
			Replies:   d.PublicMetrics.ReplyCount,
			Retweets:  d.PublicMetrics.RetweetCount,
			Views:     0,
			URLs:      urls,
			Photos:    nil, // TODO: Extract from media attachments
			Videos:    nil, // TODO: Extract from media attachments
		}
		out = append(out, t)
	}
	return out, nil
}
