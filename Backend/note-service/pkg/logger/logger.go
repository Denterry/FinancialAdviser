package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
}

type logger struct {
	logger zerolog.Logger
}

func New(level string) Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	var logLevel zerolog.Level
	switch level {
	case "debug":
		logLevel = zerolog.DebugLevel
	case "info":
		logLevel = zerolog.InfoLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "error":
		logLevel = zerolog.ErrorLevel
	default:
		logLevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(logLevel)

	return &logger{
		logger: zerolog.New(os.Stdout).With().Timestamp().Logger(),
	}
}

func (l *logger) Debug(msg string, args ...interface{}) {
	l.logger.Debug().Msgf(msg, args...)
}

func (l *logger) Info(msg string, args ...interface{}) {
	l.logger.Info().Msgf(msg, args...)
}

func (l *logger) Warn(msg string, args ...interface{}) {
	l.logger.Warn().Msgf(msg, args...)
}

func (l *logger) Error(msg string, args ...interface{}) {
	l.logger.Error().Msgf(msg, args...)
}

func (l *logger) Fatal(msg string, args ...interface{}) {
	l.logger.Fatal().Msgf(msg, args...)
}
