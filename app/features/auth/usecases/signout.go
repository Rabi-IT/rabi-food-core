package usecases

import "context"

func (c *AuthCase) SignOut(ctx context.Context, accessToken string) error {
	return c.gateway.SignOut(ctx, accessToken)
}
