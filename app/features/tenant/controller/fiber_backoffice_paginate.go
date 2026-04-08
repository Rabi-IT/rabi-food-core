package controller

import (
	"strconv"

	g "github.com/Rabi-IT/rabi-food-core/features/tenant/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/database"

	"github.com/gofiber/fiber/v2"
)

// Paginate godoc
// @Summary Paginate tenants (backoffice)
// @Description Paginate all tenants
// @Tags tenants backoffice
// @Produce json
// @Param Page query int false "Page number"
// @Param PageSize query int false "Page size"
// @Param name query string false "Tenant name"
// @Success 200 {object} g.PaginateOutput
// @Failure 500 {string} string "Internal server error"
// @Router /backoffice/tenant/ [get].
func (c *TenantBackofficeController) Paginate(ctx *fiber.Ctx) error {
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
