package fixtures

import (
	"net/http"
	"testing"

	g "github.com/Rabi-IT/rabi-food-core/features/product/gateway"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
)

type productFixture struct {
	URI           string
	BackofficeURI string
}

var Product = productFixture{URI: "/product/", BackofficeURI: "/backoffice/product/"}

func (productFixture) Create(t *testing.T, input *g.CreateInput, auth RequestContext) string {
	t.Helper()
	Body := input
	if Body == nil {
		categoryID := Category.Create(t, nil, auth)
		Body = &g.CreateInput{
			Name:        "Name",
			Photo:       "http://example.com/photo.png",
			Description: "Description",
			CategoryID:  categoryID,
			Unit:        "Unit",
			Price:       100, //nolint:mnd
			IsActive:    true,
		}
	}

	id := ""
	httpexpect.Default(t, ApiURL).
		Request(http.MethodPost, Product.URI).
		WithHeader("Authorization", "Bearer "+auth.Token).
		WithHeader("X-Tenant-ID", auth.TenantID).
		WithJSON(Body).
		Expect().
		Status(http.StatusCreated).
		Body().Decode(&id)

	return id
}

func (productFixture) GetByID(t *testing.T, id string) (g.GetByIDOutput, int) {
	t.Helper()
	found := g.GetByIDOutput{}

	obj := httpexpect.Default(t, ApiURL).
		Request(http.MethodGet, Product.URI+id).
		Expect().Status(http.StatusOK)

	response := obj.Raw()

	obj.JSON().Object().Decode(&found)

	err := response.Body.Close()
	require.NoError(t, err)

	return found, response.StatusCode
}
