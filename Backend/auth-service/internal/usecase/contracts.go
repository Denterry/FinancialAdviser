// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/Denterry/FinancialAdviser/Backend/auth-service/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=./mocks_usecase_test.go -package=usecase_test

// AuthUseCase defines the interface for authentication operations
type AuthUseCase interface {
	// SignUp creates a new user with the provided credentials
	// Returns a JWT token and error if any
	SignUp(ctx context.Context, email, password, username string) (string, error)

	// SignIn authenticates a user with the provided credentials
	// Returns a JWT token and error if any
	SignIn(ctx context.Context, email, password string) (string, error)

	// ValidateToken validates a JWT token and returns the associated user
	// Returns nil user and error if token is invalid or user not found
	ValidateToken(ctx context.Context, token string) (*entity.User, error)
}
