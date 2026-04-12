package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/Rabi-IT/rabi-food-core/domain/payment_status"
	"github.com/Rabi-IT/rabi-food-core/features/order"
	"github.com/jackc/pgx/v5"
)

type orderGetByIDRow struct {
	ID                string                  `db:"id"`
	TenantID          string                  `db:"tenant_id"`
	Code              string                  `db:"code"`
	PaymentStatus     payment_status.Status   `db:"payment_status"`
	FulfillmentStatus order.FulfillmentStatus `db:"fulfillment_status"`
	DeliveryStatus    order.DeliveryStatus    `db:"delivery_status"`
	Notes             string                  `db:"notes"`
	TotalPrice        uint                    `db:"total_price"`
	ItemsRaw          []byte                  `db:"items"`
	PaidAt            *time.Time              `db:"paid_at"`
	CreatedAt         time.Time               `db:"created_at"`
	ExternalPaymentID *string                 `db:"external_payment_id"`
}

func (g *PgxOrderGatewayAdapter) GetByID(filter GetByIDFilter) (*GetByIDOutput, error) {
	b := sq.
		Select("id", "tenant_id", "code", "payment_status", "fulfillment_status",
			"delivery_status", "notes", "total_price", "items", "paid_at",
			"created_at", "external_payment_id").
		From("commerce.orders").
		Where("deleted_at IS NULL").
		Where(sq.Eq{"id": filter.ID}).
		Limit(1).
		PlaceholderFormat(sq.Dollar)

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

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[orderGetByIDRow])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var items []OrderItem
	if err := json.Unmarshal(row.ItemsRaw, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order items: %w", err)
	}

	return &GetByIDOutput{
		ID:                row.ID,
		TenantID:          row.TenantID,
		Code:              row.Code,
		PaymentStatus:     row.PaymentStatus,
		FulfillmentStatus: row.FulfillmentStatus,
		DeliveryStatus:    row.DeliveryStatus,
		Notes:             row.Notes,
		TotalPrice:        row.TotalPrice,
		Items:             items,
		PaidAt:            row.PaidAt,
		CreatedAt:         row.CreatedAt,
		ExternalPaymentID: row.ExternalPaymentID,
	}, nil
}
