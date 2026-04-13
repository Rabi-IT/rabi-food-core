package usecases

import (
	"context"

	g "github.com/Rabi-IT/rabi-food-core/features/tenant/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
)

type CreateInput struct {
	Name         string `validate:"required"`
	UserName     string `validate:"required"`
	UserPhone    string `validate:"required"`
	UserEmail    string `validate:"required,email"`
	UserPassword string `validate:"required,min=8"`
}

type CreateOutput struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
}

func (c *TenantCase) Create(ctx context.Context, input CreateInput) (out CreateOutput, err error) {
	tenantID, err := c.gateway.Create(g.CreateInput{Name: input.Name})
	logger.GetWideEvent(ctx).SetTenantID(tenantID)

	if err != nil {
		return
	}

	userID, err := c.userCreator.SignUpTenantOwner(ctx, input.UserEmail, input.UserPassword, input.UserName, input.UserPhone)
	if err != nil {
		return
	}

	return CreateOutput{ID: tenantID, UserID: userID}, nil
}
