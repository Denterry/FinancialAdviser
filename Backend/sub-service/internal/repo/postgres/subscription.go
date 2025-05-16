package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/Denterry/FinancialAdviser/Backend/sub-service/internal/entity"
	"github.com/Denterry/FinancialAdviser/Backend/sub-service/internal/repo"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type subscriptionRepo struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) repo.SubscriptionRepository {
	return &subscriptionRepo{db: db}
}

// Plan operations
func (r *subscriptionRepo) CreatePlan(ctx context.Context, plan *entity.Plan) error {
	query := `
		INSERT INTO plans (id, name, description, type, price, currency, duration_days, features)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		plan.ID,
		plan.Name,
		plan.Description,
		plan.Type,
		plan.Price,
		plan.Currency,
		plan.DurationDays,
		pq.Array(plan.Features),
	)
	return err
}

func (r *subscriptionRepo) GetPlan(ctx context.Context, id uuid.UUID) (*entity.Plan, error) {
	query := `
		SELECT id, name, description, type, price, currency, duration_days, features, created_at, updated_at
		FROM plans
		WHERE id = $1
	`
	plan := &entity.Plan{}
	var features []string
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&plan.ID,
		&plan.Name,
		&plan.Description,
		&plan.Type,
		&plan.Price,
		&plan.Currency,
		&plan.DurationDays,
		pq.Array(&features),
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repo.ErrPlanNotFound
		}
		return nil, err
	}
	plan.Features = features
	return plan, nil
}

func (r *subscriptionRepo) ListPlans(ctx context.Context) ([]*entity.Plan, error) {
	query := `
		SELECT id, name, description, type, price, currency, duration_days, features, created_at, updated_at
		FROM plans
		ORDER BY price ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []*entity.Plan
	for rows.Next() {
		plan := &entity.Plan{}
		var features []string
		err := rows.Scan(
			&plan.ID,
			&plan.Name,
			&plan.Description,
			&plan.Type,
			&plan.Price,
			&plan.Currency,
			&plan.DurationDays,
			pq.Array(&features),
			&plan.CreatedAt,
			&plan.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		plan.Features = features
		plans = append(plans, plan)
	}
	return plans, nil
}

// Subscription operations
func (r *subscriptionRepo) CreateSubscription(ctx context.Context, sub *entity.Subscription) error {
	query := `
		INSERT INTO subscriptions (
			id, user_id, plan_id, status, start_date, end_date,
			auto_renew, amount_paid, currency, payment_method,
			last_payment_date, next_payment_date
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.ExecContext(ctx, query,
		sub.ID,
		sub.UserID,
		sub.PlanID,
		sub.Status,
		sub.StartDate,
		sub.EndDate,
		sub.AutoRenew,
		sub.AmountPaid,
		sub.Currency,
		sub.PaymentMethod,
		sub.LastPaymentDate,
		sub.NextPaymentDate,
	)
	return err
}

func (r *subscriptionRepo) GetSubscription(ctx context.Context, id uuid.UUID) (*entity.Subscription, error) {
	query := `
		SELECT id, user_id, plan_id, status, start_date, end_date,
			auto_renew, amount_paid, currency, payment_method,
			last_payment_date, next_payment_date, created_at, updated_at
		FROM subscriptions
		WHERE id = $1
	`
	sub := &entity.Subscription{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&sub.ID,
		&sub.UserID,
		&sub.PlanID,
		&sub.Status,
		&sub.StartDate,
		&sub.EndDate,
		&sub.AutoRenew,
		&sub.AmountPaid,
		&sub.Currency,
		&sub.PaymentMethod,
		&sub.LastPaymentDate,
		&sub.NextPaymentDate,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repo.ErrSubscriptionNotFound
		}
		return nil, err
	}
	return sub, nil
}

func (r *subscriptionRepo) ListSubscriptions(ctx context.Context, filter repo.SubscriptionFilter) ([]*entity.Subscription, error) {
	query := `
		SELECT s.id, s.user_id, s.plan_id, s.status, s.start_date, s.end_date,
			s.auto_renew, s.amount_paid, s.currency, s.payment_method,
			s.last_payment_date, s.next_payment_date, s.created_at, s.updated_at
		FROM subscriptions s
		WHERE ($1::uuid IS NULL OR s.user_id = $1)
		AND ($2::subscription_status IS NULL OR s.status = $2)
		AND ($3::plan_type IS NULL OR EXISTS (
			SELECT 1 FROM plans p WHERE p.id = s.plan_id AND p.type = $3
		))
		ORDER BY s.created_at DESC
		LIMIT $4 OFFSET $5
	`
	rows, err := r.db.QueryContext(ctx, query,
		filter.UserID,
		filter.Status,
		filter.PlanType,
		filter.Limit,
		filter.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*entity.Subscription
	for rows.Next() {
		sub := &entity.Subscription{}
		err := rows.Scan(
			&sub.ID,
			&sub.UserID,
			&sub.PlanID,
			&sub.Status,
			&sub.StartDate,
			&sub.EndDate,
			&sub.AutoRenew,
			&sub.AmountPaid,
			&sub.Currency,
			&sub.PaymentMethod,
			&sub.LastPaymentDate,
			&sub.NextPaymentDate,
			&sub.CreatedAt,
			&sub.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, sub)
	}
	return subscriptions, nil
}

func (r *subscriptionRepo) UpdateSubscription(ctx context.Context, sub *entity.Subscription) error {
	query := `
		UPDATE subscriptions
		SET status = $1,
			auto_renew = $2,
			payment_method = $3,
			last_payment_date = $4,
			next_payment_date = $5,
			updated_at = $6
		WHERE id = $7
	`
	_, err := r.db.ExecContext(ctx, query,
		sub.Status,
		sub.AutoRenew,
		sub.PaymentMethod,
		sub.LastPaymentDate,
		sub.NextPaymentDate,
		time.Now(),
		sub.ID,
	)
	return err
}

func (r *subscriptionRepo) CancelSubscription(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE subscriptions
		SET status = 'cancelled',
			auto_renew = false,
			updated_at = $1
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

// Payment operations
func (r *subscriptionRepo) CreatePayment(ctx context.Context, payment *entity.Payment) error {
	query := `
		INSERT INTO payments (
			id, subscription_id, amount, currency,
			status, payment_method, transaction_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query,
		payment.ID,
		payment.SubscriptionID,
		payment.Amount,
		payment.Currency,
		payment.Status,
		payment.PaymentMethod,
		payment.TransactionID,
	)
	return err
}

func (r *subscriptionRepo) GetPayment(ctx context.Context, id uuid.UUID) (*entity.Payment, error) {
	query := `
		SELECT id, subscription_id, amount, currency,
			status, payment_method, transaction_id,
			created_at, updated_at
		FROM payments
		WHERE id = $1
	`
	payment := &entity.Payment{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&payment.ID,
		&payment.SubscriptionID,
		&payment.Amount,
		&payment.Currency,
		&payment.Status,
		&payment.PaymentMethod,
		&payment.TransactionID,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repo.ErrPaymentNotFound
		}
		return nil, err
	}
	return payment, nil
}

func (r *subscriptionRepo) ListPayments(ctx context.Context, subscriptionID uuid.UUID) ([]*entity.Payment, error) {
	query := `
		SELECT id, subscription_id, amount, currency,
			status, payment_method, transaction_id,
			created_at, updated_at
		FROM payments
		WHERE subscription_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, subscriptionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*entity.Payment
	for rows.Next() {
		payment := &entity.Payment{}
		err := rows.Scan(
			&payment.ID,
			&payment.SubscriptionID,
			&payment.Amount,
			&payment.Currency,
			&payment.Status,
			&payment.PaymentMethod,
			&payment.TransactionID,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	return payments, nil
}

func (r *subscriptionRepo) UpdatePayment(ctx context.Context, payment *entity.Payment) error {
	query := `
		UPDATE payments
		SET status = $1,
			transaction_id = $2,
			updated_at = $3
		WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, query,
		payment.Status,
		payment.TransactionID,
		time.Now(),
		payment.ID,
	)
	return err
}
