package mapper

import (
	"time"

	"github.com/Denterry/FinancialAdviser/Backend/auth-service/internal/entity"
	authv1 "github.com/Denterry/FinancialAdviser/Backend/auth-service/pkg/pb/auth/v1"
	"github.com/google/uuid"
)

// UserToProto converts user entity to protobuf message
func UserToProto(u *entity.User) *authv1.User {
	return &authv1.User{
		Id:        u.ID.String(),
		Email:     u.Email,
		Username:  u.Username,
		IsAdmin:   u.IsAdmin,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}

// UserFromProto creates user entity from protobuf message
func UserFromProto(data *authv1.User) (*entity.User, error) {
	createdAt, err := time.Parse(time.RFC3339, data.CreatedAt)
	if err != nil {
		return nil, err
	}

	updatedAt, err := time.Parse(time.RFC3339, data.UpdatedAt)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(data.Id)
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:        id,
		Email:     data.Email,
		Username:  data.Username,
		IsAdmin:   data.IsAdmin,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

// UserToMap converts user entity to map for database operations
func UserToMap(u *entity.User) map[string]interface{} {
	return map[string]interface{}{
		"id":            u.ID.String(),
		"email":         u.Email,
		"password_hash": u.PasswordHash,
		"username":      u.Username,
		"is_admin":      u.IsAdmin,
		"created_at":    u.CreatedAt,
		"updated_at":    u.UpdatedAt,
	}
}

// UserFromMap creates user entity from database row
func UserFromMap(data map[string]interface{}) (*entity.User, error) {
	id, err := uuid.Parse(data["id"].(string))
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:           id,
		Email:        data["email"].(string),
		PasswordHash: data["password_hash"].(string),
		Username:     data["username"].(string),
		IsAdmin:      data["is_admin"].(bool),
		CreatedAt:    data["created_at"].(time.Time),
		UpdatedAt:    data["updated_at"].(time.Time),
	}, nil
}
