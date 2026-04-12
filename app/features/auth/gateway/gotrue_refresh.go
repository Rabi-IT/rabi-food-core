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
		return nil, errs.ErrAuthServiceFailure
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.baseURL+"/auth/v1/token?grant_type=refresh_token", bytes.NewReader(raw))
	if err != nil {
		return nil, errs.ErrAuthServiceFailure
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", g.serviceKey))

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, errs.ErrAuthServiceFailure
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusUnauthorized {
		return nil, errs.ErrInvalidRefreshToken
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errs.ErrAuthServiceFailure
	}

	var result TokenOutput
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, errs.ErrAuthServiceFailure
	}

	return &result, nil
}
