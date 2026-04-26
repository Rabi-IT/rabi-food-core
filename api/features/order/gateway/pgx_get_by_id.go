package gateway

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (g *PgxOrderGatewayAdapter) GetByID(ctx context.Context, filter GetByIDFilter) (*GetByIDOutput, error) {
	b := sq.
		Select("id", "tenant_id", "code", "payment_status", "fulfillment_status",
			"delivery_status", "notes", "total_price", "items", "paid_at",
			"created_at", "external_payment_id").
		From("commerce.orders").
		Where("deleted_at IS NULL").
		Where(sq.Eq{"id": filter.ID}).
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	if filter.UserID != "" {
		b = b.Where(sq.Eq{"user_id": filter.UserID})
	}
	if filter.TenantID != "" {
		b = b.Where(sq.Eq{"tenant_id": filter.TenantID})
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
