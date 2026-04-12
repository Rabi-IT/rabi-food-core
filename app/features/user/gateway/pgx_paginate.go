package gateway

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/Rabi-IT/rabi-food-core/libs/database"
)

func (g *PgxUserGatewayAdapter) Paginate(
	filter PaginateFilter,
	paginate database.PaginateInput,
) (PaginateOutput, error) {
	base := sq.Select().From("iam.users").PlaceholderFormat(sq.Dollar)

	if filter.Name != nil {
		base = base.Where(sq.ILike{"name": "%" + *filter.Name + "%"})
	}

	if filter.City != nil {
		base = base.Where(sq.Eq{"city": *filter.City})
	}

	if filter.State != nil {
		base = base.Where(sq.Eq{"state": *filter.State})
	}

	if filter.TenantID != nil {
		base = base.Where(sq.Eq{"tenant_id": *filter.TenantID})
	}

	countSQL, countArgs, err := base.Column("COUNT(*)").ToSql()
	if err != nil {
		return PaginateOutput{}, err
	}

	dataSQL, dataArgs, err := base.
		Columns("id", "photo", "name", "state", "city").
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
