package gateway

import (
	"rabi-food-core/features/category/model"
	"rabi-food-core/libs/database"
)

func (g *GormCategoryGatewayAdapter) Paginate(
	filter PaginateFilter,
	paginate database.PaginateInput,
) (PaginateOutput, error) {
	query := g.DB.Conn.Model(&model.Category{})

	if filter.Name != nil {
		query = query.Where("name = ?", filter.Name)
	}

	if filter.TenantID != nil {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}

	count := int64(0)
	data := []PaginateData{}
	err := database.Paginate(query, &count, &data, paginate)
	if err != nil {
		return PaginateOutput{}, err
	}

	output := PaginateOutput{
		Data:     data,
		MaxPages: paginate.CalcMaxPages(count),
	}

	return output, nil
}
