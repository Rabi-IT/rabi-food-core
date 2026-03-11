package subscription_controller

import (
	"rabi-food-core/usecases/subscription_case"
)

type SubscriptionController struct {
	usecase *subscription_case.SubscriptionCase
}

func New(usecase *subscription_case.SubscriptionCase) *SubscriptionController {
	return &SubscriptionController{usecase}
}
