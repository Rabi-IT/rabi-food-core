package tenant_gateway

import (
	"rabi-food-core/libs/database/gorm_adapter/models"
)

func (g *GormTenantGatewayAdapter) GetByID(id string) (*GetByIDOutput, error) {
	output := &models.Tenant{}
	result := g.DB.Conn.Model(&models.Tenant{}).
		Limit(1).
		Find(output, "id = ?", id)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	adapted := GetByIDOutput{
		ID:   output.ID,
		Name: output.Name,
	}

	return &adapted, nil
}
