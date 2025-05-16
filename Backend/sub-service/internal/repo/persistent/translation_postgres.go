package persistent

import (
	"context"
	"fmt"

	"github.com/Denterry/FinancialAdviser/internal/entity"
	"github.com/Denterry/FinancialAdviser/internal/repo"
	"github.com/Denterry/FinancialAdviser/pkg/postgres"
)

const _defaultEntityCap = 64

// TranslationRepo -.
type TranslationRepo struct {
	*postgres.Postgres
}

// New -.
func New(pg *postgres.Postgres) *TranslationRepo {
	return &TranslationRepo{pg}
}

// GetHistory -.
func (r *TranslationRepo) GetHistory(ctx context.Context) ([]entity.Translation, error) {
	sql, _, err := r.Builder.
		Select("source, destination, original, translation").
		From("history").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("TranslationRepo - GetHistory - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("TranslationRepo - GetHistory - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	entities := make([]entity.Translation, 0, _defaultEntityCap)

	for rows.Next() {
		e := entity.Translation{}

		err = rows.Scan(&e.Source, &e.Destination, &e.Original, &e.Translation)
		if err != nil {
			return nil, fmt.Errorf("TranslationRepo - GetHistory - rows.Scan: %w", err)
		}

		entities = append(entities, e)
	}

	return entities, nil
}

// Store -.
func (r *TranslationRepo) Store(ctx context.Context, t entity.Translation) error {
	sql, args, err := r.Builder.
		Insert("history").
		Columns("source, destination, original, translation").
		Values(t.Source, t.Destination, t.Original, t.Translation).
		ToSql()
	if err != nil {
		return fmt.Errorf("TranslationRepo - Store - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("TranslationRepo - Store - r.Pool.Exec: %w", err)
	}

	return nil
}

// CreateUser - add new user into database, table `users`
func (r *TranslationRepo) CreateUser(ctx context.Context, u *repo.User) (*repo.User, error) {
	res, err := ms.db.NamedExecContext(ctx, "INSERT INTO users (name, email, password, is_admin) VALUES (:name, :email, :password, :is_admin)", u)
	if err != nil {
		return nil, fmt.Errorf("error inserting user: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert ID: %w", err)
	}
	u.ID = id

	return u, nil
}

// GetUser - get user by email from `users` table
func (r *TranslationRepo) GetUser(ctx context.Context, email string) (*repo.User, error) {
	var u repo.User
	err := ms.db.GetContext(ctx, &u, "SELECT * FROM users WHERE email=?", email)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return &u, nil
}

// ListUsers - get list of all users from `users` table
func (r *TranslationRepo) ListUsers(ctx context.Context) ([]*repo.User, error) {
	var users []*repo.User
	err := ms.db.SelectContext(ctx, &users, "SELECT * FROM users")
	if err != nil {
		return nil, fmt.Errorf("error listing users: %w", err)
	}

	return users, nil
}

// UpdateUser - update exact user by id in `users` table
func (r *TranslationRepo) UpdateUser(ctx context.Context, u *repo.User) (*repo.User, error) {
	_, err := ms.db.NamedExecContext(ctx, "UPDATE users SET name=:name, email=:email, password=:password, is_admin=:is_admin, updated_at=:updated_at WHERE id=:id", u)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return u, nil
}

// DeleteUser - delete exact user by id from `users` table
func (r *TranslationRepo) DeleteUser(ctx context.Context, id int64) error {
	_, err := ms.db.ExecContext(ctx, "DELETE FROM users WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	return nil
}
