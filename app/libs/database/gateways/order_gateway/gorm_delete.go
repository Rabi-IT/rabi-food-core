package order_gateway

import (
	"rabi-food-core/libs/database/gorm_adapter/models"
)

func (g *GormOrderGatewayAdapter) Delete(filter DeleteFilter) (bool, error) {
	query := g.DB.Conn.Where("id = ?", filter.ID)

	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}

	result := query.Delete(&models.Order{})

	return result.RowsAffected > 0, result.Error
}
