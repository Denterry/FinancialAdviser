package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	_errInvalidEmail    = errors.New("invalid email format")
	_errInvalidPassword = errors.New("password must be at least 8 characters long")
	_errInvalidUsername = errors.New("username must be between 3 and 32 characters")
	_errEmptyEmail      = errors.New("email is required")
	_errEmptyPassword   = errors.New("password is required")
	_errEmptyUsername   = errors.New("username is required")
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID `db:"id"           json:"id"`
	Email        string    `db:"email"        json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Username     string    `db:"username"     json:"username"`
	IsAdmin      bool      `db:"is_admin"     json:"is_admin"`
	CreatedAt    time.Time `db:"created_at"   json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"   json:"updated_at"`
}

// NewUser creates a new user with validation
func NewUser(email, password, username string, isAdmin bool) (*User, error) {
	if email == "" {
		return nil, _errEmptyEmail
	}
	if len(email) < 3 || len(email) > 254 {
		return nil, _errInvalidEmail
	}
	if password == "" {
		return nil, _errEmptyPassword
	}
	if len(password) < 8 {
		return nil, _errInvalidPassword
	}
	if username == "" {
		return nil, _errEmptyUsername
	}
	if len(username) < 3 || len(username) > 32 {
		return nil, _errInvalidUsername
	}

	hash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return &User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: hash,
		Username:     username,
		IsAdmin:      isAdmin,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// SetPassword sets the password for the user
func (u *User) SetPassword(password string) error {
	hash, err := HashPassword(password)
	if err != nil {
		return err
	}

	u.PasswordHash = hash
	return nil
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// CheckPassword checks if the provided password matches the hashed password
func (u *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) == nil
}
