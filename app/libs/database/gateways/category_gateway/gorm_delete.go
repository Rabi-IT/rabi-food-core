package category_gateway

import (
	"rabi-food-core/libs/database/gorm_adapter/models"
)

func (g *GormCategoryGatewayAdapter) Delete(filter DeleteFilter) (bool, error) {
	query := g.DB.Conn.Where("id = ?", filter.ID)

	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}

	result := query.Delete(&models.Category{})

	return result.RowsAffected > 0, result.Error
}
