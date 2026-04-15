package fixtures

import (
	"net/http"
	"testing"

	g "github.com/Rabi-IT/rabi-food-core/features/category/gateway"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
)

type categoryFixture struct {
	URI          string
	BackofficeURI string
}

var Category = categoryFixture{URI: "/category/", BackofficeURI: "/backoffice/category/"}

func (categoryFixture) Create(t *testing.T, input *g.CreateInput, auth RequestContext) string {
	t.Helper()
	Body := input
	if Body == nil {
		Body = &g.CreateInput{
			Name:        "Name",
			Description: "Description",
		}
	}

	id := ""
	httpexpect.Default(t, AppURL).
		Request(http.MethodPost, Category.URI).
		WithHeader("Authorization", "Bearer "+auth.Token).
		WithHeader("X-Tenant-ID", auth.TenantID).
		WithJSON(Body).
		Expect().
		Status(http.StatusCreated).
		Body().Decode(&id)

	return id
}

func (categoryFixture) GetByID(t *testing.T, id string, auth RequestContext) (g.GetByIDOutput, int) {
	t.Helper()
	found := g.GetByIDOutput{}

	obj := httpexpect.Default(t, AppURL).
		Request(http.MethodGet, Category.URI+id).
		WithHeader("Authorization", "Bearer "+auth.Token).
		WithHeader("X-Tenant-ID", auth.TenantID).
		Expect().Status(http.StatusOK)

	response := obj.Raw()

	obj.JSON().Object().Decode(&found)

	err := response.Body.Close()
	require.NoError(t, err)

	return found, response.StatusCode
}
