package cron

import (
	"context"
	"time"

	"github.com/Rabi-IT/rabi-food-core/features/subscription/usecases"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
	"github.com/samber/do"
)

func Start(ctx context.Context, injector *do.Injector) {
	location := time.UTC
	scheduler := newCron(location)

	subscriptionCase := do.MustInvoke[*usecases.SubscriptionCase](injector)

	err := scheduler.register(ctx, Job{
		Schedule: "0 0 * * *",
		Name:     "create-deliveries",
		Run: func(ctx context.Context) error {
			tomorrow := time.Now().In(location).AddDate(0, 0, 1)

			return subscriptionCase.CreateDeliveriesByDate(ctx, usecases.CreateDeliveriesByDateInput{
				Date: tomorrow,
			})
		},
	})

	if err != nil {
		logger.Get(ctx).Fatal().Err(err).Msg("Failed to create delivery scheduler")
	}

	scheduler.start(ctx)
}
