package gateway

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (g *PgxTenantGatewayAdapter) Create(ctx context.Context, input CreateInput) (string, error) {
	id := uuid.Must(uuid.NewV7()).String()

	tenantSQL, tenantArgs, err := sq.
		Insert("iam.tenants").
		Columns("id", "name").
		Values(id, input.Name).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return "", err
	}

	wctx := context.WithoutCancel(ctx)

	tx, err := g.DB.Pool.Begin(wctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(context.Background())

	if _, err = tx.Exec(wctx, tenantSQL, tenantArgs...); err != nil {
		return "", err
	}

	memberSQL, memberArgs, err := sq.
		Insert("iam.tenant_members").
		Columns("tenant_id", "user_id", "role").
		Values(id, input.InitialMember.UserID, input.InitialMember.Role).
		Suffix("ON CONFLICT (tenant_id, user_id) DO UPDATE SET role = EXCLUDED.role").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return "", err
	}

	if _, err = tx.Exec(wctx, memberSQL, memberArgs...); err != nil {
		return "", err
	}

	return id, tx.Commit(wctx)
}
