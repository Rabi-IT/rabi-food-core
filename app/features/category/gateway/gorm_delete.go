package gateway

import "rabi-food-core/features/category/model"

func (g *GormCategoryGatewayAdapter) Delete(filter DeleteFilter) (bool, error) {
	query := g.DB.Conn.Where("id = ?", filter.ID)

	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}

	result := query.Delete(&model.Category{})

	return result.RowsAffected > 0, result.Error
}
