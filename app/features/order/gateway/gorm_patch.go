package gateway

import "github.com/Rabi-IT/rabi-food-core/features/order/model"

func (g *GormOrderGatewayAdapter) Patch(filter PatchFilter, newValues PatchValues) (bool, error) {
	query := g.DB.Conn.Model(&model.Order{}).Where("id = ?", filter.ID)

	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}

	if filter.PaymentStatus != "" {
		query = query.Where("payment_status = ?", filter.PaymentStatus)
	}

	result := query.Updates(newValues)

	return result.RowsAffected > 0, result.Error
}
