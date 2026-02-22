package fixtures

import (
	"net/http"
	"rabi-food-core/libs/database/gateways/category_gateway"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
)

type categoryFixture struct {
	URI string
}

var Category = categoryFixture{"/category/"}

func (categoryFixture) Create(t *testing.T, input *category_gateway.CreateInput, token string) string {
	t.Helper()
	Body := input
	if Body == nil {
		Body = &category_gateway.CreateInput{
			Name:        "Name",
			Description: "Description",
		}
	}

	id := ""
	httpexpect.Default(t, AppURL).
		Request(http.MethodPost, Category.URI).
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(Body).
		Expect().
		Status(http.StatusCreated).
		Body().Decode(&id)

	return id
}

func (categoryFixture) GetByID(t *testing.T, id string, token string) (category_gateway.GetByIDOutput, int) {
	t.Helper()
	found := category_gateway.GetByIDOutput{}

	obj := httpexpect.Default(t, AppURL).
		Request(http.MethodGet, Category.URI+id).
		WithHeader("Authorization", "Bearer "+token).
		Expect().Status(http.StatusOK)

	response := obj.Raw()

	obj.JSON().Object().Decode(&found)

	err := response.Body.Close()
	require.NoError(t, err)

	return found, response.StatusCode
}
