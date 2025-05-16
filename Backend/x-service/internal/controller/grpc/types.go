package grpc

import (
	"time"

	"github.com/Denterry/FinancialAdviser/Backend/x-service/internal/entity"
	adminpb "github.com/Denterry/FinancialAdviser/Backend/x-service/pkg/pb/admin/v1"
	tweetspb "github.com/Denterry/FinancialAdviser/Backend/x-service/pkg/pb/tweets/v1"
	"github.com/google/uuid"
)

func toProtoAdminTweet(t *entity.Tweet) *adminpb.Tweet {
	if t == nil {
		return nil
	}

	return &adminpb.Tweet{
		Id:        t.ID.String(),
		Text:      t.Text,
		AuthorId:  t.AuthorID,
		CreatedAt: t.CreatedAt.Unix(),
		UpdatedAt: t.UpdatedAt.Unix(),
		Sentiment: &adminpb.Sentiment{
			Score: t.SentimentScore,
			Label: t.SentimentLabel,
		},
		IsFinancial: t.IsFinancial,
		Symbols:     t.Symbols,
		Engagement: &adminpb.Engagement{
			RetweetCount:  int32(t.Retweets),
			FavoriteCount: int32(t.Likes),
			ReplyCount:    int32(t.Replies),
		},
	}
}

func toEntityAdminTweet(t *adminpb.Tweet) *entity.Tweet {
	if t == nil {
		return nil
	}

	return &entity.Tweet{
		ID:             uuid.MustParse(t.Id),
		Text:           t.Text,
		AuthorID:       t.AuthorId,
		CreatedAt:      time.Unix(t.CreatedAt, 0),
		UpdatedAt:      time.Unix(t.UpdatedAt, 0),
		SentimentScore: t.Sentiment.Score,
		SentimentLabel: t.Sentiment.Label,
		IsFinancial:    t.IsFinancial,
		Symbols:        t.Symbols,
		Retweets:       int(t.Engagement.RetweetCount),
		Likes:          int(t.Engagement.FavoriteCount),
		Replies:        int(t.Engagement.ReplyCount),
	}
}

func toProtoTweet(t *entity.Tweet) *tweetspb.Tweet {
	if t == nil {
		return nil
	}

	return &tweetspb.Tweet{
		Id:        t.ID.String(),
		AuthorId:  t.AuthorID,
		Username:  t.UserName,
		Text:      t.Text,
		Lang:      t.Lang,
		CreatedAt: t.CreatedAt.Unix(),
		FetchedAt: t.FetchedAt.Unix(),
		Likes:     int32(t.Likes),
		Replies:   int32(t.Replies),
		Retweets:  int32(t.Retweets),
		Views:     int32(t.Views),
		Urls:      t.URLs,
		Photos:    t.Photos,
		Videos:    t.Videos,
	}
}

func toEntityTweet(t *tweetspb.Tweet) *entity.Tweet {
	if t == nil {
		return nil
	}

	return &entity.Tweet{
		ID:        uuid.MustParse(t.Id),
		AuthorID:  t.AuthorId,
		UserName:  t.Username,
		Text:      t.Text,
		Lang:      t.Lang,
		CreatedAt: time.Unix(t.CreatedAt, 0),
		FetchedAt: time.Unix(t.FetchedAt, 0),
		Likes:     int(t.Likes),
		Replies:   int(t.Replies),
		Retweets:  int(t.Retweets),
		Views:     int(t.Views),
		URLs:      t.Urls,
		Photos:    t.Photos,
		Videos:    t.Videos,
	}
}
