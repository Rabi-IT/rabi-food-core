package gateway

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (g *PgxCategoryGatewayAdapter) Create(input CreateInput) (string, error) {
	id := uuid.Must(uuid.NewV7()).String()

	sql, args, err := sq.
		Insert("categories").
		Columns("id", "tenant_id", "name", "description").
		Values(id, input.TenantID, input.Name, input.Description).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return "", err
	}

	_, err = g.DB.Pool.Exec(context.Background(), sql, args...)

	return id, err
}
