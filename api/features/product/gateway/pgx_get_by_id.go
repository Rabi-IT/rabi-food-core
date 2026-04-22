package gateway

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (g *PgxProductGatewayAdapter) GetByID(ctx context.Context, filter GetByIDFilter) (*GetByIDOutput, error) {
	b := sq.
		Select("p.id", "p.tenant_id", "p.name", "p.description", "p.photo",
			"p.category_id", "c.name AS category_name", "p.unit", "p.price", "p.is_active").
		From("catalog.products p").
		Join("catalog.categories c ON c.id = p.category_id").
		Where("p.deleted_at IS NULL").
		Where(sq.Eq{"p.id": filter.ID}).
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	if filter.TenantID != "" {
		b = b.Where(sq.Eq{"p.tenant_id": filter.TenantID})
	}

	sql, args, err := b.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := g.DB.Pool.Query(ctx, sql, args...)
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
