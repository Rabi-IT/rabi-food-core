package usecases_test

import (
	"net/http"
	"testing"

	"github.com/Rabi-IT/rabi-food-core/fixtures"
	"github.com/Rabi-IT/rabi-food-core/libs/http/middlewares"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite

	app *fixtures.Api
}

func (t *TestSuite) SetupSuite() {
	t.app = fixtures.NewApi()
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
	t.Run("should enroll a customer to a tenant", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)

		customerID := fixtures.Auth.SignUp(t.T(), fixtures.SignUpInput{
			Email:    "customer@email.com",
			Password: "password123",
			Name:     "Customer",
			Phone:    "11988888888",
			TaxID:    "12345678901",
		})
		token := fixtures.Auth.UserToken(t.T(), customerID)

		httpexpect.Default(t.T(), fixtures.ApiURL).
			Request(http.MethodPost, fixtures.Tenant.URI+tenant.ID+"/enroll-customers").
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
			httpexpect.Default(t.T(), fixtures.ApiURL).
				Request(http.MethodPost, fixtures.Tenant.URI+tenant.ID+"/enroll-customers").
				WithHeader("Authorization", "Bearer "+token).
				Expect().Status(http.StatusNoContent)
		}
	})
}

func (t *TestSuite) Test_TenantIntegration_Create() {
	t.Run("should create tenant with owner as member", func() {
		tenant := fixtures.Tenant.Create(t.T(), &fixtures.TenantCreateInput{
			Name:         "My Tenant",
			UserName:     "Owner Name",
			UserPhone:    "11999999999",
			UserEmail:    "owner@email.com",
			UserPassword: "password123",
			UserTaxID:    "12345678901",
		})

		token := fixtures.Auth.TenantOwnerToken(t.T(), tenant.UserID, tenant.ID)

		httpexpect.Default(t.T(), fixtures.ApiURL).
			Request(http.MethodGet, fixtures.Tenant.URI+"me").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("X-Tenant-ID", tenant.ID).
			Expect().Status(http.StatusOK).
			JSON().Object().
			ContainsSubset(map[string]any{
				"id":   tenant.ID,
				"name": "My Tenant",
			})
	})

	t.Run("should fail if name is missing", func() {
		userID := fixtures.Auth.SignUp(t.T(), fixtures.SignUpInput{
			Email:    "owner@email.com",
			Password: "password123",
			Name:     "Owner Name",
			Phone:    "11999999999",
			TaxID:    "12345678901",
		})
		token := fixtures.Auth.UserToken(t.T(), userID)

		response := new(middlewares.ValidationErrorResponse)
		httpexpect.Default(t.T(), fixtures.ApiURL).
			Request(http.MethodPost, "/tenant").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]string{"name": ""}).
			Expect().Status(http.StatusBadRequest).
			JSON().Decode(response)

		t.Len(response.Errors, 1)
		t.Equal("name", response.Errors[0].Field)
		t.Equal("required", response.Errors[0].Tag)
	})
}
