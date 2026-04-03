package subscription_gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"rabi-food-core/features/subscription/model"
)

func (g *GormSubscriptionGatewayAdapter) GetConfig(ctx context.Context, tenantID string) (*GetConfigOutput, error) {
	output := &model.SubscriptionConfig{}
	query := g.DB.Conn.WithContext(ctx).Limit(1).Where("tenant_id = ?", tenantID)

	result := query.Find(output)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	discountRules := make([]DiscountRule, 0)
	err := json.Unmarshal(output.DiscountRules, &discountRules)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal discount rules: %w", err)
	}

	return &GetConfigOutput{
		MaxAttemptsPerOrder: output.MaxAttemptsPerOrder,
		DiscountRules:       discountRules,
		CutoffOffsetMinutes: output.CutoffOffsetMinutes,
	}, nil
}
