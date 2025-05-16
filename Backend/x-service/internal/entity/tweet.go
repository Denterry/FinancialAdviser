package entity

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type TweetOption func(*Tweet)

// Tweet represents a Twitter post with its metadata and content
type Tweet struct {
	ID   uuid.UUID `db:"id" json:"id"`
	Text string    `db:"text" json:"text"`
	Lang string    `db:"lang" json:"lang"` // ISO-language code ("en")

	AuthorID string `db:"author_id" json:"author_id"`
	UserName string `db:"username" json:"username"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	FetchedAt time.Time `db:"fetched_at" json:"fetched_at"` // when scraped / ingested it
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"` // last update in DB

	// engagement metrics – point-in-time snapshot
	Likes    int `db:"likes" json:"likes"`
	Replies  int `db:"replies" json:"replies"`
	Retweets int `db:"retweets" json:"retweets"`
	Views    int `db:"views" json:"views"`

	// media / entities
	URLs   []string `db:"urls" json:"urls"`     // string array in Postgres
	Photos []string `db:"photos" json:"photos"` // ditto
	Videos []string `db:"videos" json:"videos"`

	// financial enrichment
	IsFinancial    bool     `db:"is_financial" json:"is_financial"`
	Symbols        []string `db:"symbols" json:"symbols"`                 // ["TSLA", "BTC", ...]
	SentimentScore float64  `db:"sentiment_score" json:"sentiment_score"` // range –1 .. 1
	SentimentLabel string   `db:"sentiment_label" json:"sentiment_label"` // "POS", "NEG", "NEU"
}

// NewTweet creates a new Tweet instance
func NewTweet(text, authorID string, opts ...TweetOption) (*Tweet, error) {
	now := time.Now().UTC()
	t := &Tweet{
		ID:        uuid.New(),
		Text:      text,
		AuthorID:  authorID,
		CreatedAt: now,
		FetchedAt: now,
		UpdatedAt: now,
	}

	for _, opt := range opts {
		opt(t)
	}

	if err := t.Validate(); err != nil {
		return nil, err
	}
	return t, nil
}

// Validate checks if the tweet is valid
func (t *Tweet) Validate() error {
	switch {
	case strings.TrimSpace(t.Text) == "":
		return ErrEmptyText
	case len(t.Text) > 4096:
		return ErrTooLongText
	}
	return nil
}

// Touch updates the updated_at field to the current time
func (t *Tweet) Touch(now time.Time) {
	t.UpdatedAt = now.UTC()
}

// UpdateSentiment updates the sentiment analysis results
func (t *Tweet) UpdateSentiment(score float64, label string, now time.Time) {
	t.SentimentScore = score
	t.SentimentLabel = label
	t.UpdatedAt = now.UTC()
}

// UpdateEngagement updates the engagement metrics
func (t *Tweet) UpdateEngagement(likes, replies, retweets, views int, now time.Time) {
	t.Likes = likes
	t.Replies = replies
	t.Retweets = retweets
	t.Views = views
	t.UpdatedAt = now.UTC()
}

// AddSymbol adds a financial symbol to the tweet
func (t *Tweet) AddSymbol(symbol string, now time.Time) {
	t.Symbols = append(t.Symbols, symbol)
	t.IsFinancial = true
	t.UpdatedAt = now.UTC()
}

// AddURL adds a URL to the tweet
func (t *Tweet) AddURL(url string, now time.Time) {
	t.URLs = append(t.URLs, url)
	t.UpdatedAt = now.UTC()
}

// AddPhoto adds a photo to the tweet
func (t *Tweet) AddPhoto(photo string, now time.Time) {
	t.Photos = append(t.Photos, photo)
	t.UpdatedAt = now.UTC()
}

// AddVideo adds a video to the tweet
func (t *Tweet) AddVideo(video string, now time.Time) {
	t.Videos = append(t.Videos, video)
	t.UpdatedAt = now.UTC()
}
