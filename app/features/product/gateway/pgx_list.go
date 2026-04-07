package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

var ErrIDsRequired = errors.New("ids is required")

type productListRow struct {
	ID            string `db:"id"`
	Name          string `db:"name"`
	Price         uint   `db:"price"`
	DiscountRules []byte `db:"discount_rules"`
}

func (g *PgxProductGatewayAdapter) List(filter ListFilter) ([]ListOutput, error) {
	if len(filter.IDs) == 0 {
		return nil, ErrIDsRequired
	}

	b := sq.
		Select("id", "name", "price", "discount_rules").
		From("products").
		Where("deleted_at IS NULL").
		Where(sq.Eq{"id": filter.IDs}).
		PlaceholderFormat(sq.Dollar)

	if filter.TenantID != "" {
		b = b.Where(sq.Eq{"tenant_id": filter.TenantID})
	}

	if filter.IsActive != nil {
		b = b.Where(sq.Eq{"is_active": *filter.IsActive})
	}

	sql, args, err := b.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := g.DB.Pool.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}

	scanned, err := pgx.CollectRows(rows, pgx.RowToStructByName[productListRow])
	if err != nil {
		return nil, err
	}

	output := make([]ListOutput, len(scanned))

	for i, row := range scanned {
		var discountRules []DiscountRule
		if err := json.Unmarshal(row.DiscountRules, &discountRules); err != nil {
			return nil, fmt.Errorf("failed to unmarshal discount rules for product %s: %w", row.ID, err)
		}

		output[i] = ListOutput{
			ID:            row.ID,
			Name:          row.Name,
			Price:         row.Price,
			DiscountRules: discountRules,
		}
	}

	return output, nil
}
