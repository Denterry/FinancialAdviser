//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ../../controller/grpc/proto/auth/v1/auth_service.proto

package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Denterry/FinancialAdviser/Backend/auth-service/internal/entity"
	"github.com/Denterry/FinancialAdviser/Backend/auth-service/internal/repo"
)

// UseCase implements the authentication use case
type UseCase struct {
	userRepo repo.UserRepository
	jwtKey   []byte
	tokenTTL time.Duration
}

// NewUseCase creates a new instance of authentication UseCase
func NewUseCase(userRepo repo.UserRepository, jwtKey string, tokenTTL int) *UseCase {
	return &UseCase{
		userRepo: userRepo,
		jwtKey:   []byte(jwtKey),
		tokenTTL: time.Duration(tokenTTL) * time.Minute,
	}
}

// SignUp registers a new user
func (uc *UseCase) SignUp(ctx context.Context, email, password, username string) (string, error) {
	existingUser, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("uc.userRepo.GetByEmail(): %w", err)
	}
	if existingUser != nil {
		return "", ErrUserAlreadyExists
	}

	// create & hash
	user := &entity.User{
		ID:       uuid.New(),
		Email:    email,
		Username: username,
		IsAdmin:  false,
	}
	if err := user.SetPassword(password); err != nil {
		return "", fmt.Errorf("user.SetPassword(): %w", err)
	}
	now := time.Now().UTC()
	user.CreatedAt = now

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return "", fmt.Errorf("uc.userRepo.Create(): %w", err)
	}

	return uc.makeToken(user, now)
}

// SignIn authenticates a user and returns a JWT token
func (uc *UseCase) SignIn(ctx context.Context, email, password string) (string, error) {
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("uc.userRepo.GetByEmail(): %w", err)
	}
	if user == nil {
		return "", ErrUserNotFound
	}

	// password check
	if user == nil || !user.CheckPassword(password) { // user == nil exactly means that email is not right
		return "", ErrInvalidCredentials
	}

	now := time.Now().UTC()
	return uc.makeToken(user, now)
}
