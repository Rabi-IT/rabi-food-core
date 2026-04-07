package gateway

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

func (g *PgxCategoryGatewayAdapter) Delete(filter DeleteFilter) (bool, error) {
	b := sq.
		Delete("categories").
		Where(sq.Eq{"id": filter.ID}).
		PlaceholderFormat(sq.Dollar)

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
