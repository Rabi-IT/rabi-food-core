package subscription_gateway

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
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

	subs, err := pgx.CollectRows(rows, pgx.RowToStructByName[ListOutput])
	if err != nil {
		return nil, fmt.Errorf("failed to scan subscriptions: %w", err)
	}

	return subs, nil
}
