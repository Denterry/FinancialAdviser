package service

import (
	"context"
	"time"

	"github.com/Denterry/FinancialAdviser/Backend/gateway-service/internal/entity"
	"google.golang.org/grpc"
)

type MLService struct{ conn *grpc.ClientConn }

func NewMLService(conn *grpc.ClientConn) *MLService { return &MLService{conn: conn} }

func (m *MLService) Analyze(ctx context.Context, payload []byte) (*entity.Analysis, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// TODO: call brain / ml proto
	return &entity.Analysis{
		Score:    0.8,
		Insights: []string{"stub insight"},
		Risk:     "medium",
	}, nil
}

func (m *MLService) Recommend(ctx context.Context) ([]*entity.Recommendation, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// TODO
	return []*entity.Recommendation{
		{Type: "BUY", Description: "Consider TSLA", Priority: 1},
	}, nil
}
