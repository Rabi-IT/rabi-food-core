package controller

import (
	"strconv"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	g "github.com/Rabi-IT/rabi-food-core/features/user/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/database"

	"github.com/gofiber/fiber/v2"
)

// Paginate godoc
// @Summary Paginate users
// @Description Paginate users using query filters
// @Tags users
// @Produce json
// @Param Page query int false "Page number"
// @Param PageSize query int false "Page size"
// @Param tenantId query string false "Tenant ID"
// @Param state query string false "State"
// @Param city query string false "City"
// @Param name query string false "Name"
// @Success 200 {object} g.PaginateOutput
// @Failure 500 {string} string "Internal server error"
// @Router /user/ [get].
func (c *UserController) Paginate(ctx *fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("Page", "0"))
	if err != nil {
		return err
	}

	pageSize, err := strconv.Atoi(ctx.Query("PageSize", "10"))
	if err != nil {
		return err
	}

	filter := g.PaginateFilter{}
	err = ctx.QueryParser(&filter)
	if err != nil {
		return err
	}

	uctx := ctx.UserContext()
	session := app_context.GetSession(uctx)
	filter.TenantID = &session.TenantID

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
