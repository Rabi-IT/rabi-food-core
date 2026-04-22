package gateway

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"
)

func (g *GoTrueGatewayAdapter) Delete(ctx context.Context, id string) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, g.baseURL+"/admin/users/"+id, nil)
	if err != nil {
		return false, fmt.Errorf("%w: %v", errs.ErrAuthServiceFailure, err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", g.serviceKey))

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("%w: %v", errs.ErrAuthServiceFailure, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return false, wrapResponseError(errs.ErrAuthServiceFailure, resp)
	}

	return true, nil
}
