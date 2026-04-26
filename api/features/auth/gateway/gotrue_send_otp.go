package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"
)

type sendOTPRequest struct {
	Email      string         `json:"email"`
	CreateUser bool           `json:"create_user"`
	Data       map[string]any `json:"data,omitempty"`
}

func (g *GoTrueGatewayAdapter) SendOTP(ctx context.Context, input SendOTPInput) error {
	body := sendOTPRequest{
		Email:      input.Email,
		CreateUser: false,
	}

	raw, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("%w: %w", errs.ErrAuthServiceFailure, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.baseURL+"/otp", bytes.NewReader(raw))
	if err != nil {
		return fmt.Errorf("%w: %w", errs.ErrAuthServiceFailure, err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.serviceKey)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %w", errs.ErrAuthServiceFailure, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return wrapResponseError(errs.ErrAuthServiceFailure, resp)
	}

	return nil
}
