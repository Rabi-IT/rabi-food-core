package controller

import "github.com/Rabi-IT/rabi-food-core/features/category/usecases"

type CategoryBackofficeController struct {
	usecase *usecases.CategoryCase
}

func NewBackoffice(usecase *usecases.CategoryCase) *CategoryBackofficeController {
	return &CategoryBackofficeController{usecase}
}
