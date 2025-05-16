package repo

import (
	"context"

	"github.com/Denterry/FinancialAdviser/Backend/sub-service/internal/entity"
	"github.com/google/uuid"
)

// SubscriptionRepository defines the interface for subscription data access
type SubscriptionRepository interface {
	// Plan operations
	CreatePlan(ctx context.Context, plan *entity.Plan) error
	GetPlan(ctx context.Context, id uuid.UUID) (*entity.Plan, error)
	ListPlans(ctx context.Context, planType entity.PlanType) ([]*entity.Plan, error)
	UpdatePlan(ctx context.Context, plan *entity.Plan) error
	DeletePlan(ctx context.Context, id uuid.UUID) error

	// Subscription operations
	CreateSubscription(ctx context.Context, subscription *entity.Subscription) error
	GetSubscription(ctx context.Context, id uuid.UUID) (*entity.Subscription, error)
	ListSubscriptions(ctx context.Context, userID uuid.UUID, status entity.SubscriptionStatus, limit, offset int) ([]*entity.Subscription, int, error)
	UpdateSubscription(ctx context.Context, subscription *entity.Subscription) error
	CancelSubscription(ctx context.Context, id uuid.UUID) error

	// Payment operations
	CreatePayment(ctx context.Context, payment *entity.Payment) error
	GetPayment(ctx context.Context, id uuid.UUID) (*entity.Payment, error)
	ListPayments(ctx context.Context, subscriptionID uuid.UUID) ([]*entity.Payment, error)
	UpdatePayment(ctx context.Context, payment *entity.Payment) error
}

// SubscriptionFilter represents the filter criteria for listing subscriptions
type SubscriptionFilter struct {
	UserID   uuid.UUID
	Status   entity.SubscriptionStatus
	PlanType entity.PlanType
	Limit    int
	Offset   int
}
