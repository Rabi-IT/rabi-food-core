package usecases

import (
	"context"

	g "github.com/Rabi-IT/rabi-food-core/features/auth/gateway"
)

type VerifyOTPInput struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code"  validate:"required"`
}

func (c *AuthCase) VerifyOTP(ctx context.Context, input *VerifyOTPInput) (*TokenOutput, error) {
	out, err := c.gateway.VerifyOTP(ctx, g.VerifyOTPInput{
		Email: input.Email,
		Code:  input.Code,
	})
	if err != nil {
		return nil, err
	}

	return &TokenOutput{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
		ExpiresIn:    out.ExpiresIn,
		TokenType:    out.TokenType,
	}, nil
}
