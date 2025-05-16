package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Services ServicesConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Log      LogConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type ServicesConfig struct {
	Auth         ServiceConfig
	Subscription ServiceConfig
	ML           ServiceConfig
}

type ServiceConfig struct {
	Host string
	Port string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	SigningMethod   string
	SigningKey      string
	VerificationKey string
}

type LogConfig struct {
	Level string
	File  string
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	config := &Config{
		Server: ServerConfig{
			Port:         viper.GetString("server.port"),
			ReadTimeout:  viper.GetDuration("server.read_timeout"),
			WriteTimeout: viper.GetDuration("server.write_timeout"),
		},
		Services: ServicesConfig{
			Auth: ServiceConfig{
				Host: viper.GetString("services.auth.host"),
				Port: viper.GetString("services.auth.port"),
			},
			Subscription: ServiceConfig{
				Host: viper.GetString("services.subscription.host"),
				Port: viper.GetString("services.subscription.port"),
			},
			ML: ServiceConfig{
				Host: viper.GetString("services.ml.host"),
				Port: viper.GetString("services.ml.port"),
			},
		},
		Redis: RedisConfig{
			Host:     viper.GetString("redis.host"),
			Port:     viper.GetString("redis.port"),
			Password: viper.GetString("redis.password"),
			DB:       viper.GetInt("redis.db"),
		},
		JWT: JWTConfig{
			Secret:          viper.GetString("jwt.secret"),
			AccessTokenTTL:  viper.GetDuration("jwt.access_token_ttl"),
			RefreshTokenTTL: viper.GetDuration("jwt.refresh_token_ttl"),
			SigningMethod:   viper.GetString("jwt.signing_method"),
			SigningKey:      viper.GetString("jwt.signing_key"),
			VerificationKey: viper.GetString("jwt.verification_key"),
		},
		Log: LogConfig{
			Level: viper.GetString("log.level"),
			File:  viper.GetString("log.file"),
		},
	}

	return config, nil
}
