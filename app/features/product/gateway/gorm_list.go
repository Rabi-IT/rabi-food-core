package gateway

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Rabi-IT/rabi-food-core/features/product/model"
)

var (
	ErrIDsRequired = errors.New("ids is required")
)

func (g *GormProductGatewayAdapter) List(filter ListFilter) ([]ListOutput, error) {
	if len(filter.IDs) == 0 {
		return nil, ErrIDsRequired
	}

	query := g.DB.Conn.Model(&model.Product{}).
		Where("id IN ?", filter.IDs)

	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", filter.IsActive)
	}

	output := &[]model.Product{}
	result := query.Find(&output)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return []ListOutput{}, nil
	}

	adapted := make([]ListOutput, len(*output))
	for i, product := range *output {
		discountRules := make([]DiscountRule, 0)
		err := json.Unmarshal(product.DiscountRules, &discountRules)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal discount rules for product %s: %w", product.ID, err)
		}

		adapted[i] = ListOutput{
			ID:            product.ID,
			Name:          product.Name,
			Price:         product.Price,
			DiscountRules: discountRules,
		}
	}

	return adapted, nil
}
