package subscription_gateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/Rabi-IT/rabi-food-core/features/subscription"
	"github.com/jackc/pgx/v5"
)

type subscriptionRow struct {
	ID                  string              `db:"id"`
	TenantID            string              `db:"tenant_id"`
	UserID              string              `db:"user_id"`
	Status              subscription.Status `db:"status"`
	DeliveryDaysRaw     []byte              `db:"delivery_days"`
	ItemsRaw            []byte              `db:"items"`
	Notes               string              `db:"notes"`
	TotalCycles         uint                `db:"total_cycles"`
	RemainingCycles     uint                `db:"remaining_cycles"`
	CycleDiscount       uint                `db:"cycle_discount"`
	CutoffOffsetMinutes uint16              `db:"cutoff_offset_minutes"`
	AutoRenew           bool                `db:"auto_renew"`
	ItemsTotal          uint                `db:"items_total"`
	ItemsDiscount       uint                `db:"items_discount"`
	PaymentAmount       uint                `db:"payment_amount"`
	CreatedAt           time.Time           `db:"created_at"`
}

func (g *PgxSubscriptionGatewayAdapter) GetByID(ctx context.Context, filter GetByIDFilter) (*GetByIDOutput, error) {
	b := sq.
		Select("id", "tenant_id", "user_id", "status", "delivery_days", "items", "notes",
			"total_cycles", "remaining_cycles", "cycle_discount", "cutoff_offset_minutes",
			"auto_renew", "items_total", "items_discount", "payment_amount", "created_at").
		From("subscriptions").
		Where(sq.Eq{"id": filter.ID}).
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	if filter.TenantID != "" {
		b = b.Where(sq.Eq{"tenant_id": filter.TenantID})
	}

	if filter.UserID != "" {
		b = b.Where(sq.Eq{"user_id": filter.UserID})
	}

	sql, args, err := b.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := g.DB.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[subscriptionRow])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var items []SubscriptionItem
	if err := json.Unmarshal(row.ItemsRaw, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal subscription items: %w", err)
	}

	var deliveryDays []DeliveryDay
	if len(row.DeliveryDaysRaw) > 0 {
		if err := json.Unmarshal(row.DeliveryDaysRaw, &deliveryDays); err != nil {
			return nil, fmt.Errorf("failed to unmarshal delivery days: %w", err)
		}
	}

	return &GetByIDOutput{
		ID:                  row.ID,
		TenantID:            row.TenantID,
		UserID:              row.UserID,
		Status:              row.Status,
		DeliveryDays:        deliveryDays,
		Items:               items,
		Notes:               row.Notes,
		TotalCycles:         row.TotalCycles,
		RemainingCycles:     row.RemainingCycles,
		CycleDiscount:       row.CycleDiscount,
		CutoffOffsetMinutes: row.CutoffOffsetMinutes,
		AutoRenew:           row.AutoRenew,
		ItemsTotal:          row.ItemsTotal,
		ItemsDiscount:       row.ItemsDiscount,
		PaymentAmount:       row.PaymentAmount,
		CreatedAt:           row.CreatedAt,
	}, nil
}
