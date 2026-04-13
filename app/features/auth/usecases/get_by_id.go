package usecases

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	g "github.com/Rabi-IT/rabi-food-core/features/auth/gateway"
)

type GetByIDOutput struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	TaxID        string `json:"taxId"`
	SocialID     string `json:"socialId"`
	City         string `json:"city"`
	State        string `json:"state"`
	ZIP          string `json:"zip"`
	Street       string `json:"street"`
	Complement   string `json:"complement"`
	Neighborhood string `json:"neighborhood"`
	Role         string `json:"role"`
}

func (c *AuthCase) GetByID(ctx context.Context, id string) (*GetByIDOutput, error) {
	session := app_context.GetSession(ctx)

	// Users can only fetch themselves; backoffice can fetch any user
	if session.Role.IsUser() {
		id = session.UserID
	}

	out, err := c.gateway.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if out == nil {
		return nil, nil
	}

	return &GetByIDOutput{
		ID:           out.ID,
		Email:        out.Email,
		Name:         out.Name,
		Phone:        out.Phone,
		TaxID:        out.TaxID,
		SocialID:     out.SocialID,
		City:         out.City,
		State:        out.State,
		ZIP:          out.ZIP,
		Street:       out.Street,
		Complement:   out.Complement,
		Neighborhood: out.Neighborhood,
		Role:         out.Role,
	}, nil
}

func (c *AuthCase) GetMe(ctx context.Context) (*GetByIDOutput, error) {
	session := app_context.GetSession(ctx)
	return c.GetByID(ctx, session.UserID)
}

func fromGateway(u *g.UserOutput) *GetByIDOutput {
	return &GetByIDOutput{
		ID:           u.ID,
		Email:        u.Email,
		Name:         u.Name,
		Phone:        u.Phone,
		TaxID:        u.TaxID,
		SocialID:     u.SocialID,
		City:         u.City,
		State:        u.State,
		ZIP:          u.ZIP,
		Street:       u.Street,
		Complement:   u.Complement,
		Neighborhood: u.Neighborhood,
		Role:         u.Role,
	}
}
