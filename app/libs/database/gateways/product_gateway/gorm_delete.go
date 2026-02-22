package product_gateway

import (
	"rabi-food-core/libs/database/gorm_adapter/models"
)

func (g *GormProductGatewayAdapter) Delete(filter DeleteFilter) (bool, error) {
	query := g.DB.Conn.Where(
		"id = ?", filter.ID,
	)

	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}

	result := query.Delete(&models.Product{})

	return result.RowsAffected > 0, result.Error
}
