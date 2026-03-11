package subscription_gateway

import (
	"context"
	"fmt"
	"rabi-food-core/libs/database/gorm_adapter/models"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

// CreateDeliveries inserts multiple subscription deliveries in a single batch.
// The operation is idempotent: duplicate (subscription_id, scheduled_date) pairs are
// silently ignored thanks to the unique constraint on the table.
func (g *GormSubscriptionGatewayAdapter) CreateDeliveries(ctx context.Context, inputs []CreateDeliveryInput) error {
	if len(inputs) == 0 {
		return nil
	}

	rows := make([]models.SubscriptionDelivery, 0, len(inputs))
	for _, in := range inputs {
		rows = append(rows, models.SubscriptionDelivery{
			ID:                  uuid.Must(uuid.NewV7()).String(),
			SubscriptionID:      in.SubscriptionID,
			ScheduledDate:       in.ScheduledDate,
			StartHour:           in.StartHour,
			EndHour:             in.EndHour,
			CutoffAt:            in.CutoffAt,
			Status:              in.Status,
			MaxDeliveryAttempts: in.MaxAttempts,
		})
	}

	result := g.DB.Conn.WithContext(ctx).Clauses(clause.OnConflict{
		OnConstraint: models.UniqueIndex__SubscriptionDelivery,
		DoNothing:    true,
	}).Create(&rows)

	if result.Error != nil {
		return fmt.Errorf("failed to create subscription deliveries: %w", result.Error)
	}

	return nil
}
