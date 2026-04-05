package controller

import (
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/features/category/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
	"github.com/Rabi-IT/rabi-food-core/libs/validator"

	"github.com/gofiber/fiber/v2"
)

// Create godoc
// @Summary Create category
// @Description Create a new category
// @Tags categories
// @Accept json
// @Produce text/plain
// @Param category body gateway.CreateInput true "Category data"
// @Success 201 {string} string "Created category ID"
// @Failure 400 {object} middlewares.ValidationErrorResponse "Validation errors"
// @Failure 500 {string} string "Internal server error"
// @Router /category/ [post].
func (c *CategoryController) Create(ctx *fiber.Ctx) error {
	data := gateway.CreateInput{}
	err := ctx.BodyParser(&data)
	if err != nil {
		return ctx.JSON(err)
	}

	err = validator.V.Struct(data)
	if err != nil {
		return err
	}

	uctx := ctx.UserContext()
	logger.GetWideEvent(uctx).Event = "create-category"
	id, err := c.usecase.Create(uctx, data)

	if err != nil {
		return err
	}

	return ctx.Status(http.StatusCreated).SendString(id)
}
