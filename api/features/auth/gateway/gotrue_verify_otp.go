package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"
)

type verifyOTPRequest struct {
	Type  string `json:"type"`
	Token string `json:"token"`
	Email string `json:"email"`
}

func (g *GoTrueGatewayAdapter) VerifyOTP(ctx context.Context, input VerifyOTPInput) (*TokenOutput, error) {
	body := verifyOTPRequest{
		Type:  "email",
		Token: input.Code,
		Email: input.Email,
	}

	raw, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrAuthServiceFailure, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.baseURL+"/verify", bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrAuthServiceFailure, err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.serviceKey)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrAuthServiceFailure, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnprocessableEntity || resp.StatusCode == http.StatusBadRequest {
		return nil, wrapResponseError(errs.ErrInvalidOTP, resp)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, wrapResponseError(errs.ErrAuthServiceFailure, resp)
	}

	var result TokenOutput
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrAuthServiceFailure, err)
	}

	return &result, nil
}
