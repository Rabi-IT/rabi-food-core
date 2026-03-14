package subscription_controller

import (
	"net/http"
	"rabi-food-core/libs/validator"
	"rabi-food-core/usecases/subscription_case"

	"github.com/gofiber/fiber/v2"
)

// CreateSubscription godoc
// @Summary Create subscription
// @Description Create a new subscription
// @Tags subscriptions
// @Accept json
// @Produce text/plain
// @Param subscription body subscription_case.CreateInput true "Subscription data"
// @Success 201 {string} string "Created subscription ID"
// @Failure 400 {object} middlewares.ValidationErrorResponse "Validation errors"
// @Failure 500 {string} string "Internal server error"
// @Router /subscription/ [post]
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
