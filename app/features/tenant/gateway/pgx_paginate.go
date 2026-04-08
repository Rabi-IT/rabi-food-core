package gateway

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/Rabi-IT/rabi-food-core/libs/database"
)

func (g *PgxTenantGatewayAdapter) Paginate(
	filter PaginateFilter,
	paginate database.PaginateInput,
) (PaginateOutput, error) {
	base := sq.Select().From("tenants").PlaceholderFormat(sq.Dollar)

	if filter.Name != nil {
		base = base.Where(sq.ILike{"name": "%" + *filter.Name + "%"})
	}

	countSQL, countArgs, err := base.Column("COUNT(*)").ToSql()
	if err != nil {
		return PaginateOutput{}, err
	}

	dataSQL, dataArgs, err := base.
		Columns("id", "name").
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
