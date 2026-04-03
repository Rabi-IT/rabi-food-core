package gateway

import "github.com/Rabi-IT/rabi-food-core/features/tenant/model"

func (g *GormTenantGatewayAdapter) Patch(filter PatchFilter, newValues PatchValues) (bool, error) {
	query := g.DB.Conn.Model(&model.Tenant{}).Where("id = ?", filter.ID)

	result := query.Updates(newValues)

	return result.RowsAffected > 0, result.Error
}
