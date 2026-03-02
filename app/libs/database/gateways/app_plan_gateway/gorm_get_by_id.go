package app_plan_gateway

import (
	"rabi-food-core/libs/database/gorm_adapter/models"
)

func (g *GormAppPlanGatewayAdapter) GetByID(filter GetByIDFilter) (*GetByIDOutput, error) {
	output := &models.AppPlan{}
	query := g.DB.Conn.Limit(1).Where("id = ?", filter.ID)

	result := query.Find(output)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	adapted := GetByIDOutput{
		ID:          output.ID,
		Name:        output.Name,
		Description: output.Description,
	}

	return &adapted, nil
}
