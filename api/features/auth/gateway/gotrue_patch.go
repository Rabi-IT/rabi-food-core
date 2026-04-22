package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"
)

type adminUpdateUserRequest struct {
	Email        *string        `json:"email,omitempty"`
	UserMetadata map[string]any `json:"user_metadata,omitempty"`
}

func (g *GoTrueGatewayAdapter) Patch(ctx context.Context, id string, input PatchInput) (bool, error) {
	meta := map[string]any{}

	if input.Name != nil {
		meta["name"] = *input.Name
	}

	if input.Phone != nil {
		meta["phone"] = *input.Phone
	}

	if input.TaxID != nil {
		meta["tax_id"] = *input.TaxID
	}

	if input.SocialID != nil {
		meta["social_id"] = *input.SocialID
	}

	if input.City != nil {
		meta["city"] = *input.City
	}

	if input.State != nil {
		meta["state"] = *input.State
	}

	if input.ZIP != nil {
		meta["zip"] = *input.ZIP
	}

	if input.Street != nil {
		meta["street"] = *input.Street
	}

	if input.Complement != nil {
		meta["complement"] = *input.Complement
	}

	if input.Neighborhood != nil {
		meta["neighborhood"] = *input.Neighborhood
	}

	body := adminUpdateUserRequest{
		Email:        input.Email,
		UserMetadata: meta,
	}

	raw, err := json.Marshal(body)
	if err != nil {
		return false, fmt.Errorf("%w: %w", errs.ErrAuthServiceFailure, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, g.baseURL+"/admin/users/"+id, bytes.NewReader(raw))
	if err != nil {
		return false, fmt.Errorf("%w: %w", errs.ErrAuthServiceFailure, err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.serviceKey)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("%w: %w", errs.ErrAuthServiceFailure, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	if resp.StatusCode != http.StatusOK {
		return false, wrapResponseError(errs.ErrAuthServiceFailure, resp)
	}

	return true, nil
}
