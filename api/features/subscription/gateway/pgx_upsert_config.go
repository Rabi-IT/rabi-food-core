package subscription_gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
)

func (g *PgxSubscriptionGatewayAdapter) UpsertConfig(
	ctx context.Context,
	tenantID string,
	input UpsertConfigInput,
) error {
	discountRules, err := json.Marshal(input.DiscountRules)
	if err != nil {
		return fmt.Errorf("failed to marshal discount rules: %w", err)
	}

	sql, args, err := sq.
		Insert("subscription.subscription_configs").
		Columns("tenant_id", "max_attempts_per_order", "discount_rules", "cutoff_offset_minutes", "is_open", "updated_at").
		Values(tenantID, input.MaxAttemptsPerOrder, discountRules, input.CutoffOffsetMinutes, input.IsOpen, time.Now().UTC()).
		Suffix(`ON CONFLICT (tenant_id) DO UPDATE SET
			max_attempts_per_order  = EXCLUDED.max_attempts_per_order,
			discount_rules          = EXCLUDED.discount_rules,
			cutoff_offset_minutes   = EXCLUDED.cutoff_offset_minutes,
			is_open                 = EXCLUDED.is_open,
			updated_at              = EXCLUDED.updated_at`).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = g.DB.Pool.Exec(ctx, sql, args...)

	return err
}
