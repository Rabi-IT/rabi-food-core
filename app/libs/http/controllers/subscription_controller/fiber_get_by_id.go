package subscription_controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// GetByID godoc
// @Summary Get subscription by ID
// @Description Get a subscription by ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /subscription/{id} [get].
func (c *SubscriptionController) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	data, err := c.usecase.GetByID(ctx.UserContext(), id)
	if err != nil {
		return err
	}

	if data == nil {
		return ctx.SendStatus(http.StatusNotFound)
	}

	return ctx.JSON(data)
}
