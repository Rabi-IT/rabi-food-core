package controller

import (
	"net/http"

	gateway "github.com/Rabi-IT/rabi-food-core/features/subscription/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/validator"

	"github.com/gofiber/fiber/v2"
)

// UpsertConfig godoc
// @Summary Upsert subscription config
// @Description Create or update subscription config for the current tenant
// @Tags subscriptions
// @Accept json
// @Produce text/plain
// @Param config body gateway.UpsertConfigInput true "Subscription config"
// @Success 200 {string} string "Updated"
// @Failure 400 {object} middlewares.ValidationErrorResponse "Validation errors"
// @Failure 500 {string} string "Internal server error"
// @Router /subscription/config [put].
func (c *SubscriptionController) UpsertConfig(ctx *fiber.Ctx) error {
	data := gateway.UpsertConfigInput{}
	err := ctx.BodyParser(&data)
	if err != nil {
		return ctx.JSON(err)
	}

	err = validator.V.Struct(data)
	if err != nil {
		return err
	}

	err = c.usecase.UpsertConfig(ctx.UserContext(), data)
	if err != nil {
		return err
	}

	return ctx.SendStatus(http.StatusOK)
}
