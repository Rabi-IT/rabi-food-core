package controller

import (
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	"github.com/Rabi-IT/rabi-food-core/features/subscription/usecases"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
	"github.com/Rabi-IT/rabi-food-core/libs/validator"

	"github.com/gofiber/fiber/v2"
)

func (c *SubscriptionController) Create(ctx *fiber.Ctx) error {
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
	session := app_context.GetSession(uctx)
	data.TenantID = ctx.Get("X-Tenant-ID")
	data.UserID = session.UserID

	logger.GetWideEvent(uctx).Event = "create-subscription"
	ids, err := c.usecase.Create(uctx, data)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusCreated).JSON(ids)
}
