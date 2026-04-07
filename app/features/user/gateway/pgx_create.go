package gateway

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (g *PgxUserGatewayAdapter) Create(input CreateInput) (string, error) {
	id := uuid.Must(uuid.NewV7()).String()

	sql, args, err := sq.
		Insert("users").
		Columns("id", "tenant_id", "social_id", "street", "complement", "name",
			"email", "photo", "tax_id", "phone", "city", "state", "zip", "neighborhood", "role").
		Values(id, input.TenantID, input.SocialID, input.Street, input.Complement, input.Name,
			input.Email, input.Photo, input.TaxID, input.Phone, input.City, input.State,
			input.ZIP, input.Neighborhood, input.Role).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return "", err
	}

	_, err = g.DB.Pool.Exec(context.Background(), sql, args...)

	return id, err
}
