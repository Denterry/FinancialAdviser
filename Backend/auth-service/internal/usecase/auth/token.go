package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/Denterry/FinancialAdviser/Backend/auth-service/internal/entity"
	"github.com/Denterry/FinancialAdviser/Backend/auth-service/internal/usecase"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// makeToken creates a JWT-token for user
func (uc *UseCase) makeToken(user *entity.User, currentTime time.Time) (string, error) {
	claims := jwt.MapClaims{
		"sub":      user.ID.String(),
		"email":    user.Email,
		"username": user.Username,
		"exp":      currentTime.Add(uc.tokenTTL).Unix(),
		"iat":      currentTime.Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(uc.jwtKey)
}

// ValidateToken validates a JWT token and returns the associated user
func (uc *UseCase) ValidateToken(ctx context.Context, tokenString string) (*entity.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, usecase.ErrInvalidToken
		}
		return uc.jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, usecase.ErrInvalidToken
	}

	// get claims & validate
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, usecase.ErrInvalidToken
	}

	if exp, ok := claims["exp"].(float64); !ok || float64(time.Now().Unix()) > exp {
		return nil, usecase.ErrInvalidToken
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return nil, usecase.ErrInvalidToken
	}
	userID, err := uuid.Parse(subject)
	if err != nil {
		return nil, usecase.ErrInvalidToken
	}

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("uc.userRepo.GetByID(): %w", err)
	}
	if user == nil {
		return nil, usecase.ErrUserNotFound
	}

	return user, nil
}
