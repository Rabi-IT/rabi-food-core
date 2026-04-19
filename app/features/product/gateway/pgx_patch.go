package gateway

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
)

func (g *PgxProductGatewayAdapter) Patch(ctx context.Context, filter PatchFilter, values PatchValues) (bool, error) {
	b := sq.
		Update("catalog.products").
		Set("updated_at", time.Now().UTC()).
		Where("deleted_at IS NULL").
		Where(sq.Eq{"id": filter.ID}).
		PlaceholderFormat(sq.Dollar)

	if values.Name != nil {
		b = b.Set("name", *values.Name)
	}
	if values.Description != nil {
		b = b.Set("description", *values.Description)
	}
	if values.Photo != nil {
		b = b.Set("photo", *values.Photo)
	}
	if values.CategoryID != nil {
		b = b.Set("category_id", *values.CategoryID)
	}
	if values.Unit != nil {
		b = b.Set("unit", *values.Unit)
	}
	if values.Price != nil {
		b = b.Set("price", *values.Price)
	}
	if values.IsActive != nil {
		b = b.Set("is_active", *values.IsActive)
	}

	if filter.TenantID != "" {
		b = b.Where(sq.Eq{"tenant_id": filter.TenantID})
	}

	sql, args, err := b.ToSql()
	if err != nil {
		return false, err
	}

	tag, err := g.DB.Pool.Exec(context.WithoutCancel(ctx), sql, args...)

	return tag.RowsAffected() > 0, err
}
