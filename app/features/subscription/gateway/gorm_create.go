package subscription_gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"rabi-food-core/features/subscription/model"

	"github.com/google/uuid"
)

func (g *GormSubscriptionGatewayAdapter) Create(ctx context.Context, input CreateInput) (string, error) {
	id := uuid.Must(uuid.NewV7()).String()

	items, err := json.Marshal(input.Items)
	if err != nil {
		return "", fmt.Errorf("failed to marshal subscription items: %w", err)
	}

	deliveryDays, err := json.Marshal(input.DeliveryDays)
	if err != nil {
		return "", fmt.Errorf("failed to marshal delivery days: %w", err)
	}

	weekdaysMask := uint8(0)
	for _, day := range input.DeliveryDays {
		weekdaysMask |= weekdayBit(day.Weekday)
	}

	result := g.DB.Conn.WithContext(ctx).Create(&model.Subscription{
		ID:                   id,
		RootSubscriptionID:   input.RootID,
		TenantID:             input.TenantID,
		UserID:               input.UserID,
		Status:               input.Status,
		Items:                items,
		DeliveryDays:         deliveryDays,
		DeliveryWeekdaysMask: weekdaysMask,
		Notes:                input.Notes,
		TotalCycles:          input.TotalCycles,
		RemainingCycles:      input.RemainingCycles,
		CycleDiscount:        input.CycleDiscount,
		CutoffOffsetMinutes:  input.CutoffOffsetMinutes,
		AutoRenew:            input.AutoRenew,
		MaxAttemptsPerOrder:  input.MaxAttemptsPerOrder,
		ItemsTotal:           input.ItemsTotal,
		ItemsDiscount:        input.ItemsDiscount,
		PaymentAmount:        input.PaymentAmount,
		PaymentStatus:        input.PaymentStatus,
		ExternalPaymentID:    input.ExternalPaymentID,
	})

	return id, result.Error
}
