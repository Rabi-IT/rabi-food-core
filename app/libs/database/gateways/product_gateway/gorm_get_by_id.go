package product_gateway

import (
	"rabi-food-core/libs/database/gorm_adapter/models"
)

func (g *GormProductGatewayAdapter) GetByID(filter GetByIDFilter) (*GetByIDOutput, error) {
	output := &models.Product{}
	query := g.DB.Conn.Limit(1).
		Where("id = ?", filter.ID).
		Preload("Category")

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

	adapted := GetByIDOutput{
		ID:           output.ID,
		TenantID:     output.TenantID,
		Name:         output.Name,
		Description:  output.Description,
		Photo:        output.Photo,
		CategoryID:   output.CategoryID,
		CategoryName: output.Category.Name,
		Unit:         output.Unit,
		Price:        output.Price,
		IsActive:     output.IsActive,
	}

	return &adapted, nil
}
