package gateway

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

func (g *PgxTenantGatewayAdapter) Patch(filter PatchFilter, values PatchValues) (bool, error) {
	b := sq.
		Update("iam.tenants").
		Where(sq.Eq{"id": filter.ID}).
		PlaceholderFormat(sq.Dollar)

	if values.Name != nil {
		b = b.Set("name", *values.Name)
	}

	sql, args, err := b.ToSql()
	if err != nil {
		return false, err
	}

	tag, err := g.DB.Pool.Exec(context.Background(), sql, args...)

	return tag.RowsAffected() > 0, err
}
