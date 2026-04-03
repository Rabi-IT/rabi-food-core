package gateway

import "github.com/Rabi-IT/rabi-food-core/features/product/model"

func (g *GormProductGatewayAdapter) Delete(filter DeleteFilter) (bool, error) {
	query := g.DB.Conn.Where(
		"id = ?", filter.ID,
	)

	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}

	result := query.Delete(&model.Product{})

	return result.RowsAffected > 0, result.Error
}
