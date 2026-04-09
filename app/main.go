package main

import (
	"context"
	"time"

	"github.com/Rabi-IT/rabi-food-core/config"
	appCron "github.com/Rabi-IT/rabi-food-core/libs/cron"
	"github.com/Rabi-IT/rabi-food-core/libs/database"
	"github.com/Rabi-IT/rabi-food-core/libs/di"
	"github.com/Rabi-IT/rabi-food-core/libs/http"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"

	"github.com/samber/do"
)

func main() {
	time.Local = time.UTC
	l := logger.New().Logger()
	ctx := logger.WithContext(context.Background(), &l)

	l.
		Info().
		Str("env", config.Env.String()).
		Msg("Starting Rabi Food Core Server...")

	injector := di.Injector()

	db := do.MustInvoke[database.Database](injector)

	err := db.Start(ctx)
	if err != nil {
		l.Fatal().Err(err).Msg("Failed to start database")
	}

	// Cron scheduler
	scheduler := do.MustInvoke[*appCron.Scheduler](injector)
	scheduler.Start(ctx)

	httpServer := do.MustInvoke[http.HTTPServer](injector)

	err = httpServer.Start()
	if err != nil {
		l.Fatal().Err(err).Msg("Failed to start HTTP server")
	}
}
