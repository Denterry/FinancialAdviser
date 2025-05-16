package postgres

import "time"

// Option -.
type Option func(*Postgres)

// MaxPoolSize -.
func MaxPoolSize(size int) Option {
	return func(c *Postgres) {
		c.maxPoolSize = size
	}
}

// ConnAttempts -.
func ConnAttempts(attempts int) Option {
	return func(c *Postgres) {
		c.connAttempts = attempts
	}
}

// ConnTimeout -.
func ConnTimeout(timeout time.Duration) Option {
	return func(c *Postgres) {
		c.connTimeout = timeout
	}
}

// MinConns -.
func MinConns(minConns int) Option {
	return func(c *Postgres) {
		c.minConns = minConns
	}
}

// MaxRetries -.
func MaxRetries(maxRetries int) Option {
	return func(c *Postgres) {
		c.maxRetries = maxRetries
	}
}

// RetryDelay -.
func RetryDelay(retryDelay time.Duration) Option {
	return func(c *Postgres) {
		c.retryDelay = retryDelay
	}
}
