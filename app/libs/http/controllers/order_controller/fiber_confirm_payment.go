package order_controller

import (
	"rabi-food-core/libs/http/fiber_adapter/parser"
	"rabi-food-core/libs/validator"
	"rabi-food-core/usecases/order_case"

	"github.com/gofiber/fiber/v2"
)

func (c *OrderController) ConfirmPayment(ctx *fiber.Ctx) error {
	orderID := ctx.Params("id")

	body := order_case.ConfirmPaymentInput{}

	if err := parser.ParseBody(ctx, &body); err != nil {
		return ctx.JSON(err)
	}

	if err := validator.V.Struct(body); err != nil {
		return err
	}

	ok, err := c.usecase.ConfirmPayment(ctx.UserContext(), order_case.ConfirmPaymentInput{
		OrderID:           orderID,
		ExternalPaymentID: body.ExternalPaymentID,
		Provider:          body.Provider,
		PaidAt:            body.PaidAt,
	})

	if err != nil {
		return err
	}

	if !ok {
		return ctx.SendStatus(404)
	}

	return ctx.SendStatus(200)
}
