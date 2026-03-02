package app_plan_gateway

import (
	"rabi-food-core/libs/database/gorm_adapter/models"
)

func (g *GormAppPlanGatewayAdapter) Delete(filter DeleteFilter) (bool, error) {
	query := g.DB.Conn.Where("id = ?", filter.ID)

	result := query.Delete(&models.AppPlan{})

	return result.RowsAffected > 0, result.Error
}
