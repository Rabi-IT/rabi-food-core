package usecases

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"
)

type GetBySlugOutput struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Language string `json:"language"`
}

func (c *TenantCase) GetBySlug(ctx context.Context, slug string) (*GetBySlugOutput, error) {
	out, err := c.gateway.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if out == nil {
		return nil, errs.ErrTenantNotFound
	}

	return &GetBySlugOutput{
		ID:       out.ID,
		Name:     out.Name,
		Slug:     out.Slug,
		Language: out.Language,
	}, nil
}
