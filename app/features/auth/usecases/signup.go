package usecases

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/domain/auth"
	g "github.com/Rabi-IT/rabi-food-core/features/auth/gateway"
)

type TenantInput struct {
	Name string `json:"name" validate:"required"`
}

type SignUpInput struct {
	Email        string       `json:"email"        validate:"required,email"`
	Password     string       `json:"password"     validate:"required,min=8"`
	Name         string       `json:"name"         validate:"required"`
	Phone        string       `json:"phone"        validate:"required"`
	TaxID        string       `json:"taxId"        validate:"required"`
	SocialID     string       `json:"socialId"`
	City         string       `json:"city"`
	State        string       `json:"state"`
	ZIP          string       `json:"zip"`
	Street       string       `json:"street"`
	Complement   string       `json:"complement"`
	Neighborhood string       `json:"neighborhood"`
	Tenant       *TenantInput `json:"tenant"`
}

type SignUpOutput struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	TenantID string `json:"tenantId,omitempty"`
}

func (c *AuthCase) SignUp(ctx context.Context, input *SignUpInput) (*SignUpOutput, error) {
	role := auth.User
	var tenantID string

	if input.Tenant != nil {
		id, err := c.tenantCase.Create(ctx, input.Tenant.Name)
		if err != nil {
			return nil, err
		}
		tenantID = id
		role = auth.TenantOwner
	}

	out, err := c.gateway.SignUp(ctx, g.SignUpInput{
		Email:        input.Email,
		Password:     input.Password,
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
		Role:         string(role),
	})
	if err != nil {
		return nil, err
	}

	if tenantID != "" {
		if err := c.tenantCase.CreateMember(ctx, tenantID, out.ID, string(auth.TenantOwner)); err != nil {
			return nil, err
		}
	}

	return &SignUpOutput{ID: out.ID, Email: out.Email, TenantID: tenantID}, nil
}
