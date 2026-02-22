package user_gateway

import (
	"rabi-food-core/libs/database/gorm_adapter/models"
)

func (g *GormUserGatewayAdapter) Delete(filter DeleteFilter) (bool, error) {
	result := g.DB.Conn.Where(
		"id = ?", filter.ID,
	)

	if filter.TenantID != "" {
		result = result.Where("tenant_id = ?", filter.TenantID)
	}

	result = result.Delete(&models.User{})

	return result.RowsAffected > 0, result.Error
}
