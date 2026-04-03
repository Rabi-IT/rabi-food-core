package gateway

import "github.com/Rabi-IT/rabi-food-core/features/category/model"

func (g *GormCategoryGatewayAdapter) GetByID(filter GetByIDFilter) (*GetByIDOutput, error) {
	output := &model.Category{}
	query := g.DB.Conn.Limit(1).Where("id = ?", filter.ID)

	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}

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
		TenantID:    output.TenantID,
	}

	return &adapted, nil
}
