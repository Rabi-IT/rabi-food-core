package subscription_gateway

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

// CreateDeliveries inserts multiple subscription deliveries in a single batch.
// The operation is idempotent: duplicate (subscription_id, scheduled_date) pairs are
// silently ignored thanks to the unique constraint on the table.
func (g *PgxSubscriptionGatewayAdapter) CreateDeliveries(ctx context.Context, inputs []CreateDeliveryInput) error {
	if len(inputs) == 0 {
		return nil
	}

	now := time.Now().UTC()

	b := sq.
		Insert("subscription_deliveries").
		Columns("id", "subscription_id", "scheduled_date", "start_hour", "end_hour",
			"cutoff_at", "status", "max_delivery_attempts", "created_at", "updated_at").
		Suffix("ON CONFLICT ON CONSTRAINT subscription_delivery_unique DO NOTHING").
		PlaceholderFormat(sq.Dollar)

	for _, in := range inputs {
		b = b.Values(
			uuid.Must(uuid.NewV7()).String(),
			in.SubscriptionID,
			in.ScheduledDate,
			in.StartHour,
			in.EndHour,
			in.CutoffAt,
			in.Status,
			in.MaxAttempts,
			now,
			now,
		)
	}

	sql, args, err := b.ToSql()
	if err != nil {
		return err
	}

	_, err = g.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to create subscription deliveries: %w", err)
	}

	return nil
}
