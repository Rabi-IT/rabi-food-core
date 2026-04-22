package gateway

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"
)

func (g *GoTrueGatewayAdapter) SignOut(ctx context.Context, accessToken string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.baseURL+"/logout", nil)
	if err != nil {
		return fmt.Errorf("%w: %w", errs.ErrAuthServiceFailure, err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %w", errs.ErrAuthServiceFailure, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return wrapResponseError(errs.ErrAuthServiceFailure, resp)
	}

	return nil
}
