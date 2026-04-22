package usecases

import (
	"context"

	g "github.com/Rabi-IT/rabi-food-core/features/auth/gateway"
)

type SignInInput struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type TokenOutput struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
	TokenType    string `json:"tokenType"`
}

func (c *AuthCase) SignIn(ctx context.Context, input *SignInInput) (*TokenOutput, error) {
	out, err := c.gateway.SignIn(ctx, g.SignInInput{
		Email:    input.Email,
		Password: input.Password,
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
