package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/Denterry/FinancialAdviser/Backend/x-service/internal/repo"
	"github.com/Denterry/FinancialAdviser/Backend/x-service/internal/usecase/admin"
	adminpb "github.com/Denterry/FinancialAdviser/Backend/x-service/pkg/pb/admin/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AdminTweetService is a gRPC service for admin operations on tweets
type AdminTweetService struct {
	adminpb.UnimplementedAdminTweetServiceServer
	adminTweetUseCase *admin.UseCase
}

// NewAdminTweetService creates a new AdminTweetService
func NewAdminTweetService(adminTweetUseCase *admin.UseCase) *AdminTweetService {
	return &AdminTweetService{
		adminTweetUseCase: adminTweetUseCase,
	}
}

// CreateTweet creates a new tweet in the database
func (s *AdminTweetService) CreateTweet(ctx context.Context, req *adminpb.CreateTweetRequest) (*adminpb.CreateTweetResponse, error) {
	if req.GetText() == "" {
		return nil, status.Error(codes.InvalidArgument, "text is required")
	}
	if req.GetAuthorId() == "" {
		return nil, status.Error(codes.InvalidArgument, "author_id is required")
	}

	tweet, err := s.adminTweetUseCase.Create(ctx, req.Text, req.AuthorId)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("s.adminTweetUseCase.Create(): %v", err))
	}

	return &adminpb.CreateTweetResponse{
		Tweet: toProtoAdminTweet(tweet),
	}, nil
}

// GetTweet retrieves a tweet from the database
func (s *AdminTweetService) GetTweet(ctx context.Context, req *adminpb.GetTweetRequest) (*adminpb.GetTweetResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id format")
	}

	tweet, err := s.adminTweetUseCase.Get(ctx, id)
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("s.adminTweetUseCase.Get(): %v", err))
	}

	return &adminpb.GetTweetResponse{
		Tweet: toProtoAdminTweet(tweet),
	}, nil
}

// ListTweets lists tweets from the database
func (s *AdminTweetService) ListTweets(ctx context.Context, req *adminpb.ListTweetsRequest) (*adminpb.ListTweetsResponse, error) {
	filter := repo.TweetFilter{
		AuthorID:       req.AuthorId,
		IsFinancial:    &req.IsFinancial,
		SentimentLabel: req.SentimentLabel,
		Symbols:        req.Symbols,
		Limit:          int32(req.Limit),
		Offset:         int32(req.Offset),
	}

	if req.StartTime != 0 {
		startTime := time.Unix(req.StartTime, 0)
		filter.StartTime = &startTime
	}
	if req.EndTime != 0 {
		endTime := time.Unix(req.EndTime, 0)
		filter.EndTime = &endTime
	}

	tweets, err := s.adminTweetUseCase.List(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("s.adminTweetUseCase.List(): %v", err))
	}

	response := &adminpb.ListTweetsResponse{
		Tweets: make([]*adminpb.Tweet, len(tweets)),
	}

	for i, t := range tweets {
		response.Tweets[i] = toProtoAdminTweet(t)
	}

	return response, nil
}

// UpdateTweet updates a tweet in the database
func (s *AdminTweetService) UpdateTweet(ctx context.Context, req *adminpb.UpdateTweetRequest) (*adminpb.UpdateTweetResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id format")
	}

	tweet, err := s.adminTweetUseCase.Get(ctx, id)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("s.adminTweetUseCase.Get(): %v", err))
	}

	if req.GetText() != "" {
		tweet.Text = req.GetText()
	}
	if req.GetSentiment() != nil {
		tweet.SentimentScore = req.GetSentiment().GetScore()
		tweet.SentimentLabel = req.GetSentiment().GetLabel()
	}
	if req.GetEngagement() != nil {
		tweet.Retweets = int(req.GetEngagement().GetRetweetCount())
		tweet.Likes = int(req.GetEngagement().GetFavoriteCount())
		tweet.Replies = int(req.GetEngagement().GetReplyCount())
	}

	if err := s.adminTweetUseCase.Update(ctx, tweet); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("s.adminTweetUseCase.Update(): %v", err))
	}

	return &adminpb.UpdateTweetResponse{Tweet: toProtoAdminTweet(tweet)}, nil
}

// DeleteTweet deletes a tweet from the database
func (s *AdminTweetService) DeleteTweet(ctx context.Context, req *adminpb.DeleteTweetRequest) (*adminpb.DeleteTweetResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id format")
	}

	if err := s.adminTweetUseCase.Delete(ctx, id); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("s.adminTweetUseCase.Delete(): %v", err))
	}

	return &adminpb.DeleteTweetResponse{}, nil
}

func (s *AdminTweetService) GetTweetsBySymbol(ctx context.Context, req *adminpb.GetTweetsBySymbolRequest) (*adminpb.GetTweetsBySymbolResponse, error) {
	if req.GetSymbol() == "" {
		return nil, status.Error(codes.InvalidArgument, "symbol is required")
	}

	tweets, err := s.adminTweetUseCase.GetTweetsBySymbol(ctx, req.GetSymbol(), req.GetLimit(), req.GetOffset())
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("s.adminTweetUseCase.GetTweetsBySymbol(): %v", err))
	}

	response := &adminpb.GetTweetsBySymbolResponse{
		Tweets: make([]*adminpb.Tweet, len(tweets)),
	}

	for i, t := range tweets {
		response.Tweets[i] = toProtoAdminTweet(t)
	}

	return response, nil
}

// GetTweetsBySentiment retrieves tweets by sentiment from the database
func (s *AdminTweetService) GetTweetsBySentiment(ctx context.Context, req *adminpb.GetTweetsBySentimentRequest) (*adminpb.GetTweetsBySentimentResponse, error) {
	if req.GetLabel() == "" {
		return nil, status.Error(codes.InvalidArgument, "label is required")
	}

	tweets, err := s.adminTweetUseCase.GetTweetsBySentiment(ctx, req.GetLabel(), req.GetLimit(), req.GetOffset())
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("s.adminTweetUseCase.GetTweetsBySentiment(): %v", err))
	}

	response := &adminpb.GetTweetsBySentimentResponse{
		Tweets: make([]*adminpb.Tweet, len(tweets)),
	}

	for i, t := range tweets {
		response.Tweets[i] = toProtoAdminTweet(t)
	}

	return response, nil
}
