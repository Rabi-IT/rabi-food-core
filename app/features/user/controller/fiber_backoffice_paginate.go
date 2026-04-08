package controller

import (
	"strconv"

	g "github.com/Rabi-IT/rabi-food-core/features/user/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/database"

	"github.com/gofiber/fiber/v2"
)

// Paginate godoc
// @Summary Paginate users (backoffice)
// @Description Paginate users without tenant restriction; TenantID filter is optional
// @Tags users backoffice
// @Produce json
// @Param Page query int false "Page number"
// @Param PageSize query int false "Page size"
// @Param tenantId query string false "Optional tenant ID filter"
// @Param state query string false "State"
// @Param city query string false "City"
// @Param name query string false "Name"
// @Success 200 {object} g.PaginateOutput
// @Failure 500 {string} string "Internal server error"
// @Router /backoffice/user/ [get].
func (c *UserBackofficeController) Paginate(ctx *fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("Page", "0"))
	if err != nil {
		return err
	}

	pageSize, err := strconv.Atoi(ctx.Query("PageSize", "10"))
	if err != nil {
		return err
	}

	filter := g.PaginateFilter{}
	if err = ctx.QueryParser(&filter); err != nil {
		return err
	}

	uctx := ctx.UserContext()

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
