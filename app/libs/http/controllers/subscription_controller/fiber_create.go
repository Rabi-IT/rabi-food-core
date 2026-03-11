package subscription_controller

import (
	"net/http"
	"rabi-food-core/libs/validator"
	"rabi-food-core/usecases/subscription_case"

	"github.com/gofiber/fiber/v2"
)

func (c *SubscriptionController) Create(ctx *fiber.Ctx) error {
	data := subscription_case.CreateInput{}
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
