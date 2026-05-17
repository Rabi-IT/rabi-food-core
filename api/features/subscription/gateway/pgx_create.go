package subscription_gateway

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (g *PgxSubscriptionGatewayAdapter) BulkCreate(ctx context.Context, inputs []CreateInput) ([]string, error) {
	ids := make([]string, 0, len(inputs))
	now := time.Now().UTC()

	qb := sq.
		Insert("subscription.subscriptions").
		Columns(
			"id", "root_subscription_id",
			"tenant_id", "user_id", "product_id", "quantity", "unit_price",
			"status", "delivery_days", "delivery_weekdays_mask", "notes",
			"total_cycles", "remaining_cycles", "cycle_discount", "cutoff_offset_minutes",
			"auto_renew", "max_attempts_per_order",
			"items_total", "items_discount", "payment_amount", "payment_status",
			"created_at", "updated_at",
		)

	for _, input := range inputs {
		id := uuid.Must(uuid.NewV7()).String()
		ids = append(ids, id)

		weekdaysMask := uint8(0)
		for _, day := range input.DeliveryDays {
			weekdaysMask |= weekdayBit(day.Weekday)
		}

		qb = qb.Values(
			id, input.RootID,
			input.TenantID, input.UserID, input.ProductID, input.Quantity, input.UnitPrice,
			input.Status, input.DeliveryDays, weekdaysMask, input.Notes,
			input.TotalCycles, input.RemainingCycles, input.CycleDiscount, input.CutoffOffsetMinutes,
			input.AutoRenew, input.MaxAttemptsPerOrder,
			input.ItemsTotal, input.ItemsDiscount, input.PaymentAmount, input.PaymentStatus,
			now, now,
		)
	}

	sql, args, err := qb.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	_, err = g.DB.Pool.Exec(context.WithoutCancel(ctx), sql, args...)

	return ids, err
}
