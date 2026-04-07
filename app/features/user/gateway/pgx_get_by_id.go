package gateway

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

var ErrMustFilterByIDOrTenantID = errors.New("must filter by either ID or TenantID")

func (g *PgxUserGatewayAdapter) GetByID(filter GetByIDFilter) (*GetByIDOutput, error) {
	if filter.ID == "" && filter.TenantID == "" {
		return nil, ErrMustFilterByIDOrTenantID
	}

	b := sq.
		Select("id", "tenant_id", "social_id", "street", "complement", "name",
			"email", "photo", "tax_id", "phone", "city", "state", "zip").
		From("users").
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	if filter.ID != "" {
		b = b.Where(sq.Eq{"id": filter.ID})
	}

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
