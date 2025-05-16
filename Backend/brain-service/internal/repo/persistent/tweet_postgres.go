package persistent

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Denterry/FinancialAdviser/Backend/x-service/internal/entity"
	"github.com/Denterry/FinancialAdviser/Backend/x-service/internal/repo"
	"github.com/Denterry/FinancialAdviser/Backend/x-service/pkg/postgres"
	"github.com/google/uuid"
)

const tweetsTbl = "tweets"

// TweetRepository implements repo.TweetRepository backed by Postgres
type TweetRepository struct {
	*postgres.Postgres
}

// NewTweetPostgres returns TweetRepository
func NewTweetPostgres(pg *postgres.Postgres) *TweetRepository {
	return &TweetRepository{pg}
}

// Create inserts tweet + optional symbol links in one tx
func (r *TweetRepository) Create(ctx context.Context, t *entity.Tweet) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	now := time.Now().UTC()
	t.FetchedAt, t.UpdatedAt = now, now

	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("r.Pool.Begin(): %w", err)
	}
	defer tx.Rollback(ctx)

	const queryTweets = ` -- Create(ctx context.Context, t *entity.Tweet) error 
		INSERT INTO tweets (
			id, text, lang, author_id, username,
			created_at, fetched_at, updated_at,
			likes, replies, retweets, views,
			urls, photos, videos,
			is_financial, sentiment_score, sentiment_label,
			raw_payload
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, '{}')`

	_, err = tx.Exec(ctx, queryTweets,
		t.ID, t.Text, t.Lang, t.AuthorID, t.UserName,
		t.CreatedAt, t.FetchedAt, t.UpdatedAt,
		t.Likes, t.Replies, t.Retweets, t.Views,
		t.URLs, t.Photos, t.Videos,
		t.IsFinancial, t.SentimentScore, t.SentimentLabel,
	)
	if err != nil {
		return fmt.Errorf("tx.Exec(): %w", err)
	}

	if len(t.Symbols) > 0 {
		const querySymbols = `-- Create(ctx context.Context, t *entity.Tweet) error 
			INSERT INTO symbols (ticker, type, display_name) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING
		`

		const queryTweetSymbols = `-- Create(ctx context.Context, t *entity.Tweet) error 
			INSERT INTO tweet_symbols (tweet_id, symbol) VALUES ($1, $2)
		`

		for _, symbol := range t.Symbols {
			s := strings.ToUpper(symbol)

			_, err = tx.Exec(ctx, querySymbols, s, entity.SymbolTypeEquity, s)
			if err != nil {
				return fmt.Errorf("tx.Exec(symbols): %w", err)
			}

			_, err = tx.Exec(ctx, queryTweetSymbols, t.ID, s)
			if err != nil {
				return fmt.Errorf("tx.Exec(tweet_symbols): %w", err)
			}
		}
	}

	return tx.Commit(ctx)
}

// Get fetches tweet by ID
func (r *TweetRepository) Get(ctx context.Context, id uuid.UUID) (*entity.Tweet, error) {
	const query = ` -- Get(ctx context.Context, id uuid.UUID) (*entity.Tweet, error) 
		SELECT 
			id, text, lang, author_id, username,
			created_at, fetched_at, updated_at,
			likes, replies, retweets, views,
			urls, photos, videos,
			is_financial, sentiment_score, sentiment_label
		FROM tweets 
		WHERE id = $1`

	row := r.Pool.QueryRow(ctx, query, id)
	return scanTweet(row)
}

// Delete removes tweet (cascade cleans tweet_symbols via FK)
func (r *TweetRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = ` -- Delete(ctx context.Context, id uuid.UUID) error 
		DELETE FROM tweets WHERE id = $1
	`
	_, err := r.Pool.Exec(ctx, query, id)
	return err
}

// Update updates basic editable fields + engagement / sentiment
func (r *TweetRepository) Update(ctx context.Context, t *entity.Tweet) error {
	const query = ` -- Update(ctx context.Context, t *entity.Tweet) error 
		UPDATE tweets 
		SET 
			text = $1, likes = $2, replies = $3, retweets = $4, views = $5,
			sentiment_score = $6, sentiment_label = $7,
			updated_at = $8
		WHERE id = $9`

	_, err := r.Pool.Exec(ctx, query,
		t.Text, t.Likes, t.Replies, t.Retweets, t.Views,
		t.SentimentScore, t.SentimentLabel,
		time.Now().UTC(), t.ID,
	)
	return err
}

// List returns tweets by various optional filters
func (r *TweetRepository) List(ctx context.Context, f repo.TweetFilter) ([]*entity.Tweet, error) {
	const query = ` -- List(ctx context.Context, f repo.TweetFilter) ([]*entity.Tweet, error) 
		SELECT 
			id, text, lang, author_id, username,
			created_at, fetched_at, updated_at,
			likes, replies, retweets, views,
			urls, photos, videos,
			is_financial, sentiment_score, sentiment_label
		FROM tweets
	`

	sqlSuffix, args := buildFilter(f)
	rows, err := r.Pool.Query(ctx, query+sqlSuffix, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*entity.Tweet
	for rows.Next() {
		t, err := scanTweet(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

// ListBySymbol quickly fetches tweets linked with given ticker
func (r *TweetRepository) ListBySymbol(ctx context.Context, symbol string, limit, offset int32) ([]*entity.Tweet, error) {
	query := `
		SELECT
			t.id, t.text, t.lang, t.author_id, t.username,
			t.created_at, t.fetched_at, t.updated_at,
			t.likes, t.replies, t.retweets, t.views,
			t.urls, t.photos, t.videos,
			t.is_financial, t.sentiment_score, t.sentiment_label
		FROM tweets t
		JOIN tweet_symbols ts ON ts.tweet_id = t.id
		WHERE ts.symbol = $1
		ORDER BY t.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.Pool.Query(ctx, query, strings.ToUpper(symbol), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []*entity.Tweet
	for rows.Next() {
		t, err := scanTweet(rows)
		if err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	return res, rows.Err()
}

// ListBySentiment returns tweets by sentiment label.
func (r *TweetRepository) ListBySentiment(ctx context.Context, label string, limit, offset int32) ([]*entity.Tweet, error) {
	query := `
		SELECT
			id, text, lang, author_id, username,
			created_at, fetched_at, updated_at,
			likes, replies, retweets, views,
			urls, photos, videos,
			is_financial, sentiment_score, sentiment_label
		FROM tweets
		WHERE sentiment_label = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.Pool.Query(ctx, query, label, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []*entity.Tweet
	for rows.Next() {
		t, err := scanTweet(rows)
		if err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	return res, rows.Err()
}
