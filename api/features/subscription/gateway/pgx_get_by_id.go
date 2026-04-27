package subscription_gateway

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (g *PgxSubscriptionGatewayAdapter) GetByID(ctx context.Context, filter GetByIDFilter) (*GetByIDOutput, error) {
	b := sq.
		Select(
			"id", "tenant_id", "user_id", "status", "subscription_group_id",
			"product_id", "quantity", "unit_price",
			"delivery_days", "notes",
			"total_cycles", "remaining_cycles", "cycle_discount", "cutoff_offset_minutes",
			"auto_renew", "items_total", "items_discount", "payment_amount", "created_at",
		).
		From("subscription.subscriptions").
		Where(sq.Eq{"id": filter.ID}).
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	if filter.TenantID != "" {
		b = b.Where(sq.Eq{"tenant_id": filter.TenantID})
	}

	if filter.UserID != "" {
		b = b.Where(sq.Eq{"user_id": filter.UserID})
	}

	sql, args, err := b.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := g.DB.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[GetByIDOutput])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return &row, err
}
