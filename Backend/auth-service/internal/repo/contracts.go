// Package repo implements application outer layer logic. Each logic group in own file.
package repo

import (
	"context"

	"github.com/Denterry/FinancialAdviser/Backend/auth-service/internal/entity"
	"github.com/google/uuid"
)

//go:generate mockgen -source=contracts.go -destination=../usecase/mocks_repo_test.go -package=usecase_test

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// Create creates a new user in the database
	// Returns error if user creation fails
	Create(ctx context.Context, user *entity.User) error

	// GetByEmail retrieves a user by their email address
	// Returns nil, nil if user is not found
	// Returns error if database operation fails
	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	// GetByID retrieves a user by their UUID
	// Returns nil, nil if user is not found
	// Returns error if database operation fails
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)

	// List retrieves all users from the database
	// Returns empty slice if no users found
	// Returns error if database operation fails
	List(ctx context.Context) ([]*entity.User, error)

	// Update updates an existing user's information
	// Returns error if user update fails
	Update(ctx context.Context, user *entity.User) error

	// Delete removes a user from the database by their UUID
	// Returns error if user deletion fails
	Delete(ctx context.Context, id uuid.UUID) error
}
