package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	_defaultMaxPoolSize  = 10
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
	_defaultMinConns     = 2
	_defaultMaxRetries   = 3
	_defaultRetryDelay   = 3 * time.Second
	_defaultHealthCheck  = 15 * time.Second

	serviceName = "x-service"
	timeZone    = "UTC"
)

// Postgres represents PostgreSQL connection pool
type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration
	minConns     int
	maxRetries   int
	retryDelay   time.Duration
	healthCheck  time.Duration

	Pool *pgxpool.Pool
}

// New returns new Postgres instance
func New(url string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
		minConns:     _defaultMinConns,
		maxRetries:   _defaultMaxRetries,
		retryDelay:   _defaultRetryDelay,
		healthCheck:  _defaultHealthCheck,
	}

	// custom options
	for _, opt := range opts {
		opt(pg)
	}

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = int32(pg.maxPoolSize) //nolint:gosec // skip integer overflow conversion int -> int32
	poolConfig.MinConns = int32(pg.minConns)    //nolint:gosec // skip integer overflow conversion int -> int32
	poolConfig.HealthCheckPeriod = pg.healthCheck

	poolConfig.ConnConfig.RuntimeParams["application_name"] = serviceName
	poolConfig.ConnConfig.RuntimeParams["timezone"] = timeZone

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			if err = pg.Pool.Ping(context.Background()); err == nil {
				break
			}
		}

		log.Printf("Postgres is trying to connect, attempts left: %d, error: %v", pg.connAttempts, err)
		time.Sleep(pg.connTimeout)
		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
	}

	return pg, nil
}

// Close closes Postgres connection pool
func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

// GetPool returns Postgres connection pool
func (p *Postgres) GetPool() *pgxpool.Pool {
	return p.Pool
}

// Ping checks the database connection
func (p *Postgres) Ping(ctx context.Context) error {
	if p.Pool == nil {
		return fmt.Errorf("postgres - Ping - pool is nil")
	}
	return p.Pool.Ping(ctx)
}

// Stats returns pool statistics
func (p *Postgres) Stats() *pgxpool.Stat {
	if p.Pool == nil {
		return nil
	}
	return p.Pool.Stat()
}
