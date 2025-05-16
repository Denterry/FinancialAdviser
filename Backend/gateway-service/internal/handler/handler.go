package handler

import (
	"github.com/Denterry/FinancialAdviser/Backend/gateway-service/internal/service"
)

type Handlers struct {
	Auth         *AuthHandler
	Subscription *SubscriptionHandler
	ML           *MLHandler
}

func NewHandlers(services *service.Services) *Handlers {
	return &Handlers{
		Auth:         NewAuthHandler(services.Auth),
		Subscription: NewSubscriptionHandler(services.Subscription),
		ML:           NewMLHandler(services.ML),
	}
}
