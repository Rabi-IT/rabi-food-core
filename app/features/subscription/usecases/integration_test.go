package usecases_test

import (
	"net/http"
	"rabi-food-core/features/subscription"
	g "rabi-food-core/features/subscription/gateway"
	"rabi-food-core/features/subscription/usecases"
	"rabi-food-core/fixtures"
	"rabi-food-core/libs/http/middlewares"
	"testing"

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
