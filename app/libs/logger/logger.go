package logger

import (
	"context"
	"os"
	"rabi-food-core/config"
	"time"

	"github.com/rs/zerolog"
)

type contextKey string

var (
	base      zerolog.Logger
	loggerKey = contextKey("loggerKey")
)

func init() {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if config.Env != "production" {
		base = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			NoColor:    true,
		}).With().Timestamp().Logger()
	} else {
		base = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
}

// Get retrieves the logger from the given context.
func Get(c context.Context) *zerolog.Logger {
	v := c.Value(loggerKey)
	if l, ok := v.(*zerolog.Logger); ok {
		return l
	}

	return L()
}

// L returns the base logger instance.
func L() *zerolog.Logger {
	return &base
}

func WithContext(c context.Context, logger zerolog.Logger) context.Context {
	return context.WithValue(c, loggerKey, &logger)
}
