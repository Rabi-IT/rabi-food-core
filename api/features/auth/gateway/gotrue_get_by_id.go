package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"
)

type adminUserResponse struct {
	ID           string         `json:"id"`
	Email        string         `json:"email"`
	UserMetadata map[string]any `json:"user_metadata"`
	AppMetadata  map[string]any `json:"app_metadata"`
}

func (u *adminUserResponse) toUserOutput() *UserOutput {
	return &UserOutput{
		ID:           u.ID,
		Email:        u.Email,
		Name:         stringVal(u.UserMetadata, "name"),
		Phone:        stringVal(u.UserMetadata, "phone"),
		TaxID:        stringVal(u.UserMetadata, "tax_id"),
		SocialID:     stringVal(u.UserMetadata, "social_id"),
		City:         stringVal(u.UserMetadata, "city"),
		State:        stringVal(u.UserMetadata, "state"),
		ZIP:          stringVal(u.UserMetadata, "zip"),
		Street:       stringVal(u.UserMetadata, "street"),
		Complement:   stringVal(u.UserMetadata, "complement"),
		Neighborhood: stringVal(u.UserMetadata, "neighborhood"),
		Role:         stringVal(u.AppMetadata, "role"),
	}
}

func stringVal(m map[string]any, key string) string {
	if m == nil {
		return ""
	}

	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}

	return fmt.Sprint(v)
}

func (g *GoTrueGatewayAdapter) GetByID(ctx context.Context, id string) (*UserOutput, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, g.baseURL+"/admin/users/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrAuthServiceFailure, err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", g.serviceKey))

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrAuthServiceFailure, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, wrapResponseError(errs.ErrAuthServiceFailure, resp)
	}

	var result adminUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrAuthServiceFailure, err)
	}

	return result.toUserOutput(), nil
}
