package controller

import "github.com/Rabi-IT/rabi-food-core/features/subscription/usecases"

type SubscriptionBackofficeController struct {
	usecase *usecases.SubscriptionCase
}

func NewBackoffice(usecase *usecases.SubscriptionCase) *SubscriptionBackofficeController {
	return &SubscriptionBackofficeController{usecase}
}
