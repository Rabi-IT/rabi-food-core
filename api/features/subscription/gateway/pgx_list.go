package subscription_gateway

import (
	"context"
	"encoding/json"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

func (g *PgxSubscriptionGatewayAdapter) List(ctx context.Context, filter ListFilter) ([]ListOutput, error) {
	qb := sq.
		Select("id", "delivery_days", "cutoff_offset_minutes", "max_attempts_per_order").
		From("subscription.subscriptions").
		Where(sq.Eq{"status": filter.Status}).
		Where("remaining_cycles >= ?", filter.RemainingCyclesGTE).
		Where("delivery_weekdays_mask & ? != 0", 1<<uint8(filter.Weekday))

	if filter.ID != "" {
		qb = qb.Where(sq.Eq{"id": filter.ID})
	}

	query, args, err := qb.
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build schedulable query: %w", err)
	}

	rows, err := g.DB.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query schedulable subscriptions: %w", err)
	}
	defer rows.Close()

	var subs []ListOutput
	for rows.Next() {
		var (
			s            ListOutput
			deliveryJSON json.RawMessage
		)
		if err := rows.Scan(&s.ID, &deliveryJSON, &s.CutoffOffsetMinutes, &s.MaxAttemptsPerOrder); err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		if err := json.Unmarshal(deliveryJSON, &s.DeliveryDays); err != nil {
			return nil, fmt.Errorf("failed to unmarshal delivery_days for subscription %s: %w", s.ID, err)
		}
		subs = append(subs, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating subscriptions: %w", err)
	}

	return subs, nil
}
