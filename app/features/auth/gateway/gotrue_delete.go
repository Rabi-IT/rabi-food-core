package gateway

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"
)

func (g *GoTrueGatewayAdapter) Delete(ctx context.Context, id string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, g.baseURL+"/auth/v1/admin/users/"+id, nil)
	if err != nil {
		return errs.ErrAuthServiceFailure
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", g.serviceKey))

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return errs.ErrAuthServiceFailure
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return errs.ErrUserNotFound
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return errs.ErrAuthServiceFailure
	}

	return nil
}
