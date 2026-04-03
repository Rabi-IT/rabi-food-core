package usecases

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	"github.com/Rabi-IT/rabi-food-core/domain/auth"
	g "github.com/Rabi-IT/rabi-food-core/features/tenant/gateway"
	user "github.com/Rabi-IT/rabi-food-core/features/user/usecases"
)

type CreateInput struct {
	Name     string `validate:"required"`
	UserName string `validate:"required"`
	Phone    string `validate:"required"`
	Email    string `validate:"required,email"`
}

type CreateOutput struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
}

func (c *TenantCase) Create(ctx context.Context, input CreateInput) (out CreateOutput, err error) {
	tenantId, err := c.gateway.Create(g.CreateInput{
		Name: input.Name,
	})

	if err != nil {
		return
	}

	ctx = app_context.WithSession(ctx, &app_context.UserSession{
		TenantID: tenantId,
		Role:     auth.User,
	})
	userId, err := c.userCase.Create(ctx, &user.CreateInput{
		Name:  input.UserName,
		Phone: input.Phone,
		Email: input.Email,
	})

	if err != nil {
		return
	}

	return CreateOutput{
		UserID: userId,
		ID:     tenantId,
	}, nil
}
