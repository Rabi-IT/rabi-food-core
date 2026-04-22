package cron

import (
	"context"
	"time"

	"github.com/Rabi-IT/rabi-food-core/libs/logger"
	"github.com/robfig/cron/v3"
)

type Job struct {
	Schedule string
	Name     string
	Run      func(ctx context.Context) error
}

type Scheduler struct {
	cron *cron.Cron
}

func newCron(location *time.Location) *Scheduler {
	return &Scheduler{
		cron: cron.New(cron.WithLocation(location)),
	}
}

func (s *Scheduler) register(ctx context.Context, job Job) error {
	l := logger.Get(ctx).
		With().
		Str("job", job.Name).
		Str("schedule", job.Schedule).
		Logger()

	l.Info().Msg("registering cron job")

	_, err := s.cron.AddFunc(job.Schedule, func() {
		we := logger.Acquire()
		defer logger.Release(we)

		we.Event = job.Name

		ctx := logger.WithWideEvent(ctx, we)

		start := time.Now()
		err := job.Run(ctx)
		latency := time.Since(start).Milliseconds()
		we.FinishCronJob(err, latency)

		if err != nil {
			l.Error().EmbedObject(we).Msg("cron job failed")

			return
		}

		l.Info().EmbedObject(we).Msg("cron job completed")
	})

	if err != nil {
		l.Error().Err(err).Msg("failed to register cron job")
	}

	return err
}

func (s *Scheduler) start(ctx context.Context) {
	logger.Get(ctx).Info().Msg("starting cron")
	s.cron.Start()
}
