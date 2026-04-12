package usecases

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	g "github.com/Rabi-IT/rabi-food-core/features/auth/gateway"
)

type PaginateOutput struct {
	Data     []GetByIDOutput `json:"data"`
	MaxPages int             `json:"maxPages"`
}

func (c *AuthCase) Paginate(ctx context.Context, page, pageSize int) (*PaginateOutput, error) {
	session := app_context.GetSession(ctx)

	input := g.PaginateInput{
		Page:     page,
		PageSize: pageSize,
	}

	// Users see only their tenant; backoffice sees all
	if session.Role.IsUser() {
		input.TenantID = &session.TenantID
	}

	out, err := c.gateway.Paginate(ctx, input)
	if err != nil {
		return nil, err
	}

	data := make([]GetByIDOutput, 0, len(out.Data))
	for i := range out.Data {
		data = append(data, *fromGateway(&out.Data[i]))
	}

	return &PaginateOutput{
		Data:     data,
		MaxPages: out.MaxPages,
	}, nil
}
