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
// @Description Create a new subscription group, one subscription per selected product
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body usecases.CreateInput true "Subscription data"
// @Success 201 {array} string "Created subscription IDs"
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
	groupID, err := c.usecase.Create(uctx, data)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusCreated).SendString(groupID)
}
