package subscription_controller

import (
	"net/http"
	"rabi-food-core/libs/database/gateways/subscription_gateway"
	"rabi-food-core/libs/validator"

	"github.com/gofiber/fiber/v2"
)

func (c *SubscriptionController) UpsertConfig(ctx *fiber.Ctx) error {
	data := subscription_gateway.UpsertConfigInput{}
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
