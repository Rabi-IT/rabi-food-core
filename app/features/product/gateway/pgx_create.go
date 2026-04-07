package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (g *PgxProductGatewayAdapter) Create(input CreateInput) (string, error) {
	id := uuid.Must(uuid.NewV7()).String()

	discountRules, err := json.Marshal(input.DiscountRules)
	if err != nil {
		return "", fmt.Errorf("failed to marshal discount rules: %w", err)
	}

	now := time.Now().UTC()

	sql, args, err := sq.
		Insert("products").
		Columns("id", "tenant_id", "name", "description", "photo", "category_id",
			"discount_rules", "unit", "price", "is_active", "created_at", "updated_at").
		Values(id, input.TenantID, input.Name, input.Description, input.Photo, input.CategoryID,
			discountRules, input.Unit, input.Price, input.IsActive, now, now).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return "", err
	}

	_, err = g.DB.Pool.Exec(context.Background(), sql, args...)

	return id, err
}
