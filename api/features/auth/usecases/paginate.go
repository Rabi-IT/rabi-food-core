package usecases

import (
	"context"

	g "github.com/Rabi-IT/rabi-food-core/features/auth/gateway"
)

type PaginateOutput struct {
	Data     []GetByIDOutput `json:"data"`
	MaxPages int             `json:"maxPages"`
}

func (c *AuthCase) Paginate(ctx context.Context, page, pageSize int) (*PaginateOutput, error) {
	out, err := c.gateway.Paginate(ctx, g.PaginateInput{
		Page:     page,
		PageSize: pageSize,
	})
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
