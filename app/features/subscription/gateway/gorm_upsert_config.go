package subscription_gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"rabi-food-core/features/subscription/model"

	"gorm.io/gorm/clause"
)

func (g *GormSubscriptionGatewayAdapter) UpsertConfig(
	ctx context.Context,
	tenantID string,
	input UpsertConfigInput,
) error {
	discountRules, err := json.Marshal(input.DiscountRules)
	if err != nil {
		return fmt.Errorf("failed to marshal discount rules: %w", err)
	}

	result := g.DB.Conn.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).
		Create(&model.SubscriptionConfig{
			TenantID:            tenantID,
			MaxAttemptsPerOrder: input.MaxAttemptsPerOrder,
			DiscountRules:       discountRules,
			CutoffOffsetMinutes: input.CutoffOffsetMinutes,
			IsOpen:              input.IsOpen,
		})

	return result.Error
}
