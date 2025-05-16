package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type (
	// Config -.
	Config struct {
		App       App
		HTTP      HTTP
		GRPC      GRPC
		Log       Log
		PG        PG
		Metrics   Metrics
		Swagger   Swagger
		XProvider XProvider
		TLS       TLS
	}

	// App -.
	App struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
	}

	// HTTP -.
	HTTP struct {
		Port            string        `env:"HTTP_PORT,required"`
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

	// XProvider -.
	XProvider struct {
		Type string `env:"X_PROVIDER_TYPE" envDefault:"scraper"`
		XAPI
		XScraper
	}
	// XAPI -.
	XAPI struct {
		BaseURL           string   `env:"X_API_BASE_URL,required"`
		BearerToken       string   `env:"X_API_BEARER_TOKEN,required"`
		StreamRules       []string `env:"X_API_STREAM_RULES" envSeparator:","`
		ConsumerKey       string   `env:"X_API_CONSUMER_KEY"`
		ConsumerSecret    string   `env:"X_API_CONSUMER_SECRET"`
		AccessToken       string   `env:"X_API_ACCESS_TOKEN"`
		AccessTokenSecret string   `env:"X_API_ACCESS_TOKEN_SECRET"`
	}
	// XScraper -.
	XScraper struct {
		DelaySec     int    `env:"X_SCRAPER_DELAY"      envDefault:"2"`
		UseAppLogin  bool   `env:"X_SCRAPER_APP_LOGIN"  envDefault:"true"`
		Username     string `env:"X_SCRAPER_USER"`
		Password     string `env:"X_SCRAPER_PASS"`
		Email        string `env:"X_SCRAPER_EMAIL"`
		Proxy        string `env:"X_SCRAPER_PROXY"`
		IncludeReply bool   `env:"X_SCRAPER_REPLIES"    envDefault:"false"`
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
