package gateway

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

func (g *PgxUserGatewayAdapter) Patch(filter PatchFilter, values PatchValues) (bool, error) {
	b := sq.
		Update("iam.users").
		Where(sq.Eq{"id": filter.ID}).
		PlaceholderFormat(sq.Dollar)

	if values.ZIP != nil {
		b = b.Set("zip", *values.ZIP)
	}
	if values.Phone != nil {
		b = b.Set("phone", *values.Phone)
	}
	if values.City != nil {
		b = b.Set("city", *values.City)
	}
	if values.State != nil {
		b = b.Set("state", *values.State)
	}
	if values.TaxID != nil {
		b = b.Set("tax_id", *values.TaxID)
	}
	if values.SocialID != nil {
		b = b.Set("social_id", *values.SocialID)
	}
	if values.Street != nil {
		b = b.Set("street", *values.Street)
	}
	if values.Complement != nil {
		b = b.Set("complement", *values.Complement)
	}
	if values.Name != nil {
		b = b.Set("name", *values.Name)
	}
	if values.Email != nil {
		b = b.Set("email", *values.Email)
	}
	if values.Photo != nil {
		b = b.Set("photo", *values.Photo)
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
