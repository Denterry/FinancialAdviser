package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/Denterry/FinancialAdviser/Backend/auth-service/internal/entity"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// makeToken creates a JWT-token for user
func (uc *UseCase) makeToken(user *entity.User, currentTime time.Time) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   user.ID.String(),
		IssuedAt:  jwt.NewNumericDate(currentTime),
		ExpiresAt: jwt.NewNumericDate(currentTime.Add(uc.tokenTTL)),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(uc.jwtKey)
}

// ValidateToken validates a JWT token and returns the associated user
func (uc *UseCase) ValidateToken(ctx context.Context, tokenString string) (*entity.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return uc.jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	// get claims & check
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, ErrInvalidToken
	}

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("uc.userRepo.GetByID(): %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}
