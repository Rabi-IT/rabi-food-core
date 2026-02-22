package product_case_test

import (
	"net/http"
	"rabi-food-core/fixtures"
	"rabi-food-core/libs/database"
	"rabi-food-core/libs/database/gateways/product_gateway"
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

func (t *TestSuite) Test_ProductIntegration_Create() {
	t.Run("should be able to create", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		categoryID := fixtures.Category.Create(t.T(), nil, token)

		Body := product_gateway.CreateInput{
			Name:        "Name",
			Photo:       "http://example.com/photo.png",
			Description: "Description",
			CategoryID:  categoryID,
			Unit:        "Unit",
			Price:       100,
			IsActive:    true,
		}

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.Product.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusCreated).
			Body().NotEmpty()
	})

	t.Run("should ignore provided tenantID and use token tenant when user is not backoffice", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		categoryID := fixtures.Category.Create(t.T(), nil, token)

		anotherTenant := fixtures.Tenant.Create(t.T(), nil)
		body := product_gateway.CreateInput{
			TenantID:    anotherTenant.ID,
			Name:        "Name",
			Description: "Description",
			CategoryID:  categoryID,
		}

		productID := fixtures.Product.Create(t.T(), &body, token)

		productFound, httpStatus := fixtures.Product.GetByID(t.T(), productID, token)
		t.Equal(http.StatusOK, httpStatus)
		t.Equal(tenant.ID, productFound.TenantID)
	})
}

func (t *TestSuite) Test_ProductIntegration_GetByID() {
	t.Run("should be able to get by id", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		productID := fixtures.Product.Create(t.T(), nil, token)

		found, status := fixtures.Product.GetByID(t.T(), productID, token)

		t.Equal(http.StatusOK, status)
		t.Equal(productID, found.ID)
		t.Equal("Name", found.Name)
		t.Equal("Description", found.Description)
	})

	t.Run("should return NotFound when get by id not found", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		_ = fixtures.Product.Create(t.T(), nil, token)

		NON_EXISTING_ID := uuid.New().String()

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Product.URI+NON_EXISTING_ID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound).
			Body().NotEmpty()
	})

	t.Run("should not be able to get a product from another tenant", func() {
		anotherTenant := fixtures.Tenant.Create(t.T(), nil)
		anotherToken := fixtures.Auth.UserToken(t.T(), anotherTenant.UserID)
		anotherProductID := fixtures.Product.Create(t.T(), nil, anotherToken)

		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Product.URI+anotherProductID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound).
			Body().NotEmpty()
	})
}

func (t *TestSuite) Test_ProductIntegration_Paginate() {
	t.Run("should be able to paginate", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		for range 15 {
			fixtures.Product.Create(t.T(), nil, token)
		}

		response := new(product_gateway.PaginateOutput)
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Product.URI).
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
			fixtures.Product.Create(t.T(), nil, anotherToken)
		}

		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		response := new(product_gateway.PaginateOutput)
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Product.URI).
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

func (t *TestSuite) Test_ProductIntegration_Patch() {
	t.Run("should be able to patch", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		productID := fixtures.Product.Create(t.T(), nil, token)

		Body := product_gateway.PatchValues{
			Name:        "Updated Name",
			Description: "Updated Description",
		}

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPatch, fixtures.Product.URI+productID).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusOK).
			Body().NotEmpty()

		found, status := fixtures.Product.GetByID(t.T(), productID, token)

		t.Equal(http.StatusOK, status)
		t.Equal(productID, found.ID)
		t.Equal("Updated Name", found.Name)
		t.Equal("Updated Description", found.Description)
	})

	t.Run("should not be able to patch a product from another tenant", func() {
		anotherTenant := fixtures.Tenant.Create(t.T(), nil)
		anotherToken := fixtures.Auth.UserToken(t.T(), anotherTenant.UserID)
		anotherProductID := fixtures.Product.Create(t.T(), nil, anotherToken)

		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		Body := product_gateway.PatchValues{
			Name:        "Updated Name",
			Description: "Updated Description",
		}

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPatch, fixtures.Product.URI+anotherProductID).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusNotFound).
			Body().NotEmpty()
	})
}

func (t *TestSuite) Test_ProductIntegration_Delete() {
	t.Run("should be able to delete", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		productID := fixtures.Product.Create(t.T(), nil, token)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodDelete, fixtures.Product.URI+productID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNoContent)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Product.URI+productID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("should not be able to delete a product from another tenant", func() {
		anotherTenant := fixtures.Tenant.Create(t.T(), nil)
		anotherToken := fixtures.Auth.UserToken(t.T(), anotherTenant.UserID)
		anotherProductID := fixtures.Product.Create(t.T(), nil, anotherToken)

		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodDelete, fixtures.Product.URI+anotherProductID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound).
			Body().NotEmpty()
	})
}
