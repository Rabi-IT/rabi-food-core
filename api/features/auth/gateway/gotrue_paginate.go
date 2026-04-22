package gateway

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Rabi-IT/rabi-food-core/libs/database"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"
)

type authUserRow struct {
	ID           string `db:"id"`
	Email        string `db:"email"`
	UserMetadata []byte `db:"raw_user_meta_data"`
	AppMetadata  []byte `db:"raw_app_meta_data"`
}

func (r authUserRow) toUserOutput() UserOutput {
	var userMeta, appMeta map[string]any
	_ = json.Unmarshal(r.UserMetadata, &userMeta)
	_ = json.Unmarshal(r.AppMetadata, &appMeta)

	return UserOutput{
		ID:           r.ID,
		Email:        r.Email,
		Name:         stringVal(userMeta, "name"),
		Phone:        stringVal(userMeta, "phone"),
		TaxID:        stringVal(userMeta, "tax_id"),
		SocialID:     stringVal(userMeta, "social_id"),
		City:         stringVal(userMeta, "city"),
		State:        stringVal(userMeta, "state"),
		ZIP:          stringVal(userMeta, "zip"),
		Street:       stringVal(userMeta, "street"),
		Complement:   stringVal(userMeta, "complement"),
		Neighborhood: stringVal(userMeta, "neighborhood"),
		Role:         stringVal(appMeta, "role"),
	}
}

func (g *GoTrueGatewayAdapter) Paginate(ctx context.Context, input PaginateInput) (*PaginateOutput, error) {
	page := max(input.Page, 1)
	pageSize := max(input.PageSize, 10)

	paginate := database.PaginateInput{Page: page - 1, PageSize: pageSize}

	const countSQL = `SELECT COUNT(*) FROM auth.users`
	const dataSQL = `SELECT id, email, raw_user_meta_data, raw_app_meta_data FROM auth.users LIMIT $1 OFFSET $2`

	data, count, err := database.Paginate[authUserRow](
		ctx, g.db.Pool,
		countSQL, nil,
		dataSQL, []any{pageSize, paginate.Offset()},
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrAuthServiceFailure, err)
	}

	users := make([]UserOutput, len(data))
	for i, row := range data {
		users[i] = row.toUserOutput()
	}

	return &PaginateOutput{
		Data:     users,
		MaxPages: paginate.CalcMaxPages(count),
	}, nil
}
