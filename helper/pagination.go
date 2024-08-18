package helper

import (
	"github.com/gofiber/fiber/v2"
)

// GetPaginationParams gets page and limit from query params
func GetPaginationParams(c *fiber.Ctx) (int, int, int, error) {
	// Default values
	page := 1
	limit := 10

	// Get page from query param
	pageQuery := c.QueryInt("page")
	if pageQuery > 0 {
		page = pageQuery
	}

	// Get limit from query param
	limitQuery := c.QueryInt("per_page")
	if limitQuery > 0 {
		limit = limitQuery
	}

	// Calculate offset
	offset := (page - 1) * limit
	return page, limit, offset, nil
}

// page, limit, offset, err := helper.GetPaginationParams(c)
// db.Offset(offset).Limit(limit).Find(&products)
