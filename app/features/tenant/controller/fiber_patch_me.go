package controller

import (
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	"github.com/Rabi-IT/rabi-food-core/features/tenant/usecases"
	"github.com/Rabi-IT/rabi-food-core/libs/http/parser"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"

	"github.com/gofiber/fiber/v2"
)

func (c *TenantController) PatchMe(ctx *fiber.Ctx) error {
	uctx := ctx.UserContext()
	id := app_context.GetSession(uctx).TenantID
	filter := usecases.PatchFilter{
		ID: &id,
	}

	data := usecases.PatchValues{}
	err := parser.ParseBody(ctx, &data)
	if err != nil {
		return ctx.JSON(err)
	}

	logger.GetWideEvent(uctx).Event = "patch-tenant"
	updated, err := c.usecase.Patch(uctx, filter, data)

	if err != nil {
		return err
	}

	if updated {
		return ctx.SendStatus(http.StatusOK)
	}

	return ctx.SendStatus(http.StatusNotFound)
}
