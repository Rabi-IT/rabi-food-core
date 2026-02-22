package order_gateway

import (
	"rabi-food-core/libs/database/gorm_adapter/models"
)

func (g *GormOrderGatewayAdapter) Patch(filter PatchFilter, newValues PatchValues) (bool, error) {
	query := g.DB.Conn.Model(&models.Order{}).Where("id = ?", filter.ID)

	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}

	if filter.PaymentStatus != "" {
		query = query.Where("payment_status = ?", filter.PaymentStatus)
	}

	result := query.Updates(newValues)

	return result.RowsAffected > 0, result.Error
}
