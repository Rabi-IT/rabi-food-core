package subscription_gateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type subscriptionConfigRow struct {
	MaxAttemptsPerOrder uint8  `db:"max_attempts_per_order"`
	DiscountRulesRaw    []byte `db:"discount_rules"`
	CutoffOffsetMinutes uint16 `db:"cutoff_offset_minutes"`
	IsOpen              bool   `db:"is_open"`
}

func (g *PgxSubscriptionGatewayAdapter) GetConfig(ctx context.Context, tenantID string) (*GetConfigOutput, error) {
	sql, args, err := sq.
		Select("max_attempts_per_order", "discount_rules", "cutoff_offset_minutes", "is_open").
		From("subscription_configs").
		Where(sq.Eq{"tenant_id": tenantID}).
		Limit(1).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := g.DB.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[subscriptionConfigRow])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	discountRules := make([]DiscountRule, 0)
	if len(row.DiscountRulesRaw) > 0 {
		if err := json.Unmarshal(row.DiscountRulesRaw, &discountRules); err != nil {
			return nil, fmt.Errorf("failed to unmarshal discount rules: %w", err)
		}
	}

	return &GetConfigOutput{
		MaxAttemptsPerOrder: row.MaxAttemptsPerOrder,
		DiscountRules:       discountRules,
		CutoffOffsetMinutes: row.CutoffOffsetMinutes,
		IsOpen:              row.IsOpen,
	}, nil
}
