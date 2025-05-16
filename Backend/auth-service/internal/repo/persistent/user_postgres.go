package persistent

import (
	"context"
	"fmt"
	"time"

	"github.com/Denterry/FinancialAdviser/Backend/auth-service/internal/entity"
	"github.com/Denterry/FinancialAdviser/Backend/auth-service/pkg/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const (
	_defaultEntityCap = 64
)

// UserRepository implements interface for user data operations
type UserRepository struct {
	*postgres.Postgres
}

// NewUserPostgres creates a new instance of UserPostgres
func NewUserPostgres(pg *postgres.Postgres) *UserRepository {
	return &UserRepository{pg}
}

// Create creates a new user in the database
func (r *UserRepository) Create(ctx context.Context, user *entity.User) (uuid.UUID, error) {
	var returnedID uuid.UUID
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	now := time.Now().UTC()
	user.CreatedAt = now
	returnedID = user.ID

	const query = `
		INSERT INTO users (id, email, password_hash, username, is_admin, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id
	`
	err := r.Pool.QueryRow(ctx, query,
		user.ID,           // $1
		user.Email,        // $2
		user.PasswordHash, // $3
		user.Username,     // $4
		user.IsAdmin,      // $5
		user.CreatedAt,    // $6
		user.UpdatedAt,    // $7
	).Scan(&returnedID)
	if err != nil {
		return returnedID, fmt.Errorf("UserRepository - Create - r.Pool.Exec: %w", err)
	}

	return returnedID, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	const query = `
		SELECT id, email, password_hash, username, is_admin, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	var user entity.User
	err := r.Pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Username,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("UserRepository - GetByEmail - r.Pool.QueryRow: %w", err)
	}

	return &user, nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	const query = `
		SELECT id, email, password_hash, username, is_admin, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var user entity.User
	err := r.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Username,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("UserRepository - GetByID - r.Pool.QueryRow: %w", err)
	}

	return &user, nil
}

// List retrieves all users
func (r *UserRepository) List(ctx context.Context) ([]*entity.User, error) {
	const query = `
		SELECT id, email, password_hash, username, is_admin, created_at, updated_at
		FROM users
	`
	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("UserRepository - List - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	users := make([]*entity.User, 0, _defaultEntityCap)

	for rows.Next() {
		var user entity.User
		err = rows.Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.Username,
			&user.IsAdmin,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("UserRepository - List - rows.Scan: %w", err)
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("UserRepository - List - rows.Err: %w", err)
	}

	return users, nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	const query = `
		UPDATE users
		SET email = $1,
			password_hash = $2,
			username = $3,
			is_admin = $4,
			updated_at = $5
		WHERE id = $6
	`
	res, err := r.Pool.Exec(ctx, query,
		user.Email,        // $1
		user.PasswordHash, // $2
		user.Username,     // $3
		user.IsAdmin,      // $4
		user.UpdatedAt,    // $5
		user.ID,           // $6
	)
	if err != nil {
		return fmt.Errorf("UserRepository - Update - r.Pool.Exec: %w", err)
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("UserRepository - Update - r.Pool.Exec: user not found")
	}

	return nil
}

// Delete deletes a user by ID
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM users WHERE id = $1`
	_, err := r.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("UserRepository - Delete - r.Pool.Exec: %w", err)
	}

	return nil
}
