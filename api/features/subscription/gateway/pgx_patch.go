package subscription_gateway

import (
	"context"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
)

func (g *PgxSubscriptionGatewayAdapter) Patch(ctx context.Context, filter PatchFilter, values PatchValues) (bool, error) {
	if filter.ID == "" && filter.TenantID == "" && len(filter.Statuses) == 0 && filter.PaymentStatus == "" {
		return false, errors.New("patch requires at least one filter")
	}

	b := sq.
		Update("subscription.subscriptions").
		Set("updated_at", time.Now().UTC()).
		PlaceholderFormat(sq.Dollar)

	if values.Status != "" {
		b = b.Set("status", values.Status)
	}
	if values.PaymentStatus != nil {
		b = b.Set("payment_status", *values.PaymentStatus)
	}
	if values.PaidAt != nil {
		b = b.Set("paid_at", *values.PaidAt)
	}

	if filter.ID != "" {
		b = b.Where(sq.Eq{"id": filter.ID})
	}
	if filter.TenantID != "" {
		b = b.Where(sq.Eq{"tenant_id": filter.TenantID})
	}
	if len(filter.Statuses) > 0 {
		b = b.Where(sq.Eq{"status": filter.Statuses})
	}
	if filter.PaymentStatus != "" {
		b = b.Where(sq.Eq{"payment_status": filter.PaymentStatus})
	}

	sql, args, err := b.ToSql()
	if err != nil {
		return false, err
	}

	tag, err := g.DB.Pool.Exec(context.WithoutCancel(ctx), sql, args...)

	return tag.RowsAffected() > 0, err
}
