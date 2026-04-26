package subscription_gateway

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (g *PgxSubscriptionGatewayAdapter) GetConfig(ctx context.Context, tenantID string) (*GetConfigOutput, error) {
	sql, args, err := sq.
		Select("max_attempts_per_order", "discount_rules", "cutoff_offset_minutes", "order_lead_minutes", "is_open").
		From("subscription.subscription_configs").
		Where(sq.Eq{"tenant_id": tenantID}).
		Limit(1).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := g.DB.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[GetConfigOutput])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return &row, err
}
