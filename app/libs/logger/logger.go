package logger

import (
	"context"
	"os"
	"time"

	"github.com/Rabi-IT/rabi-food-core/config"

	"github.com/rs/zerolog"
)

var base zerolog.Logger

const (
	loggerKey = "loggerKey"
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
// If no logger is found, returns a pointer to a fresh copy of the base logger
// so that callers can safely call UpdateContext without mutating global state.
func Get(c context.Context) *zerolog.Logger {
	v := c.Value(loggerKey)
	if l, ok := v.(*zerolog.Logger); ok {
		return l
	}

	l := base.With().Logger()

	return &l
}

func New() zerolog.Context {
	return base.With()
}

// withContext stores the logger in the context.
func withContext(c context.Context, logger *zerolog.Logger) context.Context {
	return context.WithValue(c, loggerKey, logger)
}
