package gateway

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"
)

func (g *GoTrueGatewayAdapter) SignOut(ctx context.Context, accessToken string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.baseURL+"/auth/v1/logout", nil)
	if err != nil {
		return errs.ErrAuthServiceFailure
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return errs.ErrAuthServiceFailure
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return errs.ErrAuthServiceFailure
	}

	return nil
}
