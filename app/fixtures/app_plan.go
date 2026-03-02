package fixtures

import (
	"net/http"
	"rabi-food-core/libs/database/gateways/app_plan_gateway"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
)

type appPlanFixture struct {
	URI string
}

var AppPlan = appPlanFixture{"/category/"}

func (appPlanFixture) Create(t *testing.T, input *app_plan_gateway.CreateInput, token string) string {
	t.Helper()
	Body := input
	if Body == nil {
		Body = &app_plan_gateway.CreateInput{
			Name:        "Name",
			Description: "Description",
		}
	}

	id := ""
	httpexpect.Default(t, AppURL).
		Request(http.MethodPost, AppPlan.URI).
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(Body).
		Expect().
		Status(http.StatusCreated).
		Body().Decode(&id)

	return id
}

func (appPlanFixture) GetByID(t *testing.T, id string, token string) (app_plan_gateway.GetByIDOutput, int) {
	t.Helper()
	found := app_plan_gateway.GetByIDOutput{}

	obj := httpexpect.Default(t, AppURL).
		Request(http.MethodGet, AppPlan.URI+id).
		WithHeader("Authorization", "Bearer "+token).
		Expect().Status(http.StatusOK)

	response := obj.Raw()

	obj.JSON().Object().Decode(&found)

	err := response.Body.Close()
	require.NoError(t, err)

	return found, response.StatusCode
}
