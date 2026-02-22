package product_gateway

import (
	"errors"
	"rabi-food-core/libs/database/gorm_adapter/models"
)

var (
	ErrIDsRequired = errors.New("ids is required")
)

func (g *GormProductGatewayAdapter) List(filter ListFilter) ([]ListOutput, error) {
	if len(filter.IDs) == 0 {
		return nil, ErrIDsRequired
	}

	var output []ListOutput
	query := g.DB.Conn.Model(&models.Product{}).
		Where("id IN ?", filter.IDs)

	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", filter.IsActive)
	}

	result := query.Find(&output)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return output, nil
}
