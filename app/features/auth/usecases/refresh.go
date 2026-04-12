package usecases

import "context"

type RefreshInput struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

func (c *AuthCase) Refresh(ctx context.Context, input *RefreshInput) (*TokenOutput, error) {
	out, err := c.gateway.RefreshToken(ctx, input.RefreshToken)
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
