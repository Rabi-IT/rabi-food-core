package category_gateway

import (
	"rabi-food-core/libs/database/gorm_adapter/models"
)

func (g *GormCategoryGatewayAdapter) Patch(filter PatchFilter, newValues PatchValues) (bool, error) {
	query := g.DB.Conn.Model(&models.Category{}).Where("id = ?", filter.ID)

	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}

	result := query.Updates(newValues)

	return result.RowsAffected > 0, result.Error
}
