package service

import (
	"context"
	"errors"
	"time"

	authpb "github.com/Denterry/FinancialAdviser/Backend/auth-service/pkg/pb/auth/v1"
	"github.com/Denterry/FinancialAdviser/Backend/gateway-service/internal/entity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	Client authpb.AuthServiceClient
}

func NewAuthService(conn *grpc.ClientConn) *AuthService {
	return &AuthService{Client: authpb.NewAuthServiceClient(conn)}
}

// Validate a JWT and translate proto → entity.
func (a *AuthService) ValidateToken(ctx context.Context, token string) (*entity.Claims, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := a.Client.ValidateToken(ctx, &authpb.ValidateTokenRequest{Token: token})
	if err != nil {
		return nil, err
	}
	if !resp.IsValid {
		return nil, errors.New("token invalid")
	}

	return &entity.Claims{
		UserID:   resp.UserId,
		Email:    resp.Email,
		Username: resp.Username,
		IsAdmin:  resp.IsAdmin,
	}, nil
}

// Refresh takes a *refresh* token, validates it, then issues a new access
// token by simply calling SignIn again (one could design a dedicated RPC).
func (a *AuthService) Refresh(ctx context.Context, refresh string) (bool, string, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	vr, err := a.Client.ValidateToken(ctx, &authpb.ValidateTokenRequest{Token: refresh})
	if err != nil {
		return false, "", err
	}
	if !vr.IsValid {
		return false, "", nil
	}

	// In a real world flow we'd have a Refresh RPC.  For demo we re-sign-in.
	sr, err := a.Client.SignIn(ctx, &authpb.SignInRequest{
		Email:    vr.Email,
		Password: "", // passwordless refresh – the server trusts the old token
	})
	if err != nil {
		return true, "", err
	}
	return true, sr.Token, nil
}

// Helper to translate gRPC error → bool unauthorised.
func isUnauth(err error) bool {
	return status.Code(err).String() == "Unauthenticated"
}
