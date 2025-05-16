package entity

import "errors"

var (
	ErrEmptyText   = errors.New("tweet text must not be empty")
	ErrTooLongText = errors.New("tweet text exceeds 4 096 characters")
)
