package controller

import "rabi-food-core/features/subscription/usecases"

type SubscriptionController struct {
	usecase *usecases.SubscriptionCase
}

func New(usecase *usecases.SubscriptionCase) *SubscriptionController {
	return &SubscriptionController{usecase}
}
