package gateway

import "rabi-food-core/features/product/model"

func (g *GormProductGatewayAdapter) Patch(filter PatchFilter, newValues PatchValues) (bool, error) {
	query := g.DB.Conn.Model(&model.Product{}).Where("id = ?", filter.ID)

	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}

	result := query.Updates(newValues)

	return result.RowsAffected > 0, result.Error
}
