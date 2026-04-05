package usecases

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	"github.com/Rabi-IT/rabi-food-core/domain/auth"
	g "github.com/Rabi-IT/rabi-food-core/features/user/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
)

type CreateInput struct {
	Name         string `validate:"required"`
	Email        string `validate:"required,email"`
	TenantID     string
	Photo        string `validate:"omitempty,url"`
	TaxID        string
	City         string
	State        string
	Phone        string
	ZIP          string
	SocialID     string
	Neighborhood string
	Street       string
	Complement   string
}

func (c *UserCase) Create(ctx context.Context, input *CreateInput) (string, error) {
	session := app_context.GetSession(ctx)

	tenantId := ""
	if session.Role.IsUser() {
		tenantId = session.TenantID
	} else if session.Role.IsBackoffice() {
		tenantId = input.TenantID
	}

	id, err := c.gateway.Create(g.CreateInput{
		TenantID:     tenantId,
		City:         input.City,
		State:        input.State,
		ZIP:          input.ZIP,
		Phone:        input.Phone,
		Email:        input.Email,
		Photo:        input.Photo,
		TaxID:        input.TaxID,
		SocialID:     input.SocialID,
		Street:       input.Street,
		Complement:   input.Complement,
		Neighborhood: input.Neighborhood,
		Name:         input.Name,
		Role:         auth.User,
	})
	logger.GetWideEvent(ctx).SetUserID(id)

	if err != nil {
		return "", err
	}

	return id, nil
}
