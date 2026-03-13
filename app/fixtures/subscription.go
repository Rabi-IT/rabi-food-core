package fixtures

import (
	"net/http"
	"rabi-food-core/libs/database/gateways/subscription_gateway"
	"rabi-food-core/usecases/subscription_case"
	"testing"

	"github.com/stretchr/testify/require"
)

type subscriptionFixture struct {
	URI                            string
	DEFAULT_DISCOUNT_RULE          subscription_gateway.DiscountRule
	DEFAULT_CUTOFF_OFFSET_MINUTES  uint16
	DEFAULT_MAX_ATTEMPTS_PER_ORDER uint8

	DEFAULT_TOTAL_CYCLES        uint
	DEFAULT_AUTO_RENEW          bool
	DEFAULT_NOTES               string
	DEFAULT_DELIVERY_WEEKDAY    uint8
	DEFAULT_DELIVERY_START_HOUR uint8
	DEFAULT_DELIVERY_END_HOUR   uint8
}

var (
	Subscription = subscriptionFixture{
		URI: "/subscription/",
		DEFAULT_DISCOUNT_RULE: subscription_gateway.DiscountRule{
			CyclesThreshold: 5,  //nolint:mnd
			Discount:        10, //nolint:mnd
		},
		DEFAULT_CUTOFF_OFFSET_MINUTES:  uint16(0),
		DEFAULT_MAX_ATTEMPTS_PER_ORDER: uint8(1),

		DEFAULT_TOTAL_CYCLES:        uint(1),
		DEFAULT_AUTO_RENEW:          true,
		DEFAULT_NOTES:               "Notes",
		DEFAULT_DELIVERY_WEEKDAY:    uint8(0),
		DEFAULT_DELIVERY_START_HOUR: uint8(10),
		DEFAULT_DELIVERY_END_HOUR:   uint8(12),
	}
)

func (subscriptionFixture) UpsertConfig(t *testing.T, input *subscription_gateway.UpsertConfigInput, token string) {
	t.Helper()
	Body := input
	if Body == nil {
		Body = &subscription_gateway.UpsertConfigInput{
			MaxAttemptsPerOrder: Subscription.DEFAULT_MAX_ATTEMPTS_PER_ORDER,
			DiscountRules: []subscription_gateway.DiscountRule{
				Subscription.DEFAULT_DISCOUNT_RULE,
			},
			CutoffOffsetMinutes: Subscription.DEFAULT_CUTOFF_OFFSET_MINUTES,
		}
	}

	DefaultHTTP(t).
		Request(http.MethodPut, Subscription.URI+"config").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(Body).
		Expect().
		StatusList(http.StatusOK, http.StatusCreated)
}

func (subscriptionFixture) Create(t *testing.T, input *subscription_case.CreateInput, token string) string {
	t.Helper()
	Body := input
	if Body == nil {
		Body = &subscription_case.CreateInput{
			Items: []subscription_case.SubscriptionItemInput{
				{
					ProductID: Product.Create(t, nil, token),
					Quantity:  1,
				},
			},
			DeliveryDays: []subscription_gateway.DeliveryDay{
				{
					Weekday:   Subscription.DEFAULT_DELIVERY_WEEKDAY,
					StartHour: Subscription.DEFAULT_DELIVERY_START_HOUR,
					EndHour:   Subscription.DEFAULT_DELIVERY_END_HOUR,
				},
			},
			TotalCycles: Subscription.DEFAULT_TOTAL_CYCLES,
			AutoRenew:   Subscription.DEFAULT_AUTO_RENEW,
			Notes:       Subscription.DEFAULT_NOTES,
		}
	}

	id := ""
	DefaultHTTP(t).
		Request(http.MethodPost, Subscription.URI).
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(Body).
		Expect().
		Status(http.StatusCreated).
		Body().Decode(&id)

	return id
}

func (subscriptionFixture) GetByID(t *testing.T, id string, token string) (subscription_gateway.GetByIDOutput, int) {
	t.Helper()
	found := subscription_gateway.GetByIDOutput{}

	obj := DefaultHTTP(t).
		Request(http.MethodGet, Subscription.URI+id).
		WithHeader("Authorization", "Bearer "+token).
		Expect().Status(http.StatusOK)

	response := obj.Raw()

	obj.JSON().Object().Decode(&found)

	err := response.Body.Close()
	require.NoError(t, err)

	return found, response.StatusCode
}
