package order_controller

import (
	"net/http"
	"rabi-food-core/libs/http/fiber_adapter/parser"
	"rabi-food-core/libs/validator"
	"rabi-food-core/usecases/order_case"

	"github.com/gofiber/fiber/v2"
)

// ConfirmPayment godoc
// @Summary Confirm order payment
// @Description Confirm payment for a specific order
// @Tags orders
// @Accept json
// @Produce text/plain
// @Param id path string true "Order ID"
// @Param payment body order_case.ConfirmPaymentInput true "Payment confirmation data"
// @Success 200 {string} string "Payment confirmed"
// @Failure 400 {object} middlewares.ValidationErrorResponse "Validation errors"
// @Failure 404 {string} string "Not found"
// @Failure 409 {object} errs.AppError "Conflict"
// @Failure 500 {string} string "Internal server error"
// @Router /order/{id}/payments/confirm [post].
func (c *OrderController) ConfirmPayment(ctx *fiber.Ctx) error {
	orderID := ctx.Params("id")

	body := order_case.ConfirmPaymentInput{}

	err := parser.ParseBody(ctx, &body)
	if err != nil {
		return ctx.JSON(err)
	}

	err = validator.V.Struct(body)
	if err != nil {
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
		return ctx.SendStatus(http.StatusNotFound)
	}

	return ctx.SendStatus(http.StatusOK)
}
