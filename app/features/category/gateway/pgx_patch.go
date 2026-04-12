package gateway

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

func (g *PgxCategoryGatewayAdapter) Patch(filter PatchFilter, values PatchValues) (bool, error) {
	b := sq.
		Update("catalog.categories").
		Where(sq.Eq{"id": filter.ID}).
		PlaceholderFormat(sq.Dollar)

	if values.Name != nil {
		b = b.Set("name", *values.Name)
	}
	if values.Description != nil {
		b = b.Set("description", *values.Description)
	}

	if filter.TenantID != "" {
		b = b.Where(sq.Eq{"tenant_id": filter.TenantID})
	}

	sql, args, err := b.ToSql()
	if err != nil {
		return false, err
	}

	tag, err := g.DB.Pool.Exec(context.Background(), sql, args...)

	return tag.RowsAffected() > 0, err
}
