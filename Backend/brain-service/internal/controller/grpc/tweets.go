package grpc

import (
	"context"

	"github.com/Denterry/FinancialAdviser/Backend/x-service/internal/usecase"
	tweetspb "github.com/Denterry/FinancialAdviser/Backend/x-service/pkg/pb/tweets/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TweetService implements the tweets.v1.TweetService gRPC service
type TweetService struct {
	tweetspb.UnimplementedTweetServiceServer
	tweetUseCase usecase.TweetUseCase
}

// NewTweetService creates a new TweetService
func NewTweetService(tweetUseCase usecase.TweetUseCase) *TweetService {
	return &TweetService{tweetUseCase: tweetUseCase}
}

// Ingest pulls fresh tweets matching the query, persists them, and returns how many were ingested
func (s *TweetService) Ingest(ctx context.Context, req *tweetspb.IngestRequest) (*tweetspb.IngestResponse, error) {
	if req.GetQuery() == "" {
		return nil, status.Error(codes.InvalidArgument, "query is required")
	}
	max := int(req.GetMax())
	if max <= 0 {
		return nil, status.Error(codes.InvalidArgument, "max must be > 0")
	}

	tweets, err := s.tweetUseCase.Ingest(ctx, req.GetQuery(), max)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.tweetUseCase.Ingest(): %v", err)
	}

	return &tweetspb.IngestResponse{
		Ingested: int32(len(tweets)),
	}, nil
}

// ListLatestTweets returns the newest stored tweets, ordered by fetched_at desc
func (s *TweetService) ListLatestTweets(ctx context.Context, req *tweetspb.ListLatestTweetsRequest) (*tweetspb.ListLatestTweetsResponse, error) {
	limit := req.GetLimit()
	if limit <= 0 {
		return nil, status.Error(codes.InvalidArgument, "limit must be > 0")
	}

	tweets, err := s.tweetUseCase.GetListLatest(ctx, limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.tweetUseCase.GetListLatest(): %v", err)
	}

	resp := &tweetspb.ListLatestTweetsResponse{
		Tweets: make([]*tweetspb.Tweet, len(tweets)),
	}

	for i, t := range tweets {
		resp.Tweets[i] = toProtoTweet(t)
	}

	return resp, nil
}

// GetTweetByID returns one stored post by its UUID.
func (s *TweetService) GetTweetByID(ctx context.Context, req *tweetspb.GetTweetByIDRequest) (*tweetspb.GetTweetByIDResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}

	tweet, err := s.tweetUseCase.GetByID(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.tweetUseCase.GetByID: %v", err)
	}
	if tweet == nil {
		return nil, status.Error(codes.NotFound, "post not found")
	}

	return &tweetspb.GetTweetByIDResponse{
		Tweet: toProtoTweet(tweet),
	}, nil
}
