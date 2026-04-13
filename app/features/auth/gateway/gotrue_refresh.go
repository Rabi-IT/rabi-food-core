package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"
)

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (g *GoTrueGatewayAdapter) RefreshToken(ctx context.Context, refreshToken string) (*TokenOutput, error) {
	body := refreshRequest{RefreshToken: refreshToken}

	raw, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrAuthServiceFailure, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.baseURL+"/token?grant_type=refresh_token", bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrAuthServiceFailure, err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", g.serviceKey))

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrAuthServiceFailure, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusUnauthorized {
		return nil, wrapResponseError(errs.ErrInvalidRefreshToken, resp)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, wrapResponseError(errs.ErrAuthServiceFailure, resp)
	}

	var result TokenOutput
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrAuthServiceFailure, err)
	}

	return &result, nil
}
