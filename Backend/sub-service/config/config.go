package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	// Config -.
	Config struct {
		App     App
		HTTP    HTTP
		Log     Log
		PG      PG
		RMQ     RMQ
		Metrics Metrics
		Swagger Swagger

		// Server configuration
		ServerPort int
		ServerHost string

		// Database configuration
		DBHost     string
		DBPort     int
		DBUser     string
		DBPassword string
		DBName     string
		DBSSLMode  string

		// JWT configuration
		JWTSecret        string
		JWTExpiration    time.Duration
		JWTRefreshExpiry time.Duration

		// Payment configuration
		StripeSecretKey      string
		StripeWebhookSecret  string
		StripeSuccessURL     string
		StripeCancelURL      string
		DefaultCurrency      string
		PaymentMethods       []string
		AutoRenewalEnabled   bool
		GracePeriodDays      int
		SubscriptionDuration int // in days

		// Logging configuration
		LogLevel string
		LogFile  string
	}

	// App -.
	App struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
	}

	// HTTP -.
	HTTP struct {
		Port           string `env:"HTTP_PORT,required"`
		UsePreforkMode bool   `env:"HTTP_USE_PREFORK_MODE" envDefault:"false"`
	}

	// Log -.
	Log struct {
		Level string `env:"LOG_LEVEL,required"`
	}

	// PG -.
	PG struct {
		PoolMax int    `env:"PG_POOL_MAX,required"`
		URL     string `env:"PG_URL,required"`
	}

	// RMQ -.
	RMQ struct {
		ServerExchange string `env:"RMQ_RPC_SERVER"`
		ClientExchange string `env:"RMQ_RPC_CLIENT"`
		URL            string `env:"RMQ_URL"`
	}

	Metrics struct {
		Enabled bool `env:"METRICS_ENABLED" envDefault:"true"`
	}

	Swagger struct {
		Enabled bool `env:"SWAGGER_ENABLED" envDefault:"false"`
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

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	config := &Config{
		// Server configuration
		ServerPort: getEnvAsInt("SERVER_PORT", 8080),
		ServerHost: getEnv("SERVER_HOST", "localhost"),

		// Database configuration
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "subscription_db"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		// JWT configuration
		JWTSecret:        getEnv("JWT_SECRET", ""),
		JWTExpiration:    time.Duration(getEnvAsInt("JWT_EXPIRATION", 24)) * time.Hour,
		JWTRefreshExpiry: time.Duration(getEnvAsInt("JWT_REFRESH_EXPIRY", 168)) * time.Hour,

		// Payment configuration
		StripeSecretKey:      getEnv("STRIPE_SECRET_KEY", ""),
		StripeWebhookSecret:  getEnv("STRIPE_WEBHOOK_SECRET", ""),
		StripeSuccessURL:     getEnv("STRIPE_SUCCESS_URL", "http://localhost:3000/success"),
		StripeCancelURL:      getEnv("STRIPE_CANCEL_URL", "http://localhost:3000/cancel"),
		DefaultCurrency:      getEnv("DEFAULT_CURRENCY", "USD"),
		PaymentMethods:       getEnvAsSlice("PAYMENT_METHODS", []string{"card"}),
		AutoRenewalEnabled:   getEnvAsBool("AUTO_RENEWAL_ENABLED", true),
		GracePeriodDays:      getEnvAsInt("GRACE_PERIOD_DAYS", 3),
		SubscriptionDuration: getEnvAsInt("SUBSCRIPTION_DURATION", 30),

		// Logging configuration
		LogLevel: getEnv("LOG_LEVEL", "info"),
		LogFile:  getEnv("LOG_FILE", "subscription.log"),
	}

	// Validate required fields
	if config.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	if config.StripeSecretKey == "" {
		return nil, fmt.Errorf("STRIPE_SECRET_KEY is required")
	}
	if config.StripeWebhookSecret == "" {
		return nil, fmt.Errorf("STRIPE_WEBHOOK_SECRET is required")
	}

	return config, nil
}

// Helper functions to get environment variables
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	if value, exists := os.LookupEnv(key); exists {
		// Split by comma and trim spaces
		values := strings.Split(value, ",")
		for i, v := range values {
			values[i] = strings.TrimSpace(v)
		}
		return values
	}
	return defaultValue
}
