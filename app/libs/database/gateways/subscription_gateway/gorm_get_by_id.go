package subscription_gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"rabi-food-core/libs/database/gorm_adapter/models"
)

func (g *GormSubscriptionGatewayAdapter) GetByID(ctx context.Context, filter GetByIDFilter) (*GetByIDOutput, error) {
	output := &models.Subscription{}
	query := g.DB.Conn.WithContext(ctx).Limit(1).Where("id = ?", filter.ID)

	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}

	if filter.UserID != "" {
		query = query.Where("user_id = ?", filter.UserID)
	}

	result := query.Find(output)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	var items []SubscriptionItem
	err := json.Unmarshal(output.Items, &items)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal subscription items: %w", err)
	}

	var deliveryDays []DeliveryDay
	if len(output.DeliveryDays) > 0 {
		err := json.Unmarshal(output.DeliveryDays, &deliveryDays)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal delivery days: %w", err)
		}
	}

	return &GetByIDOutput{
		ID:                  output.ID,
		TenantID:            output.TenantID,
		UserID:              output.UserID,
		Status:              output.Status,
		DeliveryDays:        deliveryDays,
		Items:               items,
		Notes:               output.Notes,
		TotalCycles:         output.TotalCycles,
		RemainingCycles:     output.RemainingCycles,
		CycleDiscount:       output.CycleDiscount,
		CutoffOffsetMinutes: output.CutoffOffsetMinutes,
		AutoRenew:           output.AutoRenew,
		ItemsTotal:          output.ItemsTotal,
		ItemsDiscount:       output.ItemsDiscount,
		PaymentAmount:       output.PaymentAmount,
		CreatedAt:           output.CreatedAt,
	}, nil
}
