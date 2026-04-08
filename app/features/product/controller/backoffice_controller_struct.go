package controller

import (
	"github.com/Rabi-IT/rabi-food-core/features/product/usecases"
)

type ProductBackofficeController struct {
	usecase *usecases.ProductCase
}

func NewBackoffice(usecase *usecases.ProductCase) *ProductBackofficeController {
	return &ProductBackofficeController{usecase}
}
