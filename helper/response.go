package helper

import (
	"api-fiber-gorm/types"

	"github.com/gofiber/fiber/v2"
)

// HandleSuccessResponse sends a success response with optional pagination
func HandleSuccessResponse(c *fiber.Ctx, status int, message string, data interface{}, page int, limit int, total int64) error {
	response := types.SuccessResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}

	if page > 0 && limit > 0 {
		totalPages := int((total + int64(limit) - 1) / int64(limit))
		response.Meta.Pagination.Page = page
		response.Meta.Pagination.Limit = limit
		response.Meta.Pagination.Total = total
		response.Meta.Pagination.TotalPages = totalPages
	}

	return c.Status(status).JSON(response)
}

// HandleErrorResponse sends an error response
func HandleErrorResponse(c *fiber.Ctx, status int, message string, err error) error {
	response := types.ErrorResponse{
		Status:  "error",
		Message: message,
	}
	if err != nil {
		response.Error = err.Error()
	}

	return c.Status(status).JSON(response)
}
