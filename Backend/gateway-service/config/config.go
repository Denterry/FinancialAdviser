package config

import (
	"fmt"
	"net"
	"time"

	"github.com/caarlos0/env/v11"
)

type (
	// Config -..
	Config struct {
		App      App
		HTTP     HTTP
		Log      Log
		Metrics  Metrics
		Swagger  Swagger
		Services Services
		Redis    Redis
		JWT      JWT
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
		GIN             GIN
	}

	// GIN -.
	GIN struct {
		Mode      string `env:"GIN_MODE" envDefault:"debug"`
		RateLimit int    `env:"GIN_RATE_LIMIT" envDefault:"100"`
	}

	// Log -.
	Log struct {
		Level  string `env:"LOG_LEVEL,required"`
		Format string `env:"LOG_FORMAT" envDefault:"json"`
		Output string `env:"LOG_OUTPUT" envDefault:"stdout"`
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

	// Services -..
	Services struct {
		Auth  GRPCService
		Brain GRPCService
		ML    GRPCService
		Sub   GRPCService
	}

	// GRPCService -..
	GRPCService struct {
		Host string `env:"HOST,required"`
		Port string `env:"PORT,required"`
	}

	// Redis -..
	Redis struct {
		Host     string `env:"REDIS_HOST,required"`
		Port     string `env:"REDIS_PORT,required"`
		Password string `env:"REDIS_PASSWORD"`
		DB       int    `env:"REDIS_DB" envDefault:"0"`
	}

	// JWT -..
	JWT struct {
		Secret          string        `env:"JWT_SECRET,required"`
		AccessTokenTTL  time.Duration `env:"JWT_ACCESS_TOKEN_TTL,required"`
		RefreshTokenTTL time.Duration `env:"JWT_REFRESH_TOKEN_TTL,required"`
		VerificationKey string        `env:"JWT_VERIFICATION_KEY,required"`
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

func (s GRPCService) Addr() string {
	return net.JoinHostPort(s.Host, s.Port)
}
