package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type (
	// Config -.
	Config struct {
		App     App
		HTTP    HTTP
		GRPC    GRPC
		Log     Log
		PG      PG
		Metrics Metrics
		Swagger Swagger
		JWT     JWT
		TLS     TLS
	}

	// App -.
	App struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
	}

	// HTTP -.
	HTTP struct {
		Port            string        `env:"HTTP_PORT" envDefault:"8080"`
		ShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" envDefault:"5s"`
		ReadTimeout     time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"5s"`
		WriteTimeout    time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"5s"`
		IdleTimeout     time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"60s"`
	}

	// GRPC -.
	GRPC struct {
		Port                 string        `env:"GRPC_PORT" envDefault:"9090"`
		ShutdownTimeout      time.Duration `env:"GRPC_SHUTDOWN_TIMEOUT" envDefault:"5s"`
		MaxConcurrentStreams uint32        `env:"GRPC_MAX_CONCURRENT_STREAMS" envDefault:"100"`
		MaxConnectionIdle    time.Duration `env:"GRPC_MAX_CONNECTION_IDLE" envDefault:"30s"`
	}

	// Log -.
	Log struct {
		Level  string `env:"LOG_LEVEL,required"`
		Format string `env:"LOG_FORMAT" envDefault:"json"`
		Output string `env:"LOG_OUTPUT" envDefault:"stdout"`
	}

	// PG -.
	PG struct {
		URL        string        `env:"PG_URL,required"`
		PoolMax    int           `env:"PG_POOL_MAX,required"`
		MinConns   int           `env:"PG_MIN_CONNS" envDefault:"2"`
		MaxRetries int           `env:"PG_MAX_RETRIES" envDefault:"3"`
		RetryDelay time.Duration `env:"PG_RETRY_DELAY" envDefault:"3s"`
	}

	// Metrics -.
	Metrics struct {
		Enabled bool   `env:"METRICS_ENABLED" envDefault:"true"`
		Port    string `env:"METRICS_PORT" envDefault:"9091"`
		Path    string `env:"METRICS_PATH" envDefault:"/metrics"`
	}

	// Swagger -.
	Swagger struct {
		Enabled bool   `env:"SWAGGER_ENABLED" envDefault:"false"`
		Path    string `env:"SWAGGER_PATH" envDefault:"/swagger"`
		Version string `env:"SWAGGER_VERSION" envDefault:"1.0.0"`
	}

	// JWT -.
	JWT struct {
		Secret     string `env:"JWT_SECRET,required"`
		AccessTTL  int    `env:"JWT_ACCESS_TTL_MINUTES,required"`
		RefreshTTL int    `env:"JWT_REFRESH_TTL_DAYS,required"`
		Issuer     string `env:"JWT_ISSUER" envDefault:"auth-service"`
		Algorithm  string `env:"JWT_ALGORITHM" envDefault:"HS256"`
	}

	// TLS -.
	TLS struct {
		CertFile string `env:"TLS_CERT_FILE"`
		KeyFile  string `env:"TLS_KEY_FILE"`
	}
)

// NewConfig returns app config
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
