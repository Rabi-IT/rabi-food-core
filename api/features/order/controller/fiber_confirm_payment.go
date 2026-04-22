package controller

import (
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/features/order/usecases"
	parser "github.com/Rabi-IT/rabi-food-core/libs/http/parser"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
	"github.com/Rabi-IT/rabi-food-core/libs/validator"

	"github.com/gofiber/fiber/v2"
)

// ConfirmPayment godoc
// @Summary Confirm order payment
// @Description Confirm payment for a specific order
// @Tags orders
// @Accept json
// @Produce text/plain
// @Param id path string true "Order ID"
// @Param payment body usecases.ConfirmPaymentInput true "Payment confirmation data"
// @Success 200 {string} string "Payment confirmed"
// @Failure 400 {object} middlewares.ValidationErrorResponse "Validation errors"
// @Failure 404 {string} string "Not found"
// @Failure 409 {object} errs.ApiError "Conflict"
// @Failure 500 {string} string "Internal server error"
// @Router /order/{id}/payments/confirm [post].
func (c *OrderController) ConfirmPayment(ctx *fiber.Ctx) error {
	orderID := ctx.Params("id")

	body := usecases.ConfirmPaymentInput{}

	err := parser.ParseBody(ctx, &body)
	if err != nil {
		return ctx.JSON(err)
	}

	err = validator.V.Struct(body)
	if err != nil {
		return err
	}

	uctx := ctx.UserContext()
	logger.GetWideEvent(uctx).Event = "confirm-order-payment"
	ok, err := c.usecase.ConfirmPayment(uctx, usecases.ConfirmPaymentInput{
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
