package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"
)

type adminListUsersResponse struct {
	Users []adminUserResponse `json:"users"`
}

func (g *GoTrueGatewayAdapter) Paginate(ctx context.Context, input PaginateInput) (*PaginateOutput, error) {
	page := input.Page
	if page < 1 {
		page = 1
	}

	pageSize := input.PageSize
	if pageSize < 1 {
		pageSize = 10
	}

	url := fmt.Sprintf("%s/auth/v1/admin/users?page=%d&per_page=%d", g.baseURL, page, pageSize)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errs.ErrAuthServiceFailure
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", g.serviceKey))

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, errs.ErrAuthServiceFailure
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errs.ErrAuthServiceFailure
	}

	var result adminListUsersResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, errs.ErrAuthServiceFailure
	}

	// Filter by tenant_id client-side — GoTrue Admin API does not support metadata filtering
	users := make([]UserOutput, 0, len(result.Users))
	for _, u := range result.Users {
		out := u.toUserOutput()

		if input.TenantID != nil && out.TenantID != *input.TenantID {
			continue
		}

		users = append(users, *out)
	}

	return &PaginateOutput{
		Data:     users,
		MaxPages: -1, // GoTrue does not expose total count in the Admin API
	}, nil
}
