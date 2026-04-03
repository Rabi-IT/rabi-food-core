package gateway

import "github.com/Rabi-IT/rabi-food-core/features/user/model"

func (g *GormUserGatewayAdapter) Delete(filter DeleteFilter) (bool, error) {
	result := g.DB.Conn.Where(
		"id = ?", filter.ID,
	)

	if filter.TenantID != "" {
		result = result.Where("tenant_id = ?", filter.TenantID)
	}

	result = result.Delete(&model.User{})

	return result.RowsAffected > 0, result.Error
}
