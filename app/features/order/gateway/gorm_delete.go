package gateway

import "github.com/Rabi-IT/rabi-food-core/features/order/model"

func (g *GormOrderGatewayAdapter) Delete(filter DeleteFilter) (bool, error) {
	query := g.DB.Conn.Where("id = ?", filter.ID)

	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}

	result := query.Delete(&model.Order{})

	return result.RowsAffected > 0, result.Error
}
