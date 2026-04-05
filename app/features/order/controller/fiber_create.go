package controller

import (
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/features/order/usecases"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
	"github.com/Rabi-IT/rabi-food-core/libs/validator"

	"github.com/gofiber/fiber/v2"
)

// Create godoc
// @Summary Create order
// @Description Create a new order
// @Tags orders
// @Accept json
// @Produce text/plain
// @Param order body usecases.CreateInput true "Order data"
// @Success 201 {string} string "Created order ID"
// @Failure 400 {object} middlewares.ValidationErrorResponse "Validation errors"
// @Failure 500 {string} string "Internal server error"
// @Router /order/ [post].
func (c *OrderController) Create(ctx *fiber.Ctx) error {
	data := usecases.CreateInput{}
	err := ctx.BodyParser(&data)
	if err != nil {
		return ctx.JSON(err)
	}

	err = validator.V.Struct(data)
	if err != nil {
		return err
	}

	uctx := ctx.UserContext()
	logger.GetWideEvent(uctx).Event = "create-order"
	id, err := c.usecase.Create(uctx, data)

	if err != nil {
		return err
	}

	return ctx.Status(http.StatusCreated).SendString(id)
}
