package gateway

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/Rabi-IT/rabi-food-core/libs/database"
)

func (g *PgxProductGatewayAdapter) Paginate(
	ctx context.Context,
	filter PaginateFilter,
	paginate database.PaginateInput,
) (PaginateOutput, error) {
	base := sq.Select().From("catalog.products").
		Where("deleted_at IS NULL").
		PlaceholderFormat(sq.Dollar)

	if filter.Name != "" {
		base = base.Where(sq.ILike{"name": "%" + filter.Name + "%"})
	}

	if filter.CategoryID != "" {
		base = base.Where(sq.Eq{"category_id": filter.CategoryID})
	}

	if !filter.IsActive.IsEmpty() {
		base = base.Where(sq.Eq{"is_active": filter.IsActive.Value()})
	}

	if filter.TenantID != "" {
		base = base.Where(sq.Eq{"tenant_id": filter.TenantID})
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
		ctx, g.DB.Pool,
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
