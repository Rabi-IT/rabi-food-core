package order_gateway

import (
	"encoding/json"
	"fmt"
	"rabi-food-core/libs/database/gorm_adapter/models"

	"github.com/google/uuid"
)

func (g *GormOrderGatewayAdapter) Create(input CreateInput) (string, error) {
	id := uuid.NewString()

	items, err := json.Marshal(input.Items)
	if err != nil {
		return "", fmt.Errorf("failed to marshal order items: %w", err)
	}

	result := g.DB.Conn.Create(&models.Order{
		ID:                id,
		TenantID:          input.TenantID,
		UserID:            input.UserID,
		Code:              input.Code,
		FulfillmentStatus: input.FulfillmentStatus,
		DeliveryStatus:    input.DeliveryStatus,
		PaymentStatus:     input.PaymentStatus,
		Notes:             input.Notes,
		TotalPrice:        input.TotalPrice,
		Items:             items,
	})

	return id, result.Error
}
