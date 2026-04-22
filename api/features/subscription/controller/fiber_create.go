package controller

import (
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	"github.com/Rabi-IT/rabi-food-core/features/subscription/usecases"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
	"github.com/Rabi-IT/rabi-food-core/libs/validator"

	"github.com/gofiber/fiber/v2"
)

// Create godoc
// @Summary Create subscription
// @Description Create a new subscription
// @Tags subscriptions
// @Accept json
// @Produce text/plain
// @Param subscription body usecases.CreateInput true "Subscription data"
// @Success 201 {string} string "Created subscription ID"
// @Failure 400 {object} middlewares.ValidationErrorResponse "Validation errors"
// @Failure 500 {string} string "Internal server error"
// @Router /subscription/ [post].
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
	id, err := c.usecase.Create(uctx, data)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusCreated).SendString(id)
}
