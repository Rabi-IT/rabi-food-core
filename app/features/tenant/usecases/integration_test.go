package usecases_test

import (
	"net/http"
	"testing"

	tenant_case "github.com/Rabi-IT/rabi-food-core/features/tenant/usecases"
	"github.com/Rabi-IT/rabi-food-core/fixtures"
	"github.com/Rabi-IT/rabi-food-core/libs/http/middlewares"

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

func (t *TestSuite) Test_TenantIntegration_RegisterCustomer() {
	t.Run("should register a customer to a tenant", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)

		customerID := fixtures.Auth.SignUp(t.T(), fixtures.SignUpInput{
			Email:    "customer@email.com",
			Password: "password123",
			Name:     "Customer",
			Phone:    "11988888888",
			TaxID:    "12345678901",
		})
		token := fixtures.Auth.UserToken(t.T(), customerID)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.Tenant.URI+tenant.ID+"/customers").
			WithHeader("Authorization", "Bearer "+token).
			Expect().Status(http.StatusNoContent)
	})

	t.Run("should be idempotent", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)

		customerID := fixtures.Auth.SignUp(t.T(), fixtures.SignUpInput{
			Email:    "customer@email.com",
			Password: "password123",
			Name:     "Customer",
			Phone:    "11988888888",
			TaxID:    "12345678901",
		})
		token := fixtures.Auth.UserToken(t.T(), customerID)

		for range 2 {
			httpexpect.Default(t.T(), fixtures.AppURL).
				Request(http.MethodPost, fixtures.Tenant.URI+tenant.ID+"/customers").
				WithHeader("Authorization", "Bearer "+token).
				Expect().Status(http.StatusNoContent)
		}
	})
}

func (t *TestSuite) Test_TenantIntegration_Create() {
	t.Run("should be able to create", func() {
		body := tenant_case.CreateInput{
			Name:         "Name",
			UserName:     "UserName",
			UserPhone:    "11999999999",
			UserEmail:    "email@email.com",
			UserPassword: "password123",
		}

		var response tenant_case.CreateOutput
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.Tenant.URI).
			WithJSON(body).
			Expect().Status(http.StatusCreated).
			JSON().Decode(&response)

		token := fixtures.Auth.UserToken(t.T(), response.UserID)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Tenant.URI+response.ID).
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("X-Tenant-ID", response.ID).
			Expect().Status(http.StatusOK).
			JSON().Object().
			ContainsSubset(map[string]any{
				"id":   response.ID,
				"name": body.Name,
			})
	})

	t.Run("should fail if required fields are missing", func() {
		body := tenant_case.CreateInput{}

		response := new(middlewares.ValidationErrorResponse)
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.Tenant.URI).
			WithJSON(body).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Decode(response)

		t.Len(response.Errors, 5)
		for _, e := range response.Errors {
			switch e.Field {
			case "Name":
				t.Equal("required", e.Tag)
			case "UserName":
				t.Equal("required", e.Tag)
			case "UserPhone":
				t.Equal("required", e.Tag)
			case "UserEmail":
				t.Equal("required", e.Tag)
			case "UserPassword":
				t.Equal("required", e.Tag)
			default:
				t.Fail("unexpected validation error field: " + e.Field)
			}
		}
	})
}
