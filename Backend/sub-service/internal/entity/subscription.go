package entity

import (
	"time"

	"github.com/google/uuid"
)

// PlanType represents the type of subscription plan
type PlanType string

const (
	PlanTypeBasic      PlanType = "basic"
	PlanTypePro        PlanType = "pro"
	PlanTypeEnterprise PlanType = "enterprise"
)

// SubscriptionStatus represents the status of a subscription
type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusCancelled SubscriptionStatus = "cancelled"
	SubscriptionStatusExpired   SubscriptionStatus = "expired"
	SubscriptionStatusPending   SubscriptionStatus = "pending"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

// Plan represents a subscription plan
type Plan struct {
	ID           uuid.UUID
	Name         string
	Description  string
	Type         PlanType
	Price        float64
	Currency     string
	DurationDays int
	Features     []string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Subscription represents a user's subscription
type Subscription struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	PlanID          uuid.UUID
	Status          SubscriptionStatus
	StartDate       time.Time
	EndDate         time.Time
	AutoRenew       bool
	AmountPaid      float64
	Currency        string
	PaymentMethod   string
	LastPaymentDate time.Time
	NextPaymentDate time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Payment represents a subscription payment
type Payment struct {
	ID             uuid.UUID
	SubscriptionID uuid.UUID
	Amount         float64
	Currency       string
	Status         PaymentStatus
	PaymentMethod  string
	TransactionID  string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
