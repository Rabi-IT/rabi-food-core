package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"
)

type adminCreateUserRequest struct {
	Email        string         `json:"email"`
	Password     string         `json:"password"`
	EmailConfirm bool           `json:"email_confirm"`
	UserMetadata map[string]any `json:"user_metadata"`
	AppMetadata  map[string]any `json:"app_metadata"`
}

type adminCreateUserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func (g *GoTrueGatewayAdapter) SignUp(ctx context.Context, input SignUpInput) (*SignUpOutput, error) {
	body := adminCreateUserRequest{
		Email:        input.Email,
		Password:     input.Password,
		EmailConfirm: true,
		UserMetadata: map[string]any{
			"name":         input.Name,
			"phone":        input.Phone,
			"tax_id":       input.TaxID,
			"social_id":    input.SocialID,
			"tenant_id":    input.TenantID,
			"city":         input.City,
			"state":        input.State,
			"zip":          input.ZIP,
			"street":       input.Street,
			"complement":   input.Complement,
			"neighborhood": input.Neighborhood,
		},
		AppMetadata: map[string]any{
			"role":      input.Role,
			"tenant_id": input.TenantID,
		},
	}

	raw, err := json.Marshal(body)
	if err != nil {
		return nil, errs.ErrAuthServiceFailure
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.baseURL+"/auth/v1/admin/users", bytes.NewReader(raw))
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

	if resp.StatusCode != http.StatusOK {
		return nil, errs.ErrAuthServiceFailure
	}

	var result adminCreateUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, errs.ErrAuthServiceFailure
	}

	return &SignUpOutput{ID: result.ID, Email: result.Email}, nil
}
