package gateway

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/Rabi-IT/rabi-food-core/libs/database"
)

func (g *PgxProductGatewayAdapter) Paginate(
	filter PaginateFilter,
	paginate database.PaginateInput,
) (PaginateOutput, error) {
	base := sq.Select().From("catalog.products").
		Where("deleted_at IS NULL").
		PlaceholderFormat(sq.Dollar)

	if filter.Name != nil {
		base = base.Where(sq.ILike{"name": "%" + *filter.Name + "%"})
	}

	if filter.CategoryID != nil {
		base = base.Where(sq.Eq{"category_id": *filter.CategoryID})
	}

	if filter.IsActive != nil {
		base = base.Where(sq.Eq{"is_active": *filter.IsActive})
	}

	if filter.TenantID != nil {
		base = base.Where(sq.Eq{"tenant_id": *filter.TenantID})
	}

	countSQL, countArgs, err := base.Column("COUNT(*)").ToSql()
	if err != nil {
		return PaginateOutput{}, err
	}

	dataSQL, dataArgs, err := base.
		Columns("id", "name", "description", "photo", "category_id", "unit", "price", "is_active").
		Limit(uint64(paginate.PageSize)).
		Offset(uint64(paginate.Offset())).
		ToSql()
	if err != nil {
		return PaginateOutput{}, err
	}

	data, count, err := database.Paginate[PaginateData](
		context.Background(), g.DB.Pool,
		countSQL, countArgs,
		dataSQL, dataArgs,
	)
	if err != nil {
		return PaginateOutput{}, err
	}

	return PaginateOutput{
		Data:     data,
		MaxPages: paginate.CalcMaxPages(count),
	}, nil
}
