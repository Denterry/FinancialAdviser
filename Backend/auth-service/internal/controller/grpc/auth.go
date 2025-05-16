package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Denterry/FinancialAdviser/Backend/auth-service/internal/usecase"
	authv1 "github.com/Denterry/FinancialAdviser/Backend/auth-service/pkg/pb/auth/v1"
)

// AuthService implements the gRPC auth service
type AuthService struct {
	authv1.UnimplementedAuthServiceServer
	authUseCase usecase.AuthUseCase
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(authUseCase usecase.AuthUseCase) *AuthService {
	return &AuthService{
		authUseCase: authUseCase,
	}
}

// SignIn implements the SignIn RPC method
func (s *AuthService) SignIn(ctx context.Context, req *authv1.SignInRequest) (*authv1.SignInResponse, error) {
	if req.GetEmail() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	token, err := s.authUseCase.SignIn(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		switch err {
		case usecase.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, "user not found")
		case usecase.ErrInvalidPassword:
			return nil, status.Error(codes.InvalidArgument, "invalid password")
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	user, err := s.authUseCase.ValidateToken(ctx, token)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Errorf("s.authUseCase.ValidateToken(): %w", err).Error())
	}

	return &authv1.SignInResponse{
		Token:  token,
		UserId: user.ID.String(),
	}, nil
}

// SignUp implements the SignUp RPC method
func (s *AuthService) SignUp(ctx context.Context, req *authv1.SignUpRequest) (*authv1.SignUpResponse, error) {
	if req.GetEmail() == "" || req.GetPassword() == "" || req.GetUsername() == "" {
		return nil, status.Error(codes.InvalidArgument, "email, password, and username are required")
	}

	token, err := s.authUseCase.SignUp(ctx, req.GetEmail(), req.GetPassword(), req.GetUsername())
	if err != nil {
		switch err {
		case usecase.ErrUserAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		case usecase.ErrInvalidEmail:
			return nil, status.Error(codes.InvalidArgument, "invalid email format")
		case usecase.ErrInvalidPassword:
			return nil, status.Error(codes.InvalidArgument, "invalid password format")
		case usecase.ErrInvalidUsername:
			return nil, status.Error(codes.InvalidArgument, "invalid username format")
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	user, err := s.authUseCase.ValidateToken(ctx, token)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to validate token")
	}

	return &authv1.SignUpResponse{
		Token:  token,
		UserId: user.ID.String(),
	}, nil
}

// ValidateToken implements the ValidateToken RPC method
func (s *AuthService) ValidateToken(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
	if req.GetToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	user, err := s.authUseCase.ValidateToken(ctx, req.GetToken())
	if err != nil {
		switch err {
		case usecase.ErrInvalidToken:
			return &authv1.ValidateTokenResponse{IsValid: false}, status.Error(codes.InvalidArgument, "invalid token")
		case usecase.ErrUserNotFound:
			return &authv1.ValidateTokenResponse{IsValid: false}, status.Error(codes.NotFound, "user not found")
		default:
			return &authv1.ValidateTokenResponse{IsValid: false}, status.Error(codes.Internal, "internal server error")
		}
	}

	return &authv1.ValidateTokenResponse{
		IsValid:  true,
		UserId:   user.ID.String(),
		Email:    user.Email,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
	}, nil
}
