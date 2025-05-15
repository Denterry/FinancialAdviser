package service

import (
	"context"
	"fmt"

	"github.com/Denterry/FinancialAdviser/Backend/gateway-service/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Services struct {
	Auth         *AuthService
	Subscription *SubscriptionService
	ML           *MLService
}

type AuthService struct {
	client AuthClient
}

type SubscriptionService struct {
	client SubscriptionClient
}

type MLService struct {
	client MLClient
}

func NewServices(cfg *config.Config) (*Services, error) {
	// Initialize auth service
	authConn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", cfg.Services.Auth.Host, cfg.Services.Auth.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	// Initialize subscription service
	subConn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", cfg.Services.Subscription.Host, cfg.Services.Subscription.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to subscription service: %w", err)
	}

	// Initialize ML service
	mlConn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", cfg.Services.ML.Host, cfg.Services.ML.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ML service: %w", err)
	}

	return &Services{
		Auth:         NewAuthService(authConn),
		Subscription: NewSubscriptionService(subConn),
		ML:           NewMLService(mlConn),
	}, nil
}

// Auth service methods
func (s *AuthService) ValidateToken(token string) (*Claims, error) {
	// TODO: Implement token validation
	return nil, nil
}

// Subscription service methods
func (s *SubscriptionService) GetPlans(ctx context.Context) ([]*Plan, error) {
	// TODO: Implement get plans
	return nil, nil
}

func (s *SubscriptionService) Subscribe(ctx context.Context, planID string) error {
	// TODO: Implement subscription
	return nil
}

func (s *SubscriptionService) GetStatus(ctx context.Context) (*SubscriptionStatus, error) {
	// TODO: Implement get status
	return nil, nil
}

func (s *SubscriptionService) Cancel(ctx context.Context) error {
	// TODO: Implement cancellation
	return nil
}

// ML service methods
func (s *MLService) Analyze(ctx context.Context, data []byte) (*Analysis, error) {
	// TODO: Implement analysis
	return nil, nil
}

func (s *MLService) GetRecommendations(ctx context.Context) ([]*Recommendation, error) {
	// TODO: Implement recommendations
	return nil, nil
}

// Types
type Claims struct {
	UserID string
}

type Plan struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Features    []string
}

type SubscriptionStatus struct {
	Active    bool
	PlanID    string
	StartDate string
	EndDate   string
	AutoRenew bool
}

type Analysis struct {
	Score    float64
	Insights []string
	Risk     string
}

type Recommendation struct {
	Type        string
	Description string
	Priority    int
}
