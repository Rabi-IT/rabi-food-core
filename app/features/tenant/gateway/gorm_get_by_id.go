package gateway

import "rabi-food-core/features/tenant/model"

func (g *GormTenantGatewayAdapter) GetByID(id string) (*GetByIDOutput, error) {
	output := &model.Tenant{}
	result := g.DB.Conn.Model(&model.Tenant{}).
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
