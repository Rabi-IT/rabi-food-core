package usecases

import (
	"context"

	g "github.com/Rabi-IT/rabi-food-core/features/auth/gateway"
)

type PatchInput struct {
	Email        *string `json:"email"        validate:"omitempty,email"`
	Name         *string `json:"name"`
	Phone        *string `json:"phone"`
	TaxID        *string `json:"taxId"`
	SocialID     *string `json:"socialId"`
	City         *string `json:"city"`
	State        *string `json:"state"`
	ZIP          *string `json:"zip"`
	Street       *string `json:"street"`
	Complement   *string `json:"complement"`
	Neighborhood *string `json:"neighborhood"`
}

func (c *AuthCase) Patch(ctx context.Context, id string, input *PatchInput) (bool, error) {
	return c.gateway.Patch(ctx, id, g.PatchInput{
		Email:        input.Email,
		Name:         input.Name,
		Phone:        input.Phone,
		TaxID:        input.TaxID,
		SocialID:     input.SocialID,
		City:         input.City,
		State:        input.State,
		ZIP:          input.ZIP,
		Street:       input.Street,
		Complement:   input.Complement,
		Neighborhood: input.Neighborhood,
	})
}
