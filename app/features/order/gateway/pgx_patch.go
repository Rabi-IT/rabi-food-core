package gateway

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
)

func (g *PgxOrderGatewayAdapter) Patch(filter PatchFilter, values PatchValues) (bool, error) {
	b := sq.
		Update("commerce.orders").
		Set("updated_at", time.Now().UTC()).
		Where("deleted_at IS NULL").
		Where(sq.Eq{"id": filter.ID}).
		PlaceholderFormat(sq.Dollar)

	if values.PaymentStatus != nil {
		b = b.Set("payment_status", *values.PaymentStatus)
	}
	if values.ExternalPaymentID != nil {
		b = b.Set("external_payment_id", *values.ExternalPaymentID)
	}
	if values.FulfillmentStatus != nil {
		b = b.Set("fulfillment_status", *values.FulfillmentStatus)
	}
	if values.DeliveryStatus != nil {
		b = b.Set("delivery_status", *values.DeliveryStatus)
	}
	if values.PaidAt != nil {
		b = b.Set("paid_at", *values.PaidAt)
	}

	if filter.TenantID != "" {
		b = b.Where(sq.Eq{"tenant_id": filter.TenantID})
	}

	if filter.PaymentStatus != "" {
		b = b.Where(sq.Eq{"payment_status": filter.PaymentStatus})
	}

	sql, args, err := b.ToSql()
	if err != nil {
		return false, err
	}

	tag, err := g.DB.Pool.Exec(context.Background(), sql, args...)

	return tag.RowsAffected() > 0, err
}
