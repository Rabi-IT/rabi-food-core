package gateway

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

func (g *PgxTenantGatewayAdapter) CreateMember(ctx context.Context, input CreateMemberInput) error {
	sql, args, err := sq.
		Insert("iam.tenant_members").
		Columns("tenant_id", "user_id", "role").
		Values(input.TenantID, input.UserID, input.Role).
		Suffix("ON CONFLICT (tenant_id, user_id) DO UPDATE SET role = EXCLUDED.role").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = g.DB.Pool.Exec(ctx, sql, args...)

	return err
}
