package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (g *PgxOrderGatewayAdapter) Create(input CreateInput) (string, error) {
	id := uuid.Must(uuid.NewV7()).String()

	items, err := json.Marshal(input.Items)
	if err != nil {
		return "", fmt.Errorf("failed to marshal order items: %w", err)
	}

	now := time.Now().UTC()

	sql, args, err := sq.
		Insert("commerce.orders").
		Columns("id", "tenant_id", "user_id", "code", "fulfillment_status",
			"delivery_status", "payment_status", "notes", "total_price", "items",
			"created_at", "updated_at").
		Values(id, input.TenantID, input.UserID, input.Code, input.FulfillmentStatus,
			input.DeliveryStatus, input.PaymentStatus, input.Notes, input.TotalPrice, items,
			now, now).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return "", err
	}

	_, err = g.DB.Pool.Exec(context.Background(), sql, args...)

	return id, err
}
