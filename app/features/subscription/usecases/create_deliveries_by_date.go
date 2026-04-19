package usecases

import (
	"context"
	"time"

	"github.com/Rabi-IT/rabi-food-core/features/subscription"
	g "github.com/Rabi-IT/rabi-food-core/features/subscription/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
)

type CreateDeliveriesByDateInput struct {
	SubscriptionID string
	Date           time.Time
}

func (s *SubscriptionCase) CreateDeliveriesByDate(ctx context.Context, input CreateDeliveriesByDateInput) error {
	date := time.Date(input.Date.Year(), input.Date.Month(), input.Date.Day(),
		0, 0, 0, 0, input.Date.Location())

	weekday := date.Weekday()
	subs, err := s.gateway.List(ctx, g.ListFilter{
		ID:                 input.SubscriptionID,
		Status:             subscription.StatusActive,
		RemainingCyclesGTE: 1,
		Weekday:            weekday,
	})
	if err != nil {
		return err
	}

	deliveries := make([]g.CreateDeliveryInput, 0, len(subs))
	for _, sub := range subs {
		var day g.DeliveryDay
		for _, d := range sub.DeliveryDays {
			if d.Weekday == uint8(weekday) {
				day = d
				break
			}
		}

		scheduledDate := date.Add(time.Duration(day.StartHour) * time.Hour)

		cutoffAt := scheduledDate.Add(-time.Duration(sub.CutoffOffsetMinutes) * time.Minute)

		deliveries = append(deliveries, g.CreateDeliveryInput{
			SubscriptionID: sub.ID,
			ScheduledDate:  scheduledDate,
			StartHour:      day.StartHour,
			EndHour:        day.EndHour,
			CutoffAt:       cutoffAt,
			Status:         subscription.DeliveryScheduled,
			MaxAttempts:    sub.MaxAttemptsPerOrder,
		})
	}

	logger.
		GetWideEvent(ctx).
		TrackDeliveryScheduling(
			date.Format("2006-01-02"),
			weekday.String(),
			date.Location().String(),
			len(deliveries),
		)

	return s.gateway.CreateDeliveries(ctx, deliveries)
}
