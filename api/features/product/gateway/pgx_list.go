package gateway

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (g *PgxProductGatewayAdapter) List(ctx context.Context, filter ListFilter) ([]ListOutput, error) {
	columns := []string{
		"p.id", "p.name", "p.description", "p.photo",
		"p.category_id", "p.unit", "p.price", "p.discount_rules",
	}

	var b sq.SelectBuilder

	if filter.WithCategoryName {
		b = sq.Select(append(columns, "c.name as category_name")...).
			From("catalog.products p").
			Join("catalog.categories c ON c.id = p.category_id").
			Where("p.deleted_at IS NULL").
			PlaceholderFormat(sq.Dollar)
	} else {
		b = sq.Select(append(columns, "'' as category_name")...).
			From("catalog.products p").
			Where("p.deleted_at IS NULL").
			PlaceholderFormat(sq.Dollar)
	}

	if len(filter.IDs) > 0 {
		b = b.Where(sq.Eq{"p.id": filter.IDs})
	}

	if filter.TenantID != "" {
		b = b.Where(sq.Eq{"p.tenant_id": filter.TenantID})
	}

	if filter.Name != "" {
		b = b.Where(sq.ILike{"p.name": "%" + filter.Name + "%"})
	}

	if filter.CategoryID != "" {
		b = b.Where(sq.Eq{"p.category_id": filter.CategoryID})
	}

	if !filter.IsActive.IsEmpty() {
		b = b.Where(sq.Eq{"p.is_active": filter.IsActive.Value()})
	}

	if filter.Limit > 0 {
		b = b.Limit(filter.Limit)
	}

	sql, args, err := b.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := g.DB.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[ListOutput])
}
