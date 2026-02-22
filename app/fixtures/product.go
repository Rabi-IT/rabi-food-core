package fixtures

import (
	"net/http"
	"rabi-food-core/libs/database/gateways/product_gateway"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
)

type productFixture struct {
	URI string
}

var Product = productFixture{"/product/"}

func (productFixture) Create(t *testing.T, input *product_gateway.CreateInput, token string) string {
	t.Helper()
	Body := input
	if Body == nil {
		categoryID := Category.Create(t, nil, token)
		Body = &product_gateway.CreateInput{
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
	httpexpect.Default(t, AppURL).
		Request(http.MethodPost, Product.URI).
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(Body).
		Expect().
		Status(http.StatusCreated).
		Body().Decode(&id)

	return id
}

func (productFixture) GetByID(t *testing.T, id string, token string) (product_gateway.GetByIDOutput, int) {
	t.Helper()
	found := product_gateway.GetByIDOutput{}

	obj := httpexpect.Default(t, AppURL).
		Request(http.MethodGet, Product.URI+id).
		WithHeader("Authorization", "Bearer "+token).
		Expect().Status(http.StatusOK)

	response := obj.Raw()

	obj.JSON().Object().Decode(&found)

	err := response.Body.Close()
	require.NoError(t, err)

	return found, response.StatusCode
}
