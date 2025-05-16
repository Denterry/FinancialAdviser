package service

import (
	"fmt"

	"github.com/Denterry/FinancialAdviser/Backend/gateway-service/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Services is a DI-container exposed to handlers
type Services struct {
	Auth         *AuthService
	Subscription *SubscriptionService
	ML           *MLService
}

// NewServices -.
func NewServices(cfg *config.Config) (*Services, error) {
	dial := func(s config.GRPCService) (*grpc.ClientConn, error) {
		addr := s.Addr()
		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}
		return grpc.NewClient(addr, opts...)
	}

	authConn, err := dial(cfg.Services.Auth)
	if err != nil {
		return nil, fmt.Errorf("auth dial: %w", err)
	}

	subConn, err := dial(cfg.Services.Sub)
	if err != nil {
		return nil, fmt.Errorf("sub dial: %w", err)
	}

	mlConn, err := dial(cfg.Services.ML)
	if err != nil {
		return nil, fmt.Errorf("ml dial: %w", err)
	}

	return &Services{
		Auth:         NewAuthService(authConn),
		Subscription: NewSubscriptionService(subConn),
		ML:           NewMLService(mlConn),
	}, nil
}
