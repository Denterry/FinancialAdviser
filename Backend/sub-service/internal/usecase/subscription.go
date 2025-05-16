package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/Denterry/FinancialAdviser/Backend/sub-service/internal/entity"
	"github.com/Denterry/FinancialAdviser/Backend/sub-service/internal/repo"
	"github.com/google/uuid"
)

var (
	ErrPlanNotFound         = errors.New("plan not found")
	ErrSubscriptionNotFound = errors.New("subscription not found")
	ErrInvalidPlan          = errors.New("invalid plan")
	ErrInvalidPayment       = errors.New("invalid payment")
	ErrPaymentFailed        = errors.New("payment failed")
)

// SubscriptionUseCase implements subscription business logic
type SubscriptionUseCase struct {
	repo repo.SubscriptionRepository
}

// NewSubscriptionUseCase creates a new instance of SubscriptionUseCase
func NewSubscriptionUseCase(repo repo.SubscriptionRepository) *SubscriptionUseCase {
	return &SubscriptionUseCase{
		repo: repo,
	}
}

// CreateSubscription creates a new subscription for a user
func (uc *SubscriptionUseCase) CreateSubscription(ctx context.Context, userID, planID uuid.UUID, autoRenew bool, paymentMethod string) (*entity.Subscription, error) {
	// Get plan details
	plan, err := uc.repo.GetPlan(ctx, planID)
	if err != nil {
		return nil, ErrPlanNotFound
	}

	// Create subscription
	subscription := &entity.Subscription{
		ID:              uuid.New(),
		UserID:          userID,
		PlanID:          planID,
		Status:          entity.SubscriptionStatusPending,
		StartDate:       time.Now(),
		EndDate:         time.Now().AddDate(0, 0, plan.DurationDays),
		AutoRenew:       autoRenew,
		AmountPaid:      plan.Price,
		Currency:        plan.Currency,
		PaymentMethod:   paymentMethod,
		LastPaymentDate: time.Now(),
		NextPaymentDate: time.Now().AddDate(0, 0, plan.DurationDays),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := uc.repo.CreateSubscription(ctx, subscription); err != nil {
		return nil, err
	}

	return subscription, nil
}

// GetSubscription retrieves a subscription by ID
func (uc *SubscriptionUseCase) GetSubscription(ctx context.Context, id uuid.UUID) (*entity.Subscription, error) {
	subscription, err := uc.repo.GetSubscription(ctx, id)
	if err != nil {
		return nil, ErrSubscriptionNotFound
	}
	return subscription, nil
}

// ListSubscriptions retrieves a list of subscriptions for a user
func (uc *SubscriptionUseCase) ListSubscriptions(ctx context.Context, userID uuid.UUID, status entity.SubscriptionStatus, limit, offset int) ([]*entity.Subscription, int, error) {
	return uc.repo.ListSubscriptions(ctx, userID, status, limit, offset)
}

// UpdateSubscription updates a subscription's details
func (uc *SubscriptionUseCase) UpdateSubscription(ctx context.Context, subscription *entity.Subscription) error {
	if err := uc.repo.UpdateSubscription(ctx, subscription); err != nil {
		return err
	}
	return nil
}

// CancelSubscription cancels a subscription
func (uc *SubscriptionUseCase) CancelSubscription(ctx context.Context, id uuid.UUID) error {
	subscription, err := uc.repo.GetSubscription(ctx, id)
	if err != nil {
		return ErrSubscriptionNotFound
	}

	subscription.Status = entity.SubscriptionStatusCancelled
	subscription.AutoRenew = false
	subscription.UpdatedAt = time.Now()

	return uc.repo.UpdateSubscription(ctx, subscription)
}

// ProcessPayment processes a payment for a subscription
func (uc *SubscriptionUseCase) ProcessPayment(ctx context.Context, subscriptionID uuid.UUID, amount float64, currency, paymentMethod string) (*entity.Payment, error) {
	subscription, err := uc.repo.GetSubscription(ctx, subscriptionID)
	if err != nil {
		return nil, ErrSubscriptionNotFound
	}

	payment := &entity.Payment{
		ID:             uuid.New(),
		SubscriptionID: subscriptionID,
		Amount:         amount,
		Currency:       currency,
		Status:         entity.PaymentStatusPending,
		PaymentMethod:  paymentMethod,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := uc.repo.CreatePayment(ctx, payment); err != nil {
		return nil, err
	}

	// TODO: Integrate with payment gateway
	// For now, we'll just mark it as completed
	payment.Status = entity.PaymentStatusCompleted
	payment.TransactionID = uuid.New().String()
	payment.UpdatedAt = time.Now()

	if err := uc.repo.UpdatePayment(ctx, payment); err != nil {
		return nil, err
	}

	// Update subscription status
	subscription.Status = entity.SubscriptionStatusActive
	subscription.LastPaymentDate = time.Now()
	subscription.NextPaymentDate = time.Now().AddDate(0, 0, 30) // Assuming monthly subscription
	subscription.UpdatedAt = time.Now()

	if err := uc.repo.UpdateSubscription(ctx, subscription); err != nil {
		return nil, err
	}

	return payment, nil
}

// GetPlans retrieves available subscription plans
func (uc *SubscriptionUseCase) GetPlans(ctx context.Context, planType entity.PlanType) ([]*entity.Plan, error) {
	return uc.repo.ListPlans(ctx, planType)
}
