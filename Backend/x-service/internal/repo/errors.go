package repo

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrDuplicateTweet = errors.New("duplicate tweet")
)

var (
	UniqueViolationErr = pgconn.PgError{
		Code: "23505",
	}
)
