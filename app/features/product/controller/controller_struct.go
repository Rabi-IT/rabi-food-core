package controller

import (
	"rabi-food-core/features/product/usecases"
)

type ProductController struct {
	usecase *usecases.ProductCase
}

func New(usecase *usecases.ProductCase) *ProductController {
	return &ProductController{usecase}
}
