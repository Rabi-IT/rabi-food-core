package gateway

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

func (g *PgxTenantGatewayAdapter) CreateCustomer(ctx context.Context, tenantID, userID string) error {
	sql, args, err := sq.
		Insert("iam.tenant_customers").
		Columns("tenant_id", "user_id").
		Values(tenantID, userID).
		Suffix("ON CONFLICT (tenant_id, user_id) DO NOTHING").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = g.DB.Pool.Exec(ctx, sql, args...)

	return err
}
