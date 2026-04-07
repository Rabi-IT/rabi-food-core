package gateway

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (g *PgxCategoryGatewayAdapter) GetByID(filter GetByIDFilter) (*GetByIDOutput, error) {
	b := sq.
		Select("id", "tenant_id", "name", "description").
		From("categories").
		Where(sq.Eq{"id": filter.ID}).
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	if filter.TenantID != "" {
		b = b.Where(sq.Eq{"tenant_id": filter.TenantID})
	}

	sql, args, err := b.ToSql()
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
