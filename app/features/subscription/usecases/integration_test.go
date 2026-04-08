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

func (t *TestSuite) Test_SubscriptionIntegration_Create() {
	t.Run("should be able to create", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, token)

		Body := &usecases.CreateInput{
			Items: []usecases.SubscriptionItemInput{
				{
					ProductID: fixtures.Product.Create(t.T(), nil, token),
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
			AutoRenew:   true,
			Notes:       "Notes",
		}

		id := ""
		fixtures.DefaultHTTP(t.T()).
			Request(http.MethodPost, fixtures.Subscription.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusCreated).
			Body().Decode(&id)
	})

	t.Run("should fail if required fields are empty", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, token)

		Body := map[string]any{
			// Optional fields
			"autoRenew": false,
			"notes":     "",
		}

		response := new(middlewares.ValidationErrorResponse)
		fixtures.DefaultHTTP(t.T()).
			Request(http.MethodPost, fixtures.Subscription.URI).
			WithHeader("Authorization", "Bearer "+token).
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
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, token)

		Body := map[string]any{
			"items": []usecases.SubscriptionItemInput{
				{
					ProductID: fixtures.Product.Create(t.T(), nil, token),
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
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(Body).
			Expect().
			Status(http.StatusCreated).
			Body().NotEmpty()
	})
}

func (t *TestSuite) Test_SubscriptionIntegration_GetByID() {
	t.Run("should be able to get by id", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, token)
		subscriptionID := fixtures.Subscription.Create(t.T(), nil, token)

		response := new(g.GetByIDOutput)
		fixtures.DefaultHTTP(t.T()).
			Request(http.MethodGet, fixtures.Subscription.URI+subscriptionID).
			WithHeader("Authorization", "Bearer "+token).
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

		t.Len(response.Items, 1)
		t.Equal(fixtures.Subscription.DEFAULT_NOTES, response.Notes)
		t.Equal(fixtures.Subscription.DEFAULT_TOTAL_CYCLES, response.TotalCycles)
		t.Equal(fixtures.Subscription.DEFAULT_TOTAL_CYCLES, response.RemainingCycles)
		t.EqualValues(0, response.CycleDiscount)
		t.Equal(fixtures.Subscription.DEFAULT_CUTOFF_OFFSET_MINUTES, response.CutoffOffsetMinutes)
		t.Equal(fixtures.Subscription.DEFAULT_AUTO_RENEW, response.AutoRenew)
	})

	t.Run("should return NotFound when get by id not found", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, token)
		fixtures.Subscription.Create(t.T(), nil, token)

		NON_EXISTING_ID := uuid.New().String()

		fixtures.DefaultHTTP(t.T()).
			Request(http.MethodGet, fixtures.Subscription.URI+NON_EXISTING_ID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound).
			Body().NotEmpty()
	})
}

func (t *TestSuite) Test_SubscriptionIntegration_Paginate() {
	t.Run("should paginate own subscriptions", func() {
		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, token)

		for range 3 {
			fixtures.Subscription.Create(t.T(), nil, token)
		}

		response := new(g.PaginateOutput)
		fixtures.DefaultHTTP(t.T()).
			Request(http.MethodGet, fixtures.Subscription.URI).
			WithHeader("Authorization", "Bearer "+token).
			WithQueryObject(map[string]any{"Page": 0, "PageSize": 10}).
			Expect().
			Status(http.StatusOK).
			JSON().Decode(&response)

		t.Len(response.Data, 3)
	})

	t.Run("should not see subscriptions from another tenant", func() {
		anotherTenant := fixtures.Tenant.Create(t.T(), nil)
		anotherToken := fixtures.Auth.UserToken(t.T(), anotherTenant.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, anotherToken)

		for range 3 {
			fixtures.Subscription.Create(t.T(), nil, anotherToken)
		}

		tenant := fixtures.Tenant.Create(t.T(), nil)
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)

		response := new(g.PaginateOutput)
		fixtures.DefaultHTTP(t.T()).
			Request(http.MethodGet, fixtures.Subscription.URI).
			WithHeader("Authorization", "Bearer "+token).
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
		token1 := fixtures.Auth.UserToken(t.T(), tenant1.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, token1)

		tenant2 := fixtures.Tenant.Create(t.T(), nil)
		token2 := fixtures.Auth.UserToken(t.T(), tenant2.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, token2)

		for range 2 {
			fixtures.Subscription.Create(t.T(), nil, token1)
		}
		for range 2 {
			fixtures.Subscription.Create(t.T(), nil, token2)
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
		token1 := fixtures.Auth.UserToken(t.T(), tenant1.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, token1)

		tenant2 := fixtures.Tenant.Create(t.T(), nil)
		token2 := fixtures.Auth.UserToken(t.T(), tenant2.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, token2)

		for range 2 {
			fixtures.Subscription.Create(t.T(), nil, token1)
		}
		fixtures.Subscription.Create(t.T(), nil, token2)

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
		token := fixtures.Auth.UserToken(t.T(), tenant.UserID)
		fixtures.Subscription.UpsertConfig(t.T(), nil, token)
		subscriptionID := fixtures.Subscription.Create(t.T(), nil, token)

		fixtures.DefaultHTTP(t.T()).
			Request(http.MethodDelete, fixtures.Subscription.URI+subscriptionID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusMethodNotAllowed)
	})
}
