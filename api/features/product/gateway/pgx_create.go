package gateway

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (g *PgxProductGatewayAdapter) Create(ctx context.Context, input CreateInput) (string, error) {
	id := uuid.Must(uuid.NewV7()).String()
	now := time.Now().UTC()

	sql, args, err := sq.
		Insert("catalog.products").
		Columns("id", "tenant_id", "name", "description", "photo", "category_id",
			"discount_rules", "unit", "price", "is_active", "created_at", "updated_at").
		Values(id, input.TenantID, input.Name, input.Description, input.Photo, input.CategoryID,
			input.DiscountRules, input.Unit, input.Price, input.IsActive, now, now).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return "", err
	}

	_, err = g.DB.Pool.Exec(context.WithoutCancel(ctx), sql, args...)

	return id, err
}
