package gateway

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (g *PgxTenantGatewayAdapter) GetBySlug(ctx context.Context, slug string) (*GetBySlugOutput, error) {
	sql, args, err := sq.
		Select("id", "name", "slug", "language").
		From("iam.tenants").
		Where(sq.Eq{"slug": slug}).
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

	out, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[GetBySlugOutput])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &out, nil
}
