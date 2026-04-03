package controller

import (
	"github.com/Rabi-IT/rabi-food-core/features/product/usecases"
)

type ProductController struct {
	usecase *usecases.ProductCase
}

func New(usecase *usecases.ProductCase) *ProductController {
	return &ProductController{usecase}
}
