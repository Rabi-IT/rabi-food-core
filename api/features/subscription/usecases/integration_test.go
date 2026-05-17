package usecases_test

import (
	"net/http"
	"testing"

	"github.com/Rabi-IT/rabi-food-core/features/subscription"
	g "github.com/Rabi-IT/rabi-food-core/features/subscription/gateway"
	"github.com/Rabi-IT/rabi-food-core/features/subscription/usecases"
	"github.com/Rabi-IT/rabi-food-core/fixtures"
	"github.com/Rabi-IT/rabi-food-core/libs/http/middlewares"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite

	api *fixtures.Api
}

func (t *TestSuite) SetupSuite() {
	t.api = fixtures.NewApi()
	t.api.Start(t.T())
}

func (t *TestSuite) SetupSubTest() {
	fixtures.CleanDatabase(t.T())
}

func (t *TestSuite) TearDownSuite() {
	t.api.Stop(t.T())
}

func TestMySuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (t *TestSuite) Test_SubscriptionIntegration_Create() {
	t.Run("should be able to create", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		userAuth := fixtures.Auth.UserAuth(t.T(), tenant.ID, tenant.UserID)
		tenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, userAuth)

		Body := &usecases.CreateInput{
			Items: []usecases.SubscriptionItemInput{
				{
					ProductID: fixtures.Product.Create(t.T(), nil, tenantAuth),
					Quantity:  1,
				},
			},
			DeliveryDays: []g.DeliveryDay{
				{
					Weekday:   0,
					StartHour: 10,
					EndHour:   12,
				},
			},
			TotalCycles: 10,
			Notes:       "Notes",
		}

		var ids []string
		fixtures.DefaultHTTP(t.T()).
			Request(http.MethodPost, fixtures.Subscription.URI).
			WithHeader("Authorization", "Bearer "+userAuth.Token).
			WithHeader("X-Tenant-ID", userAuth.TenantID).
			WithJSON(Body).
			Expect().
			Status(http.StatusCreated).
			JSON().Decode(&ids)

		t.NotEmpty(ids)
	})

	t.Run("should fail if required fields are empty", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		userAuth := fixtures.Auth.UserAuth(t.T(), tenant.ID, tenant.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, userAuth)

		Body := map[string]any{
			"autoRenew": false,
			"notes":     "",
		}

		response := new(middlewares.ValidationErrorResponse)
		fixtures.DefaultHTTP(t.T()).
			Request(http.MethodPost, fixtures.Subscription.URI).
			WithHeader("Authorization", "Bearer "+userAuth.Token).
			WithHeader("X-Tenant-ID", userAuth.TenantID).
			WithJSON(Body).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Decode(&response)

		t.Len(response.Errors, 3)
		requiredFields := []string{"items", "deliveryDays", "totalCycles"}
		for i := range response.Errors {
			t.Contains(requiredFields, response.Errors[i].Field)
			t.Equal("required", response.Errors[i].Tag)
		}
	})

	t.Run("should not fail if optional fields are empty", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		userAuth := fixtures.Auth.UserAuth(t.T(), tenant.ID, tenant.UserID)
		tenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, userAuth)

		Body := map[string]any{
			"items": []usecases.SubscriptionItemInput{
				{
					ProductID: fixtures.Product.Create(t.T(), nil, tenantAuth),
					Quantity:  1,
				},
			},
			"deliveryDays": []g.DeliveryDay{
				{
					Weekday:   0,
					StartHour: 10,
					EndHour:   12,
				},
			},
			"totalCycles": 10,
		}

		fixtures.DefaultHTTP(t.T()).
			Request(http.MethodPost, fixtures.Subscription.URI).
			WithHeader("Authorization", "Bearer "+userAuth.Token).
			WithHeader("X-Tenant-ID", userAuth.TenantID).
			WithJSON(Body).
			Expect().
			Status(http.StatusCreated).
			Body().NotEmpty()
	})
}

func (t *TestSuite) Test_SubscriptionIntegration_GetByID() {
	t.Run("should be able to get by id", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		userAuth := fixtures.Auth.UserAuth(t.T(), tenant.ID, tenant.UserID)
		tenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, userAuth)
		subscriptionID := fixtures.Subscription.Create(t.T(), nil, tenantAuth)

		response := new(g.GetByIDOutput)
		fixtures.DefaultHTTP(t.T()).
			Request(http.MethodGet, fixtures.Subscription.URI+subscriptionID).
			WithHeader("Authorization", "Bearer "+userAuth.Token).
			WithHeader("X-Tenant-ID", userAuth.TenantID).
			Expect().
			Status(http.StatusOK).
			JSON().Decode(&response)

		t.Equal(subscriptionID, response.ID)
		t.Equal(tenant.ID, response.TenantID)
		t.Equal(tenant.UserID, response.UserID)
		t.Equal(subscription.StatusActive, response.Status)

		t.Len(response.DeliveryDays, 1)
		deliveryDay := response.DeliveryDays[0]
		t.Equal(fixtures.Subscription.DEFAULT_DELIVERY_WEEKDAY, deliveryDay.Weekday)
		t.Equal(fixtures.Subscription.DEFAULT_DELIVERY_START_HOUR, deliveryDay.StartHour)
		t.Equal(fixtures.Subscription.DEFAULT_DELIVERY_END_HOUR, deliveryDay.EndHour)

		t.NotEmpty(response.ProductID)
		t.Equal(fixtures.Subscription.DEFAULT_NOTES, response.Notes)
		t.Equal(fixtures.Subscription.DEFAULT_TOTAL_CYCLES, response.TotalCycles)
		t.Equal(fixtures.Subscription.DEFAULT_TOTAL_CYCLES, response.RemainingCycles)
		t.EqualValues(0, response.CycleDiscount)
		t.Equal(fixtures.Subscription.DEFAULT_CUTOFF_OFFSET_MINUTES, response.CutoffOffsetMinutes)
		t.True(response.AutoRenew)
	})

	t.Run("should return NotFound when get by id not found", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		userAuth := fixtures.Auth.UserAuth(t.T(), tenant.ID, tenant.UserID)
		tenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, userAuth)
		fixtures.Subscription.Create(t.T(), nil, tenantAuth)

		fixtures.DefaultHTTP(t.T()).
			Request(http.MethodGet, fixtures.Subscription.URI+uuid.New().String()).
			WithHeader("Authorization", "Bearer "+userAuth.Token).
			WithHeader("X-Tenant-ID", userAuth.TenantID).
			Expect().
			Status(http.StatusNotFound).
			Body().NotEmpty()
	})
}

func (t *TestSuite) Test_SubscriptionIntegration_Paginate() {
	t.Run("should paginate own subscriptions", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		userAuth := fixtures.Auth.UserAuth(t.T(), tenant.ID, tenant.UserID)
		tenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, userAuth)

		for range 3 {
			fixtures.Subscription.Create(t.T(), nil, tenantAuth)
		}

		response := new(g.PaginateOutput)
		fixtures.DefaultHTTP(t.T()).
			Request(http.MethodGet, fixtures.Subscription.URI).
			WithHeader("Authorization", "Bearer "+userAuth.Token).
			WithHeader("X-Tenant-ID", userAuth.TenantID).
			WithQueryObject(map[string]any{"Page": 0, "PageSize": 10}).
			Expect().
			Status(http.StatusOK).
			JSON().Decode(&response)

		t.Len(response.Data, 3)
	})

	t.Run("should not see subscriptions from another tenant", func() {
		anotherTenant := fixtures.Tenant.Create(t.T(), nil)
		anotherTenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), anotherTenant.ID, anotherTenant.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, anotherTenantAuth)

		for range 3 {
			fixtures.Subscription.Create(t.T(), nil, anotherTenantAuth)
		}

		tenant := fixtures.Tenant.Create(t.T(), nil)
		userAuth := fixtures.Auth.UserAuth(t.T(), tenant.ID, tenant.UserID)

		response := new(g.PaginateOutput)
		fixtures.DefaultHTTP(t.T()).
			Request(http.MethodGet, fixtures.Subscription.URI).
			WithHeader("Authorization", "Bearer "+userAuth.Token).
			WithHeader("X-Tenant-ID", userAuth.TenantID).
			WithQueryObject(map[string]any{"Page": 0, "PageSize": 10}).
			Expect().
			Status(http.StatusOK).
			JSON().Decode(&response)

		t.Empty(response.Data)
	})
}

func (t *TestSuite) Test_SubscriptionIntegration_BackofficePaginate() {
	t.Run("should see subscriptions from all tenants", func() {
		tenant1 := fixtures.Tenant.Create(t.T(), nil)
		tenantAuth1 := fixtures.Auth.TenantOwnerAuth(t.T(), tenant1.ID, tenant1.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, tenantAuth1)

		tenant2 := fixtures.Tenant.Create(t.T(), nil)
		tenantAuth2 := fixtures.Auth.TenantOwnerAuth(t.T(), tenant2.ID, tenant2.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, tenantAuth2)

		for range 2 {
			fixtures.Subscription.Create(t.T(), nil, tenantAuth1)
		}
		for range 2 {
			fixtures.Subscription.Create(t.T(), nil, tenantAuth2)
		}

		backofficeToken := fixtures.Auth.BackofficeToken(t.T(), tenant1.UserID)

		response := new(g.PaginateOutput)
		fixtures.DefaultHTTP(t.T()).
			Request(http.MethodGet, fixtures.Subscription.BackofficeURI).
			WithHeader("Authorization", "Bearer "+backofficeToken).
			WithQueryObject(map[string]any{"Page": 0, "PageSize": 10}).
			Expect().
			Status(http.StatusOK).
			JSON().Decode(&response)

		t.Len(response.Data, 4)
	})

	t.Run("should filter by tenantId when provided", func() {
		tenant1 := fixtures.Tenant.Create(t.T(), nil)
		tenantAuth1 := fixtures.Auth.TenantOwnerAuth(t.T(), tenant1.ID, tenant1.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, tenantAuth1)

		tenant2 := fixtures.Tenant.Create(t.T(), nil)
		tenantAuth2 := fixtures.Auth.TenantOwnerAuth(t.T(), tenant2.ID, tenant2.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, tenantAuth2)

		for range 2 {
			fixtures.Subscription.Create(t.T(), nil, tenantAuth1)
		}
		fixtures.Subscription.Create(t.T(), nil, tenantAuth2)

		backofficeToken := fixtures.Auth.BackofficeToken(t.T(), tenant1.UserID)

		response := new(g.PaginateOutput)
		fixtures.DefaultHTTP(t.T()).
			Request(http.MethodGet, fixtures.Subscription.BackofficeURI).
			WithHeader("Authorization", "Bearer "+backofficeToken).
			WithQueryObject(map[string]any{
				"Page":     0,
				"PageSize": 10,
				"tenantId": tenant1.ID,
			}).
			Expect().
			Status(http.StatusOK).
			JSON().Decode(&response)

		t.Len(response.Data, 2)
	})
}

func (t *TestSuite) Test_SubscriptionIntegration_Delete() {
	t.Run("should not be able to delete", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		userAuth := fixtures.Auth.UserAuth(t.T(), tenant.ID, tenant.UserID)
		tenantAuth := fixtures.Auth.TenantOwnerAuth(t.T(), tenant.ID, tenant.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, userAuth)
		subscriptionID := fixtures.Subscription.Create(t.T(), nil, tenantAuth)

		fixtures.DefaultHTTP(t.T()).
			Request(http.MethodDelete, fixtures.Subscription.URI+subscriptionID).
			WithHeader("Authorization", "Bearer "+userAuth.Token).
			WithHeader("X-Tenant-ID", userAuth.TenantID).
			Expect().
			Status(http.StatusMethodNotAllowed)
	})
}
