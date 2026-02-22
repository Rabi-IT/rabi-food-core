package product_gateway

import (
	"rabi-food-core/libs/database"
	"rabi-food-core/libs/database/gorm_adapter"
	"rabi-food-core/libs/database/gorm_adapter/models"
)

func (g *GormProductGatewayAdapter) Paginate(
	filter PaginateFilter,
	paginate database.PaginateInput,
) (PaginateOutput, error) {
	query := g.DB.Conn.Model(&models.Product{})

	if filter.Name != nil {
		query = query.Where("name = ?", filter.Name)
	}

	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", filter.CategoryID)
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", filter.IsActive)
	}

	if filter.TenantID != nil {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}

	count := int64(0)
	data := []PaginateData{}
	err := gorm_adapter.Paginate(query, &count, &data, paginate)
	if err != nil {
		return PaginateOutput{}, err
	}

	output := PaginateOutput{
		Data:     data,
		MaxPages: paginate.CalcMaxPages(count),
	}

	return output, nil
}
