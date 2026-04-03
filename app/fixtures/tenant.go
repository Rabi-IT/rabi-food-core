package fixtures

import (
	"net/http"
	g "rabi-food-core/features/tenant/gateway"
	c "rabi-food-core/features/tenant/usecases"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
)

type tenantFixture struct {
	URI string
}

var Tenant = tenantFixture{"/tenant/"}

func (tenantFixture) Create(t *testing.T, input *c.CreateInput) *c.CreateOutput {
	t.Helper()
	Body := input
	if Body == nil {
		Body = &c.CreateInput{
			Name:     "Name",
			UserName: "UserName",
			Phone:    "http://example.com/photo.png",
			Email:    "email@email.com",
		}
	}

	output := &c.CreateOutput{}
	httpexpect.Default(t, AppURL).
		Request(http.MethodPost, Tenant.URI).
		WithJSON(Body).
		Expect().Status(http.StatusCreated).
		JSON().Object().
		Decode(output)

	require.NotEqual(t, &c.CreateOutput{}, output)

	return output
}

func (tenantFixture) GetMe(t *testing.T, token string) g.GetByIDOutput {
	t.Helper()
	found := g.GetByIDOutput{}

	obj := httpexpect.Default(t, AppURL).
		Request(http.MethodGet, Tenant.URI+"me").
		WithHeader("Authorization", "Bearer "+token).
		Expect().Status(http.StatusOK)

	response := obj.Raw()

	obj.JSON().Object().Decode(&found)

	err := response.Body.Close()
	require.NoError(t, err)

	return found
}
