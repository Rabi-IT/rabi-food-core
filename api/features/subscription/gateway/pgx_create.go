package subscription_gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (g *PgxSubscriptionGatewayAdapter) Create(ctx context.Context, input CreateInput) (string, error) {
	id := uuid.Must(uuid.NewV7()).String()

	items, err := json.Marshal(input.Items)
	if err != nil {
		return "", fmt.Errorf("failed to marshal subscription items: %w", err)
	}

	deliveryDays, err := json.Marshal(input.DeliveryDays)
	if err != nil {
		return "", fmt.Errorf("failed to marshal delivery days: %w", err)
	}

	weekdaysMask := uint8(0)
	for _, day := range input.DeliveryDays {
		weekdaysMask |= weekdayBit(day.Weekday)
	}

	now := time.Now().UTC()

	sql, args, err := sq.
		Insert("subscription.subscriptions").
		Columns("id", "root_subscription_id", "tenant_id", "user_id", "status",
			"items", "delivery_days", "delivery_weekdays_mask", "notes",
			"total_cycles", "remaining_cycles", "cycle_discount", "cutoff_offset_minutes",
			"auto_renew", "max_attempts_per_order",
			"items_total", "items_discount", "payment_amount", "payment_status",
			"external_payment_id", "created_at", "updated_at").
		Values(id, input.RootID, input.TenantID, input.UserID, input.Status,
			items, deliveryDays, weekdaysMask, input.Notes,
			input.TotalCycles, input.RemainingCycles, input.CycleDiscount, input.CutoffOffsetMinutes,
			input.AutoRenew, input.MaxAttemptsPerOrder,
			input.ItemsTotal, input.ItemsDiscount, input.PaymentAmount, input.PaymentStatus,
			input.ExternalPaymentID, now, now).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return "", err
	}

	_, err = g.DB.Pool.Exec(ctx, sql, args...)

	return id, err
}
