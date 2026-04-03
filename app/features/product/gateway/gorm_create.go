package gateway

import (
	"encoding/json"
	"fmt"
	"rabi-food-core/features/product/model"

	"github.com/google/uuid"
)

func (g *GormProductGatewayAdapter) Create(input CreateInput) (string, error) {
	id := uuid.Must(uuid.NewV7()).String()

	discountRules, err := json.Marshal(input.DiscountRules)
	if err != nil {
		return "", fmt.Errorf("failed to marshal discount rules: %w", err)
	}

	result := g.DB.Conn.Create(&model.Product{
		ID:            id,
		TenantID:      input.TenantID,
		Name:          input.Name,
		Description:   input.Description,
		Photo:         input.Photo,
		CategoryID:    input.CategoryID,
		Unit:          input.Unit,
		Price:         input.Price,
		IsActive:      input.IsActive,
		DiscountRules: discountRules,
	})

	return id, result.Error
}
