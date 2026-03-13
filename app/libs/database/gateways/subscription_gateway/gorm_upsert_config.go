package subscription_gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"rabi-food-core/libs/database/gorm_adapter/models"

	"gorm.io/gorm/clause"
)

func (g *GormSubscriptionGatewayAdapter) UpsertConfig(
	ctx context.Context,
	tenatID string,
	input UpsertConfigInput,
) (bool, error) {
	discountRules, err := json.Marshal(input.DiscountRules)
	if err != nil {
		return false, fmt.Errorf("failed to marshal discount rules: %w", err)
	}

	result := g.DB.Conn.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).
		Create(&models.SubscriptionConfig{
			TenantID:            tenatID,
			MaxAttemptsPerOrder: input.MaxAttemptsPerOrder,
			DiscountRules:       discountRules,
			CutoffOffsetMinutes: input.CutoffOffsetMinutes,
		})

	return result.RowsAffected > 0, result.Error
}
