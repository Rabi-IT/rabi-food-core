package usecases_test

import (
	"net/http"
	"testing"

	product_gateway "github.com/Rabi-IT/rabi-food-core/features/product/gateway"
	"github.com/Rabi-IT/rabi-food-core/fixtures"
	"github.com/Rabi-IT/rabi-food-core/libs/database"

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
		tenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)
		categoryID := fixtures.Category.Create(t.T(), nil, tenantAuth)

		body := product_gateway.CreateInput{
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
			WithHeader("Authorization", "Bearer "+tenantAuth.Token).
			WithJSON(body).
			Expect().
			Status(http.StatusCreated).
			Body().NotEmpty()
	})

	t.Run("should use session TenantID, ignoring tenantID in body", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		tenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)
		categoryID := fixtures.Category.Create(t.T(), nil, tenantAuth)

		anotherTenant := fixtures.Tenant.Create(t.T(), nil)
		body := product_gateway.CreateInput{
			TenantID:    anotherTenant.ID,
			Name:        "Name",
			Description: "Description",
			CategoryID:  categoryID,
			Price:       100,
		}

		productID := fixtures.Product.Create(t.T(), &body, tenantAuth)

		productFound, httpStatus := fixtures.Product.GetByID(t.T(), productID)
		t.Equal(http.StatusOK, httpStatus)
		t.Equal(tenant.ID, productFound.TenantID)
	})

	t.Run("should not be able to create as a regular user", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		userAuth := fixtures.Auth.UserAuth(t.T(), tenant.ID, tenant.UserID)
		tenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)
		categoryID := fixtures.Category.Create(t.T(), nil, tenantAuth)

		body := product_gateway.CreateInput{
			Name:        "Name",
			Description: "Description",
			CategoryID:  categoryID,
			Price:       100,
		}

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPost, fixtures.Product.URI).
			WithHeader("Authorization", "Bearer "+userAuth.Token).
			WithHeader("X-Tenant-ID", userAuth.TenantID).
			WithJSON(body).
			Expect().
			Status(http.StatusForbidden)
	})
}

func (t *TestSuite) Test_ProductIntegration_GetByID() {
	t.Run("should be able to get by id", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		tenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)
		productID := fixtures.Product.Create(t.T(), nil, tenantAuth)

		found, status := fixtures.Product.GetByID(t.T(), productID)

		t.Equal(http.StatusOK, status)
		t.Equal(productID, found.ID)
		t.Equal("Name", found.Name)
		t.Equal("Description", found.Description)
	})

	t.Run("should return NotFound when get by id not found", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		userAuth := fixtures.Auth.UserAuth(t.T(), tenant.ID, tenant.UserID)
		tenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)
		_ = fixtures.Product.Create(t.T(), nil, tenantAuth)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Product.URI+uuid.New().String()).
			WithHeader("Authorization", "Bearer "+userAuth.Token).
			Expect().
			Status(http.StatusNotFound).
			Body().NotEmpty()
	})
}

func (t *TestSuite) Test_ProductIntegration_Paginate() {
	t.Run("should be able to paginate", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		tenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)

		for range 15 {
			fixtures.Product.Create(t.T(), nil, tenantAuth)
		}

		response := new(product_gateway.PaginateOutput)
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Product.URI).
			WithHeader("X-Tenant-ID", tenantAuth.TenantID).
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

	t.Run("should not be able to paginate products from another tenant", func() {
		anotherTenant := fixtures.Tenant.Create(t.T(), nil)
		anotherTenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), anotherTenant.ID, anotherTenant.UserID)

		for range 5 {
			fixtures.Product.Create(t.T(), nil, anotherTenantAuth)
		}

		tenant := fixtures.Tenant.Create(t.T(), nil)

		response := new(product_gateway.PaginateOutput)
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Product.URI).
			WithHeader("X-Tenant-ID", tenant.ID).
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
		tenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)
		productID := fixtures.Product.Create(t.T(), nil, tenantAuth)

		body := product_gateway.PatchValues{
			Name:        new("Updated Name"),
			Description: new("Updated Description"),
		}

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPatch, fixtures.Product.URI+productID).
			WithHeader("Authorization", "Bearer "+tenantAuth.Token).
			WithJSON(body).
			Expect().
			Status(http.StatusOK).
			Body().NotEmpty()

		found, status := fixtures.Product.GetByID(t.T(), productID)

		t.Equal(http.StatusOK, status)
		t.Equal(productID, found.ID)
		t.Equal("Updated Name", found.Name)
		t.Equal("Updated Description", found.Description)
	})

	t.Run("should not be able to patch a product from another tenant", func() {
		anotherTenant := fixtures.Tenant.Create(t.T(), nil)
		anotherTenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), anotherTenant.ID, anotherTenant.UserID)
		anotherProductID := fixtures.Product.Create(t.T(), nil, anotherTenantAuth)

		tenant := fixtures.Tenant.Create(t.T(), nil)
		tenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)

		body := product_gateway.PatchValues{
			Name:        new("Updated Name"),
			Description: new("Updated Description"),
		}

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPatch, fixtures.Product.URI+anotherProductID).
			WithHeader("Authorization", "Bearer "+tenantAuth.Token).
			WithJSON(body).
			Expect().
			Status(http.StatusNotFound).
			Body().NotEmpty()
	})

	t.Run("should not be able to patch as a regular user", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		userAuth := fixtures.Auth.UserAuth(t.T(), tenant.ID, tenant.UserID)
		tenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)
		productID := fixtures.Product.Create(t.T(), nil, tenantAuth)

		body := product_gateway.PatchValues{
			Name: new("Updated Name"),
		}

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodPatch, fixtures.Product.URI+productID).
			WithHeader("Authorization", "Bearer "+userAuth.Token).
			WithHeader("X-Tenant-ID", userAuth.TenantID).
			WithJSON(body).
			Expect().
			Status(http.StatusForbidden)
	})
}

func (t *TestSuite) Test_ProductIntegration_BackofficePaginate() {
	t.Run("should see products from all tenants", func() {
		tenant1 := fixtures.Tenant.Create(t.T(), nil)
		tenantAuth1 := fixtures.Auth.TenantOwnerAuth(t.T(), tenant1.ID, tenant1.UserID)
		tenant2 := fixtures.Tenant.Create(t.T(), nil)
		tenantAuth2 := fixtures.Auth.TenantOwnerAuth(t.T(), tenant2.ID, tenant2.UserID)

		for range 3 {
			fixtures.Product.Create(t.T(), nil, tenantAuth1)
		}
		for range 3 {
			fixtures.Product.Create(t.T(), nil, tenantAuth2)
		}

		backofficeToken := fixtures.Auth.BackofficeToken(t.T(), tenant1.UserID)

		response := new(product_gateway.PaginateOutput)
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Product.BackofficeURI).
			WithHeader("Authorization", "Bearer "+backofficeToken).
			WithQueryObject(database.PaginateInput{Page: 0, PageSize: 10}).
			Expect().
			Status(http.StatusOK).
			JSON().Decode(&response)

		t.Len(response.Data, 6)
	})

	t.Run("should filter by tenantId when provided", func() {
		tenant1 := fixtures.Tenant.Create(t.T(), nil)
		tenantAuth1 := fixtures.Auth.TenantOwnerAuth(t.T(), tenant1.ID, tenant1.UserID)
		tenant2 := fixtures.Tenant.Create(t.T(), nil)
		tenantAuth2 := fixtures.Auth.TenantOwnerAuth(t.T(), tenant2.ID, tenant2.UserID)

		for range 3 {
			fixtures.Product.Create(t.T(), nil, tenantAuth1)
		}
		for range 2 {
			fixtures.Product.Create(t.T(), nil, tenantAuth2)
		}

		backofficeToken := fixtures.Auth.BackofficeToken(t.T(), tenant1.UserID)

		response := new(product_gateway.PaginateOutput)
		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Product.BackofficeURI).
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

func (t *TestSuite) Test_ProductIntegration_Delete() {
	t.Run("should be able to delete", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		userAuth := fixtures.Auth.UserAuth(t.T(), tenant.ID, tenant.UserID)
		tenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)
		productID := fixtures.Product.Create(t.T(), nil, tenantAuth)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodDelete, fixtures.Product.URI+productID).
			WithHeader("Authorization", "Bearer "+tenantAuth.Token).
			Expect().
			Status(http.StatusNoContent)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodGet, fixtures.Product.URI+productID).
			WithHeader("Authorization", "Bearer "+userAuth.Token).
			WithHeader("X-Tenant-ID", userAuth.TenantID).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("should not be able to delete a product from another tenant", func() {
		anotherTenant := fixtures.Tenant.Create(t.T(), nil)
		anotherTenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), anotherTenant.ID, anotherTenant.UserID)
		anotherProductID := fixtures.Product.Create(t.T(), nil, anotherTenantAuth)

		tenant := fixtures.Tenant.Create(t.T(), nil)
		tenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodDelete, fixtures.Product.URI+anotherProductID).
			WithHeader("Authorization", "Bearer "+tenantAuth.Token).
			Expect().
			Status(http.StatusNotFound).
			Body().NotEmpty()
	})

	t.Run("should not be able to delete as a regular user", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		userAuth := fixtures.Auth.UserAuth(t.T(), tenant.ID, tenant.UserID)
		tenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)
		productID := fixtures.Product.Create(t.T(), nil, tenantAuth)

		httpexpect.Default(t.T(), fixtures.AppURL).
			Request(http.MethodDelete, fixtures.Product.URI+productID).
			WithHeader("Authorization", "Bearer "+userAuth.Token).
			WithHeader("X-Tenant-ID", userAuth.TenantID).
			Expect().
			Status(http.StatusForbidden)
	})
}
