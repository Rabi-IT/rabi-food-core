package gateway

import (
	"encoding/json"
	"fmt"

	"github.com/Rabi-IT/rabi-food-core/features/order/model"
)

func (g *GormOrderGatewayAdapter) GetByID(filter GetByIDFilter) (*GetByIDOutput, error) {
	output := &model.Order{}
	query := g.DB.Conn.Limit(1).
		Preload("User").
		Where("id = ?", filter.ID)

	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}

	result := query.Find(output)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	var items []OrderItem
	err := json.Unmarshal(output.Items, &items)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal order items: %w", err)
	}

	adapted := GetByIDOutput{
		ID:                output.ID,
		TenantID:          output.TenantID,
		Code:              output.Code,
		PaymentStatus:     output.PaymentStatus,
		FulfillmentStatus: output.FulfillmentStatus,
		DeliveryStatus:    output.DeliveryStatus,
		Notes:             output.Notes,
		TotalPrice:        output.TotalPrice,
		Items:             items,
		CreatedAt:         output.CreatedAt,
		PaidAt:            output.PaidAt,
		ExternalPaymentID: output.ExternalPaymentID,
	}

	return &adapted, nil
}
