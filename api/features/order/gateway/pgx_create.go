package gateway

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (g *PgxOrderGatewayAdapter) Create(ctx context.Context, input CreateInput) (string, error) {
	id := uuid.Must(uuid.NewV7()).String()
	now := time.Now().UTC()

	sql, args, err := sq.
		Insert("commerce.orders").
		Columns("id", "tenant_id", "user_id", "code", "fulfillment_status",
			"delivery_status", "payment_status", "notes", "total_price", "items",
			"created_at", "updated_at").
		Values(id, input.TenantID, input.UserID, input.Code, input.FulfillmentStatus,
			input.DeliveryStatus, input.PaymentStatus, input.Notes, input.TotalPrice, input.Items,
			now, now).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return "", err
	}

	_, err = g.DB.Pool.Exec(context.WithoutCancel(ctx), sql, args...)

	return id, err
}
