package service

import (
	"context"
	"time"

	"github.com/Denterry/FinancialAdviser/Backend/gateway-service/internal/entity"
	"google.golang.org/grpc"
)

type SubscriptionService struct {
	// TODO: once proto is defined, add pb client
	conn *grpc.ClientConn
}

func NewSubscriptionService(conn *grpc.ClientConn) *SubscriptionService {
	return &SubscriptionService{conn: conn}
}

func (s *SubscriptionService) GetPlans(ctx context.Context) ([]*entity.Plan, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// TODO: call real gRPC.  Placeholder:
	return []*entity.Plan{
		{ID: "free", Name: "Free", Description: "Basic", Price: 0},
		{ID: "pro", Name: "Pro", Description: "Everything", Price: 19.99},
	}, nil
}

func (s *SubscriptionService) Subscribe(ctx context.Context, planID string) error {
	_, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// TODO
	return nil
}

func (s *SubscriptionService) Status(ctx context.Context) (*entity.Status, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// TODO
	return &entity.Status{Active: false}, nil
}

func (s *SubscriptionService) Cancel(ctx context.Context) error {
	_, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// TODO
	return nil
}
