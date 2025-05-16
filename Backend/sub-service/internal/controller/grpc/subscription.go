package grpc

import (
	"context"
	"fmt"

	"github.com/Denterry/FinancialAdviser/Backend/sub-service/internal/usecase"
	subscriptionv1 "github.com/Denterry/FinancialAdviser/Backend/sub-service/pkg/pb/subscription/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SubscriptionService implements the gRPC subscription service
type SubscriptionService struct {
	subscriptionv1.UnimplementedSubscriptionServiceServer
	subscriptionUseCase *usecase.SubscriptionUseCase
}

// NewSubscriptionService creates a new instance of SubscriptionService
func NewSubscriptionService(subscriptionUseCase *usecase.SubscriptionUseCase) *SubscriptionService {
	return &SubscriptionService{
		subscriptionUseCase: subscriptionUseCase,
	}
}

// CreateSubscription implements the CreateSubscription RPC method
func (s *SubscriptionService) CreateSubscription(ctx context.Context, req *subscriptionv1.CreateSubscriptionRequest) (*subscriptionv1.CreateSubscriptionResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	planID, err := uuid.Parse(req.PlanId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid plan_id format")
	}

	subscription, err := s.subscriptionUseCase.CreateSubscription(ctx, userID, planID, req.AutoRenew, req.PaymentMethod)
	if err != nil {
		switch err {
		case usecase.ErrPlanNotFound:
			return nil, status.Error(codes.NotFound, "plan not found")
		default:
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create subscription: %v", err))
		}
	}

	return &subscriptionv1.CreateSubscriptionResponse{
		Subscription: &subscriptionv1.Subscription{
			Id:              subscription.ID.String(),
			UserId:          subscription.UserID.String(),
			PlanId:          subscription.PlanID.String(),
			Status:          subscriptionv1.SubscriptionStatus(subscription.Status),
			StartDate:       subscription.StartDate.Unix(),
			EndDate:         subscription.EndDate.Unix(),
			AutoRenew:       subscription.AutoRenew,
			AmountPaid:      subscription.AmountPaid,
			Currency:        subscription.Currency,
			PaymentMethod:   subscription.PaymentMethod,
			LastPaymentDate: subscription.LastPaymentDate.Unix(),
			NextPaymentDate: subscription.NextPaymentDate.Unix(),
		},
		PaymentUrl: fmt.Sprintf("/payment/%s", subscription.ID.String()),
	}, nil
}

// GetSubscription implements the GetSubscription RPC method
func (s *SubscriptionService) GetSubscription(ctx context.Context, req *subscriptionv1.GetSubscriptionRequest) (*subscriptionv1.GetSubscriptionResponse, error) {
	id, err := uuid.Parse(req.SubscriptionId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid subscription_id format")
	}

	subscription, err := s.subscriptionUseCase.GetSubscription(ctx, id)
	if err != nil {
		switch err {
		case usecase.ErrSubscriptionNotFound:
			return nil, status.Error(codes.NotFound, "subscription not found")
		default:
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get subscription: %v", err))
		}
	}

	return &subscriptionv1.GetSubscriptionResponse{
		Subscription: &subscriptionv1.Subscription{
			Id:              subscription.ID.String(),
			UserId:          subscription.UserID.String(),
			PlanId:          subscription.PlanID.String(),
			Status:          subscriptionv1.SubscriptionStatus(subscription.Status),
			StartDate:       subscription.StartDate.Unix(),
			EndDate:         subscription.EndDate.Unix(),
			AutoRenew:       subscription.AutoRenew,
			AmountPaid:      subscription.AmountPaid,
			Currency:        subscription.Currency,
			PaymentMethod:   subscription.PaymentMethod,
			LastPaymentDate: subscription.LastPaymentDate.Unix(),
			NextPaymentDate: subscription.NextPaymentDate.Unix(),
		},
	}, nil
}

// ListSubscriptions implements the ListSubscriptions RPC method
func (s *SubscriptionService) ListSubscriptions(ctx context.Context, req *subscriptionv1.ListSubscriptionsRequest) (*subscriptionv1.ListSubscriptionsResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	subscriptions, total, err := s.subscriptionUseCase.ListSubscriptions(ctx, userID, usecase.SubscriptionStatus(req.Status), int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to list subscriptions: %v", err))
	}

	response := &subscriptionv1.ListSubscriptionsResponse{
		Subscriptions: make([]*subscriptionv1.Subscription, len(subscriptions)),
		Total:         int32(total),
	}

	for i, subscription := range subscriptions {
		response.Subscriptions[i] = &subscriptionv1.Subscription{
			Id:              subscription.ID.String(),
			UserId:          subscription.UserID.String(),
			PlanId:          subscription.PlanID.String(),
			Status:          subscriptionv1.SubscriptionStatus(subscription.Status),
			StartDate:       subscription.StartDate.Unix(),
			EndDate:         subscription.EndDate.Unix(),
			AutoRenew:       subscription.AutoRenew,
			AmountPaid:      subscription.AmountPaid,
			Currency:        subscription.Currency,
			PaymentMethod:   subscription.PaymentMethod,
			LastPaymentDate: subscription.LastPaymentDate.Unix(),
			NextPaymentDate: subscription.NextPaymentDate.Unix(),
		}
	}

	return response, nil
}

// UpdateSubscription implements the UpdateSubscription RPC method
func (s *SubscriptionService) UpdateSubscription(ctx context.Context, req *subscriptionv1.UpdateSubscriptionRequest) (*subscriptionv1.UpdateSubscriptionResponse, error) {
	id, err := uuid.Parse(req.SubscriptionId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid subscription_id format")
	}

	subscription, err := s.subscriptionUseCase.GetSubscription(ctx, id)
	if err != nil {
		switch err {
		case usecase.ErrSubscriptionNotFound:
			return nil, status.Error(codes.NotFound, "subscription not found")
		default:
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get subscription: %v", err))
		}
	}

	subscription.AutoRenew = req.AutoRenew
	if req.PaymentMethod != "" {
		subscription.PaymentMethod = req.PaymentMethod
	}

	if err := s.subscriptionUseCase.UpdateSubscription(ctx, subscription); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update subscription: %v", err))
	}

	return &subscriptionv1.UpdateSubscriptionResponse{
		Subscription: &subscriptionv1.Subscription{
			Id:              subscription.ID.String(),
			UserId:          subscription.UserID.String(),
			PlanId:          subscription.PlanID.String(),
			Status:          subscriptionv1.SubscriptionStatus(subscription.Status),
			StartDate:       subscription.StartDate.Unix(),
			EndDate:         subscription.EndDate.Unix(),
			AutoRenew:       subscription.AutoRenew,
			AmountPaid:      subscription.AmountPaid,
			Currency:        subscription.Currency,
			PaymentMethod:   subscription.PaymentMethod,
			LastPaymentDate: subscription.LastPaymentDate.Unix(),
			NextPaymentDate: subscription.NextPaymentDate.Unix(),
		},
	}, nil
}

// CancelSubscription implements the CancelSubscription RPC method
func (s *SubscriptionService) CancelSubscription(ctx context.Context, req *subscriptionv1.CancelSubscriptionRequest) (*subscriptionv1.CancelSubscriptionResponse, error) {
	id, err := uuid.Parse(req.SubscriptionId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid subscription_id format")
	}

	if err := s.subscriptionUseCase.CancelSubscription(ctx, id); err != nil {
		switch err {
		case usecase.ErrSubscriptionNotFound:
			return nil, status.Error(codes.NotFound, "subscription not found")
		default:
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to cancel subscription: %v", err))
		}
	}

	subscription, err := s.subscriptionUseCase.GetSubscription(ctx, id)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get subscription: %v", err))
	}

	return &subscriptionv1.CancelSubscriptionResponse{
		Subscription: &subscriptionv1.Subscription{
			Id:              subscription.ID.String(),
			UserId:          subscription.UserID.String(),
			PlanId:          subscription.PlanID.String(),
			Status:          subscriptionv1.SubscriptionStatus(subscription.Status),
			StartDate:       subscription.StartDate.Unix(),
			EndDate:         subscription.EndDate.Unix(),
			AutoRenew:       subscription.AutoRenew,
			AmountPaid:      subscription.AmountPaid,
			Currency:        subscription.Currency,
			PaymentMethod:   subscription.PaymentMethod,
			LastPaymentDate: subscription.LastPaymentDate.Unix(),
			NextPaymentDate: subscription.NextPaymentDate.Unix(),
		},
	}, nil
}

// GetPlans implements the GetPlans RPC method
func (s *SubscriptionService) GetPlans(ctx context.Context, req *subscriptionv1.GetPlansRequest) (*subscriptionv1.GetPlansResponse, error) {
	plans, err := s.subscriptionUseCase.GetPlans(ctx, usecase.PlanType(req.Type))
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get plans: %v", err))
	}

	response := &subscriptionv1.GetPlansResponse{
		Plans: make([]*subscriptionv1.Plan, len(plans)),
	}

	for i, plan := range plans {
		response.Plans[i] = &subscriptionv1.Plan{
			Id:           plan.ID.String(),
			Name:         plan.Name,
			Description:  plan.Description,
			Type:         subscriptionv1.PlanType(plan.Type),
			Price:        plan.Price,
			Currency:     plan.Currency,
			DurationDays: int32(plan.DurationDays),
			Features:     plan.Features,
		}
	}

	return response, nil
}

// ProcessPayment implements the ProcessPayment RPC method
func (s *SubscriptionService) ProcessPayment(ctx context.Context, req *subscriptionv1.ProcessPaymentRequest) (*subscriptionv1.ProcessPaymentResponse, error) {
	subscriptionID, err := uuid.Parse(req.SubscriptionId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid subscription_id format")
	}

	payment, err := s.subscriptionUseCase.ProcessPayment(ctx, subscriptionID, req.Amount, req.Currency, req.PaymentMethod)
	if err != nil {
		switch err {
		case usecase.ErrSubscriptionNotFound:
			return nil, status.Error(codes.NotFound, "subscription not found")
		case usecase.ErrPaymentFailed:
			return nil, status.Error(codes.Internal, "payment processing failed")
		default:
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to process payment: %v", err))
		}
	}

	return &subscriptionv1.ProcessPaymentResponse{
		PaymentId:     payment.ID.String(),
		Status:        subscriptionv1.PaymentStatus(payment.Status),
		TransactionId: payment.TransactionID,
		Timestamp:     payment.CreatedAt.Unix(),
	}, nil
}
