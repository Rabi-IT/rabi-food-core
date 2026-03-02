package app_plan_case_test

import (
	"net/http"
	"rabi-food-core/fixtures"
	g "rabi-food-core/libs/database/gateways/app_plan_gateway"
	"rabi-food-core/libs/http/fiber_adapter/middlewares"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
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

func (t *TestSuite) Test_AppPlanIntegration_Create() {
	t.Run("should be able to create", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		Body := g.CreateInput{
			Name:        "Name",
			Description: "Description",
		}

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.AppPlan.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusCreated).
			Body().NotEmpty()
	})

	t.Run("should be able to get by id", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		categoryID := fixtures.AppPlan.Create(t.T(), nil, token)

		found, status := fixtures.AppPlan.GetByID(t.T(), categoryID, token)

		t.Equal(http.StatusOK, status)
		t.Equal(categoryID, found.ID)
		t.Equal("Name", found.Name)
		t.Equal("Description", found.Description)
	})

	t.Run("should return NotFound when get by id not found", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		_ = fixtures.AppPlan.Create(t.T(), nil, token)

		NON_EXISTING_ID := uuid.New().String()

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.AppPlan.URI+NON_EXISTING_ID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound).
			Body().NotEmpty()
	})

	t.Run("should fail if required fields are empty", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		Body := g.CreateInput{
			// Required fields
			Name: "",

			// Optional fields
			Description: "Description",
		}

		response := new(middlewares.ValidationErrorResponse)
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.AppPlan.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Decode(&response)

		t.Len(response.Errors, 1)
		t.Equal("Name", response.Errors[0].Field)
		t.Equal("required", response.Errors[0].Tag)
	})

	t.Run("should not fail if optional fields are empty", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		Body := g.CreateInput{
			// Required fields
			Name: "Name",

			// Optional fields
			Description: "",
		}

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.AppPlan.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusCreated).
			Body().NotEmpty()
	})
}

func (t *TestSuite) Test_AppPlanIntegration_GetByID() {
	t.Run("should be able to get by id", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		categoryID := fixtures.AppPlan.Create(t.T(), nil, token)

		found, status := fixtures.AppPlan.GetByID(t.T(), categoryID, token)

		t.Equal(http.StatusOK, status)
		t.Equal(categoryID, found.ID)
		t.Equal("Name", found.Name)
		t.Equal("Description", found.Description)
	})

	t.Run("should return NotFound when get by id not found", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		_ = fixtures.AppPlan.Create(t.T(), nil, token)

		NON_EXISTING_ID := uuid.New().String()

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.AppPlan.URI+NON_EXISTING_ID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound).
			Body().NotEmpty()
	})

	t.Run("should not be able to get a category from another tenant", func() {
		anotherTenant := fixtures.Tenant.Create(t.T(), nil)
		anotherToken := fixtures.Auth.UserToken(t.T(), anotherTenant.UserID)
		anotherAppPlanID := fixtures.AppPlan.Create(t.T(), nil, anotherToken)

		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.AppPlan.URI+anotherAppPlanID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound).
			Body().NotEmpty()
	})
}

func (t *TestSuite) Test_AppPlanIntegration_Patch() {
	t.Run("should be able to patch", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		categoryID := fixtures.AppPlan.Create(t.T(), nil, token)

		Body := g.PatchValues{
			Name:        "Updated Name",
			Description: "Updated Description",
		}

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPatch, fixtures.AppPlan.URI+categoryID).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusOK).
			Body().NotEmpty()

		found, status := fixtures.AppPlan.GetByID(t.T(), categoryID, token)

		t.Equal(http.StatusOK, status)
		t.Equal(categoryID, found.ID)
		t.Equal("Updated Name", found.Name)
		t.Equal("Updated Description", found.Description)
	})

	t.Run("should not be able to patch a category from another tenant", func() {
		anotherTenant := fixtures.Tenant.Create(t.T(), nil)
		anotherToken := fixtures.Auth.UserToken(t.T(), anotherTenant.UserID)
		anotherAppPlanID := fixtures.AppPlan.Create(t.T(), nil, anotherToken)

		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		Body := g.PatchValues{
			Name:        "Updated Name",
			Description: "Updated Description",
		}

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPatch, fixtures.AppPlan.URI+anotherAppPlanID).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusNotFound).
			Body().NotEmpty()
	})
}

func (t *TestSuite) Test_AppPlanIntegration_Delete() {
	t.Run("should be able to delete", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		categoryID := fixtures.AppPlan.Create(t.T(), nil, token)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodDelete, fixtures.AppPlan.URI+categoryID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNoContent)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.AppPlan.URI+categoryID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("should not be able to delete a category from another tenant", func() {
		anotherTenant := fixtures.Tenant.Create(t.T(), nil)
		anotherToken := fixtures.Auth.UserToken(t.T(), anotherTenant.UserID)
		anotherAppPlanID := fixtures.AppPlan.Create(t.T(), nil, anotherToken)

		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodDelete, fixtures.AppPlan.URI+anotherAppPlanID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound).
			Body().NotEmpty()
	})
}
