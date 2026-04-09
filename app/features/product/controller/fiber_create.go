package controller

import (
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/features/product/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
	"github.com/Rabi-IT/rabi-food-core/libs/validator"

	"github.com/gofiber/fiber/v2"
)

// Create godoc
// @Summary Create product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce text/plain
// @Param product body gateway.CreateInput true "Product data"
// @Success 201 {string} string "Created product ID"
// @Failure 400 {object} middlewares.ValidationErrorResponse "Validation errors"
// @Failure 500 {string} string "Internal server error"
// @Router /product/ [post].
func (c *ProductController) Create(ctx *fiber.Ctx) error {
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
	logger.GetWideEvent(uctx).Event = "create-product"
	id, err := c.usecase.Create(uctx, data)

	if err != nil {
		return err
	}

	return ctx.Status(http.StatusCreated).SendString(id)
}
