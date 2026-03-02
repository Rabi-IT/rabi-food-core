package app_plan_controller

import (
	"net/http"
	"rabi-food-core/libs/database/gateways/app_plan_gateway"
	"rabi-food-core/libs/http/fiber_adapter/parser"

	"github.com/gofiber/fiber/v2"
)

func (c *AppPlanController) Delete(ctx *fiber.Ctx) error {
	filter := app_plan_gateway.DeleteFilter{}

	err := parser.ParseBody(ctx, &filter)
	if err != nil {
		return ctx.JSON(err)
	}

	filter.ID = ctx.Params("id")
	deleted, err := c.usecase.Delete(ctx.UserContext(), filter)

	if err != nil {
		return err
	}

	if deleted {
		return ctx.SendStatus(http.StatusNoContent)
	}

	return ctx.SendStatus(http.StatusNotFound)
}
