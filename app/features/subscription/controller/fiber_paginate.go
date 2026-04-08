package controller

import (
	"strconv"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	g "github.com/Rabi-IT/rabi-food-core/features/subscription/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/database"

	"github.com/gofiber/fiber/v2"
)

// Paginate godoc
// @Summary Paginate subscriptions
// @Description Paginate subscriptions using query filters
// @Tags subscriptions
// @Produce json
// @Param Page query int false "Page number"
// @Param PageSize query int false "Page size"
// @Success 200 {object} g.PaginateOutput
// @Failure 500 {string} string "Internal server error"
// @Router /subscription/ [get].
func (c *SubscriptionController) Paginate(ctx *fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("Page", "0"))
	if err != nil {
		return err
	}

	pageSize, err := strconv.Atoi(ctx.Query("PageSize", "10"))
	if err != nil {
		return err
	}

	uctx := ctx.UserContext()
	session := app_context.GetSession(uctx)

	filter := g.PaginateFilter{
		TenantID: &session.TenantID,
	}

	paginate := database.PaginateInput{
		Page:     page,
		PageSize: pageSize,
	}

	result, err := c.usecase.Paginate(uctx, filter, paginate)
	if err != nil {
		return err
	}

	return ctx.JSON(result)
}
