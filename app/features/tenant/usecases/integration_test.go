package usecases_test

import (
	"net/http"
	"testing"

	auth_usecases "github.com/Rabi-IT/rabi-food-core/features/auth/usecases"
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

func (t *TestSuite) Test_TenantIntegration_EnrollCustomer() {
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
			Request(http.MethodPost, fixtures.Tenant.URI+tenant.ID+"/me/customers").
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
				Request(http.MethodPost, fixtures.Tenant.URI+tenant.ID+"/me/customers").
				WithHeader("Authorization", "Bearer "+token).
				Expect().Status(http.StatusNoContent)
		}
	})
}

func (t *TestSuite) Test_TenantIntegration_SignUpWithTenant() {
	t.Run("should create tenant and owner on signup", func() {
		body := auth_usecases.SignUpInput{
			Email:    "owner@email.com",
			Password: "password123",
			Name:     "Owner Name",
			Phone:    "11999999999",
			TaxID:    "12345678901",
			Tenant:   &auth_usecases.TenantInput{Name: "My Tenant"},
		}

		var response auth_usecases.SignUpOutput
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, "/auth/signup").
			WithJSON(body).
			Expect().Status(http.StatusCreated).
			JSON().Decode(&response)

		t.NotEmpty(response.ID)
		t.NotEmpty(response.TenantID)

		token := fixtures.Auth.TenantOwnerToken(t.T(), response.ID, response.TenantID)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Tenant.URI+"me").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("X-Tenant-ID", response.TenantID).
			Expect().Status(http.StatusOK).
			JSON().Object().
			ContainsSubset(map[string]any{
				"id":   response.TenantID,
				"name": body.Tenant.Name,
			})
	})

	t.Run("should fail if required fields are missing", func() {
		body := auth_usecases.SignUpInput{
			Tenant: &auth_usecases.TenantInput{},
		}

		response := new(middlewares.ValidationErrorResponse)
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, "/auth/signup").
			WithJSON(body).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Decode(response)

		t.Len(response.Errors, 6)
		for _, e := range response.Errors {
			switch e.Field {
			case "Email":
				t.Equal("required", e.Tag)
			case "Password":
				t.Equal("required", e.Tag)
			case "Name":
				t.Equal("required", e.Tag)
			case "Phone":
				t.Equal("required", e.Tag)
			case "TaxID":
				t.Equal("required", e.Tag)
			case "Tenant.Name":
				t.Equal("required", e.Tag)
			default:
				t.Fail("unexpected validation error field: " + e.Field)
			}
		}
	})
}
