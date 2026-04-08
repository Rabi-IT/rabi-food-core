package usecases_test

import (
	"net/http"
	"testing"

	"github.com/Rabi-IT/rabi-food-core/features/category/gateway"
	"github.com/Rabi-IT/rabi-food-core/fixtures"
	"github.com/Rabi-IT/rabi-food-core/libs/database"
	"github.com/Rabi-IT/rabi-food-core/libs/http/middlewares"

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

func (t *TestSuite) Test_CategoryIntegration_Create() {
	t.Run("should be able to create", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		Body := gateway.CreateInput{
			Name:        "Name",
			Description: "Description",
		}

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.Category.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusCreated).
			Body().NotEmpty()
	})

	t.Run("should be able to get by id", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		categoryID := fixtures.Category.Create(t.T(), nil, token)

		found, status := fixtures.Category.GetByID(t.T(), categoryID, token)

		t.Equal(http.StatusOK, status)
		t.Equal(categoryID, found.ID)
		t.Equal("Name", found.Name)
		t.Equal("Description", found.Description)
	})

	t.Run("should return NotFound when get by id not found", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		_ = fixtures.Category.Create(t.T(), nil, token)

		NON_EXISTING_ID := uuid.New().String()

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Category.URI+NON_EXISTING_ID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound).
			Body().NotEmpty()
	})

	t.Run("should fail if required fields are empty", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		Body := gateway.CreateInput{
			// Required fields
			Name: "",

			// Optional fields
			Description: "Description",
		}

		response := new(middlewares.ValidationErrorResponse)
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.Category.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Decode(&response)

		t.Len(response.Errors, 1)
		t.Equal("name", response.Errors[0].Field)
		t.Equal("required", response.Errors[0].Tag)
	})

	t.Run("should not fail if optional fields are empty", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		Body := gateway.CreateInput{
			// Required fields
			Name: "Name",

			// Optional fields
			Description: "",
		}

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.Category.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusCreated).
			Body().NotEmpty()
	})

	t.Run("should ignore provided tenantID and use token tenant when user is not backoffice", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		anotherTenant := fixtures.Tenant.Create(t.T(), nil)

		Body := gateway.CreateInput{
			TenantID:    anotherTenant.ID,
			Name:        "Name",
			Description: "Description",
		}

		categoryID := fixtures.Category.Create(t.T(), &Body, token)

		categoryFound, httpStatus := fixtures.Category.GetByID(t.T(), categoryID, token)
		t.Equal(http.StatusOK, httpStatus)
		t.Equal(tenant.ID, categoryFound.TenantID)
	})
}

func (t *TestSuite) Test_CategoryIntegration_GetByID() {
	t.Run("should be able to get by id", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		categoryID := fixtures.Category.Create(t.T(), nil, token)

		found, status := fixtures.Category.GetByID(t.T(), categoryID, token)

		t.Equal(http.StatusOK, status)
		t.Equal(categoryID, found.ID)
		t.Equal("Name", found.Name)
		t.Equal("Description", found.Description)
	})

	t.Run("should return NotFound when get by id not found", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		_ = fixtures.Category.Create(t.T(), nil, token)

		NON_EXISTING_ID := uuid.New().String()

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Category.URI+NON_EXISTING_ID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound).
			Body().NotEmpty()
	})

	t.Run("should not be able to get a category from another tenant", func() {
		anotherTenant := fixtures.Tenant.Create(t.T(), nil)
		anotherToken := fixtures.Auth.UserToken(t.T(), anotherTenant.UserID)
		anotherCategoryID := fixtures.Category.Create(t.T(), nil, anotherToken)

		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Category.URI+anotherCategoryID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound).
			Body().NotEmpty()
	})
}

func (t *TestSuite) Test_CategoryIntegration_Paginate() {
	t.Run("should be able to paginate", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		for range 15 {
			fixtures.Category.Create(t.T(), nil, token)
		}

		response := new(gateway.PaginateOutput)
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Category.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithQueryObject(database.PaginateInput{
				Page:     0,
				PageSize: 10,
			}).
			Expect().
			Status(http.StatusOK).
			JSON().Decode(&response)

		t.Len(response.Data, 10)
		t.Equal(2, response.MaxPages)

		for i := range 10 {
			t.NotEmpty(response.Data[i].ID)
			t.Equal("Name", response.Data[i].Name)
			t.Equal("Description", response.Data[i].Description)
		}
	})

	t.Run("should not be able to paginate categories from another tenant", func() {
		anotherTenant := fixtures.Tenant.Create(t.T(), nil)
		anotherToken := fixtures.Auth.UserToken(t.T(), anotherTenant.UserID)

		for range 5 {
			fixtures.Category.Create(t.T(), nil, anotherToken)
		}

		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		response := new(gateway.PaginateOutput)
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Category.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithQueryObject(database.PaginateInput{
				Page:     0,
				PageSize: 10,
			}).
			Expect().
			Status(http.StatusOK).
			JSON().Decode(&response)

		t.Empty(response.Data)
		t.Equal(0, response.MaxPages)
	})
}

func (t *TestSuite) Test_CategoryIntegration_Patch() {
	t.Run("should be able to patch", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		categoryID := fixtures.Category.Create(t.T(), nil, token)

		Body := gateway.PatchValues{
			Name:        new("Updated Name"),
			Description: new("Updated Description"),
		}

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPatch, fixtures.Category.URI+categoryID).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusOK).
			Body().NotEmpty()

		found, status := fixtures.Category.GetByID(t.T(), categoryID, token)

		t.Equal(http.StatusOK, status)
		t.Equal(categoryID, found.ID)
		t.Equal("Updated Name", found.Name)
		t.Equal("Updated Description", found.Description)
	})

	t.Run("should not be able to patch a category from another tenant", func() {
		anotherTenant := fixtures.Tenant.Create(t.T(), nil)
		anotherToken := fixtures.Auth.UserToken(t.T(), anotherTenant.UserID)
		anotherCategoryID := fixtures.Category.Create(t.T(), nil, anotherToken)

		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		Body := gateway.PatchValues{
			Name:        new("Updated Name"),
			Description: new("Updated Description"),
		}

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPatch, fixtures.Category.URI+anotherCategoryID).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusNotFound).
			Body().NotEmpty()
	})
}

func (t *TestSuite) Test_CategoryIntegration_BackofficePaginate() {
	t.Run("should see categories from all tenants", func() {
		tenant1 := fixtures.Tenant.Create(t.T(), nil)
		token1 := fixtures.Auth.UserToken(t.T(), tenant1.UserID)
		tenant2 := fixtures.Tenant.Create(t.T(), nil)
		token2 := fixtures.Auth.UserToken(t.T(), tenant2.UserID)

		for range 3 {
			fixtures.Category.Create(t.T(), nil, token1)
		}
		for range 3 {
			fixtures.Category.Create(t.T(), nil, token2)
		}

		backofficeToken := fixtures.Auth.BackofficeToken(t.T(), tenant1.UserID)

		response := new(gateway.PaginateOutput)
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Category.BackofficeURI).
			WithHeader("Authorization", "Bearer "+backofficeToken).
			WithQueryObject(database.PaginateInput{Page: 0, PageSize: 10}).
			Expect().
			Status(http.StatusOK).
			JSON().Decode(&response)

		t.Len(response.Data, 6)
	})

	t.Run("should filter by tenantId when provided", func() {
		tenant1 := fixtures.Tenant.Create(t.T(), nil)
		token1 := fixtures.Auth.UserToken(t.T(), tenant1.UserID)
		tenant2 := fixtures.Tenant.Create(t.T(), nil)
		token2 := fixtures.Auth.UserToken(t.T(), tenant2.UserID)

		for range 3 {
			fixtures.Category.Create(t.T(), nil, token1)
		}
		for range 2 {
			fixtures.Category.Create(t.T(), nil, token2)
		}

		backofficeToken := fixtures.Auth.BackofficeToken(t.T(), tenant1.UserID)

		response := new(gateway.PaginateOutput)
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Category.BackofficeURI).
			WithHeader("Authorization", "Bearer "+backofficeToken).
			WithQueryObject(map[string]any{
				"Page":     0,
				"PageSize": 10,
				"tenantId": tenant1.ID,
			}).
			Expect().
			Status(http.StatusOK).
			JSON().Decode(&response)

		t.Len(response.Data, 3)
	})
}

func (t *TestSuite) Test_CategoryIntegration_Delete() {
	t.Run("should be able to delete", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		categoryID := fixtures.Category.Create(t.T(), nil, token)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodDelete, fixtures.Category.URI+categoryID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNoContent)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Category.URI+categoryID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("should not be able to delete a category from another tenant", func() {
		anotherTenant := fixtures.Tenant.Create(t.T(), nil)
		anotherToken := fixtures.Auth.UserToken(t.T(), anotherTenant.UserID)
		anotherCategoryID := fixtures.Category.Create(t.T(), nil, anotherToken)

		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodDelete, fixtures.Category.URI+anotherCategoryID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound).
			Body().NotEmpty()
	})
}
