package gateway

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (g *PgxTenantGatewayAdapter) GetByID(id string) (*GetByIDOutput, error) {
	sql, args, err := sq.
		Select("id", "name").
		From("tenants").
		Where(sq.Eq{"id": id}).
		Limit(1).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := g.DB.Pool.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}

	out, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[GetByIDOutput])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &out, nil
}
