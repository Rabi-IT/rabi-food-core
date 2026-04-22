package gateway

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (g *PgxTenantGatewayAdapter) GetCustomer(ctx context.Context, filter GetCustomerFilter) (*GetCustomerOutput, error) {
	q := sq.
		Select("user_id", "tenant_id").
		From("iam.tenant_customers").
		PlaceholderFormat(sq.Dollar)

	if filter.UserID != "" {
		q = q.Where(sq.Eq{"user_id": filter.UserID})
	}

	if filter.TenantID != "" {
		q = q.Where(sq.Eq{"tenant_id": filter.TenantID})
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	var out GetCustomerOutput
	err = g.DB.Pool.QueryRow(ctx, sql, args...).Scan(&out.UserID, &out.TenantID)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return &out, err
}
