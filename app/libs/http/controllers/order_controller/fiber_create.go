package order_controller

import (
	"net/http"
	"rabi-food-core/libs/validator"
	"rabi-food-core/usecases/order_case"

	"github.com/gofiber/fiber/v2"
)

// CreateOrder godoc
// @Summary Create order
// @Description Create a new order
// @Tags orders
// @Accept json
// @Produce text/plain
// @Param order body order_case.CreateInput true "Order data"
// @Success 201 {string} string "Created order ID"
// @Failure 400 {object} middlewares.ValidationErrorResponse "Validation errors"
// @Failure 500 {string} string "Internal server error"
// @Router /order/ [post]
func (c *OrderController) Create(ctx *fiber.Ctx) error {
	data := order_case.CreateInput{}
	err := ctx.BodyParser(&data)
	if err != nil {
		return ctx.JSON(err)
	}

	err = validator.V.Struct(data)
	if err != nil {
		return err
	}

	id, err := c.usecase.Create(ctx.UserContext(), data)

	if err != nil {
		return err
	}

	return ctx.Status(http.StatusCreated).SendString(id)
}
