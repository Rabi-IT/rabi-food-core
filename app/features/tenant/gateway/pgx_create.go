package gateway

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (g *PgxTenantGatewayAdapter) Create(ctx context.Context, input CreateInput) (string, error) {
	id := uuid.Must(uuid.NewV7()).String()

	sql, args, err := sq.
		Insert("iam.tenants").
		Columns("id", "name").
		Values(id, input.Name).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return "", err
	}

	_, err = g.DB.Pool.Exec(context.WithoutCancel(ctx), sql, args...)

	return id, err
}
