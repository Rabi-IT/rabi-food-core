package tenant_case_test

import (
	"net/http"
	"rabi-food-core/fixtures"
	"rabi-food-core/libs/http/fiber_adapter/middlewares"
	"rabi-food-core/usecases/tenant_case"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite

	app *fixtures.App
}

func (t *TestSuite) SetupSuite() {
	t.app = fixtures.NewApp()
	t.app.Start(t.T())
}

func (t *TestSuite) SetupSubTest() {
	fixtures.CleanDatabase(t.T())
}

func (t *TestSuite) TearDownSuite() {
	t.app.Stop(t.T())
}

func TestMySuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (t *TestSuite) Test_TenantIntegration_Create() {
	t.Run("should be able to create", func() {
		Body := tenant_case.CreateInput{
			Name:     "Name",
			UserName: "UserName",
			Phone:    "Phone",
			Email:    "email@email.com",
		}

		var response tenant_case.CreateOutput
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.Tenant.URI).
			WithJSON(Body).
			Expect().Status(http.StatusCreated).
			JSON().Decode(&response)

		token := fixtures.Auth.UserToken(t.T(), response.UserID)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Tenant.URI+response.ID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().Status(http.StatusOK).
			JSON().Object().
			ContainsSubset(map[string]any{
				"id":   response.ID,
				"name": Body.Name,
			})
	})

	t.Run("should fail if required fields are empty", func() {
		Body := tenant_case.CreateInput{}

		response := new(middlewares.ValidationErrorResponse)
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.Tenant.URI).
			WithJSON(Body).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Decode(response)

		t.Len(response.Errors, 4)
		for _, e := range response.Errors {
			switch e.Field {
			case "Name":
				t.Equal("required", e.Tag)
			case "UserName":
				t.Equal("required", e.Tag)
			case "Phone":
				t.Equal("required", e.Tag)
			case "Email":
				t.Equal("required", e.Tag)
			default:
				t.Fail("unexpected validation error field: " + e.Field)
			}
		}
	})
}
