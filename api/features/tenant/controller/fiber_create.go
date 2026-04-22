package controller

import (
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
	"github.com/Rabi-IT/rabi-food-core/libs/validator"
	"github.com/gofiber/fiber/v2"
)

type createInput struct {
	Name string `json:"name" validate:"required"`
}

func (c *TenantController) Create(ctx *fiber.Ctx) error {
	data := createInput{}
	if err := ctx.BodyParser(&data); err != nil {
		return err
	}
	if err := validator.V.Struct(data); err != nil {
		return err
	}

	uctx := ctx.UserContext()
	logger.GetWideEvent(uctx).Event = "create-tenant"

	session := app_context.GetSession(uctx)
	id, err := c.usecase.Create(uctx, data.Name, session.UserID)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusCreated).SendString(id)
}
