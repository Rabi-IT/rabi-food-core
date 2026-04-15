package controller

import (
	"strconv"

	"github.com/Rabi-IT/rabi-food-core/features/category/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/database"

	"github.com/gofiber/fiber/v2"
)

// Paginate godoc
// @Summary Paginate categories
// @Description Paginate categories using query filters
// @Tags categories
// @Produce json
// @Param Page query int false "Page number"
// @Param PageSize query int false "Page size"
// @Param name query string false "Category name"
// @Success 200 {object} gateway.PaginateOutput
// @Failure 500 {string} string "Internal server error"
// @Router /category/ [get].
func (c *CategoryController) Paginate(ctx *fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("Page", "0"))
	if err != nil {
		return err
	}

	pageSize, err := strconv.Atoi(ctx.Query("PageSize", "10"))
	if err != nil {
		return err
	}

	filter := gateway.PaginateFilter{}
	err = ctx.QueryParser(&filter)
	if err != nil {
		return err
	}

	tenantID := ctx.Get("X-Tenant-ID")
	if tenantID != "" {
		filter.TenantID = &tenantID
	}

	paginate := database.PaginateInput{
		Page:     page,
		PageSize: pageSize,
	}

	result, err := c.usecase.Paginate(ctx.UserContext(), filter, paginate)
	if err != nil {
		return err
	}

	return ctx.JSON(result)
}
