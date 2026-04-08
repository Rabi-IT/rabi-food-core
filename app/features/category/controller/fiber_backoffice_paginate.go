package controller

import (
	"strconv"

	"github.com/Rabi-IT/rabi-food-core/features/category/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/database"

	"github.com/gofiber/fiber/v2"
)

// Paginate godoc
// @Summary Paginate categories (backoffice)
// @Description Paginate categories without tenant restriction; TenantID filter is optional
// @Tags categories backoffice
// @Produce json
// @Param Page query int false "Page number"
// @Param PageSize query int false "Page size"
// @Param tenantId query string false "Optional tenant ID filter"
// @Param name query string false "Category name"
// @Success 200 {object} gateway.PaginateOutput
// @Failure 500 {string} string "Internal server error"
// @Router /backoffice/category/ [get].
func (c *CategoryController) BackofficePaginate(ctx *fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("Page", "0"))
	if err != nil {
		return err
	}

	pageSize, err := strconv.Atoi(ctx.Query("PageSize", "10"))
	if err != nil {
		return err
	}

	filter := gateway.PaginateFilter{}
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
