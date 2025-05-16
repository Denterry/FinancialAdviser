package persistent

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Denterry/FinancialAdviser/Backend/x-service/internal/entity"
	"github.com/Denterry/FinancialAdviser/Backend/x-service/internal/repo"
	"github.com/jackc/pgx/v4"
)

// scanTweet scans a tweet from a database row
func scanTweet(row pgx.Row) (*entity.Tweet, error) {
	var t entity.Tweet
	err := row.Scan(
		&t.ID,
		&t.Text,
		&t.Lang,
		&t.AuthorID,
		&t.UserName,
		&t.CreatedAt,
		&t.FetchedAt,
		&t.UpdatedAt,
		&t.Likes,
		&t.Replies,
		&t.Retweets,
		&t.Views,
		&t.URLs,
		&t.Photos,
		&t.Videos,
		&t.IsFinancial,
		&t.SentimentScore,
		&t.SentimentLabel,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// buildFilter builds WHERE … LIMIT/OFFSET for List queries
func buildFilter(f repo.TweetFilter) (string, []any) {
	var (
		buf   bytes.Buffer
		args  []any
		where []string
	)

	if f.AuthorID != "" {
		args = append(args, f.AuthorID)
		where = append(where, fmt.Sprintf("author_id=$%d", len(args)))
	}
	if f.IsFinancial != nil {
		args = append(args, *f.IsFinancial)
		where = append(where, fmt.Sprintf("is_financial=$%d", len(args)))
	}
	if f.SentimentLabel != "" {
		args = append(args, f.SentimentLabel)
		where = append(where, fmt.Sprintf("sentiment_label=$%d", len(args)))
	}
	if !f.StartTime.IsZero() {
		args = append(args, f.StartTime)
		where = append(where, fmt.Sprintf("created_at>=$%d", len(args)))
	}
	if !f.EndTime.IsZero() {
		args = append(args, f.EndTime)
		where = append(where, fmt.Sprintf("created_at<=$%d", len(args)))
	}
	if len(f.Symbols) > 0 {
		args = append(args, f.Symbols)
		where = append(where,
			fmt.Sprintf(`EXISTS (SELECT 1 FROM tweet_symbols ts WHERE ts.tweet_id=tweets.id AND ts.symbol = ANY($%d))`, len(args)),
		)
	}

	if len(where) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(where, " AND "))
	}
	buf.WriteString(" ORDER BY created_at DESC")

	if f.Limit > 0 {
		args = append(args, f.Limit)
		buf.WriteString(fmt.Sprintf(" LIMIT $%d", len(args)))
	}
	if f.Offset > 0 {
		args = append(args, f.Offset)
		buf.WriteString(fmt.Sprintf(" OFFSET $%d", len(args)))
	}

	return buf.String(), args
}
