package handler

import (
	"api-fiber-gorm/database"
	"api-fiber-gorm/helper"
	"api-fiber-gorm/model"
	"api-fiber-gorm/services"
	"api-fiber-gorm/types"
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var (
	total         int64
	orderNotFound = "Order not found"
	forbidden     = "Access Forbidden"
)

const deletedNull = "deleted_at IS NULL"

// GetAllProducts query all products
func GetAllProducts(c *fiber.Ctx) error {
	page, limit, offset, err := helper.GetPaginationParams(c)
	if err != nil {
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, "Failed to get pagination params", err)
	}

	return services.GetAllProducts(c, 0, "", page, limit, offset)
}

// GetProduct query product
func GetProduct(c *fiber.Ctx) error {
	code := c.Params("code")
	db := database.DB

	var product types.Products
	err := db.Model(&model.Product{}).
		Select("products.*, users.name as user_name").
		Joins("left join users on users.id = products.refer_user").
		Where(deletedNull).
		First(&product, model.Product{Code: code}).
		Count(&total).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.HandleErrorResponse(c, fiber.StatusNotFound, err.Error(), err)
		}
		return helper.HandleErrorResponse(c, fiber.StatusOK, err.Error(), err)
	}

	return helper.HandleSuccessResponse(c, fiber.StatusOK, "Get Product by ID", product, 1, 10, total)
}

// CreateProduct new product
func CreateProduct(c *fiber.Ctx) error {
	db := database.DB
	input := new(model.Product)
	if err := c.BodyParser(input); err != nil {
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, "Couldn't create product", err)
	}

	code, err := services.GenerateProductCode()
	if err != nil {
		return nil
	}

	user, err := helper.ConvertJWT(c)
	if err != nil {
		return helper.HandleErrorResponse(c, fiber.StatusUnauthorized, "Invalid token", err)
	}

	if user.Role != "admin" {
		return helper.HandleErrorResponse(c, fiber.StatusForbidden, forbidden, nil)
	}

	product := model.Product{
		Model:         gorm.Model{},
		Code:          code,
		Name:          input.Name,
		PricePurchase: input.PricePurchase,
		PriceSelling:  input.PriceSelling,
		Quantity:      input.Quantity,
		Image:         input.Image,
		ReferUser:     user.ID,
	}

	err = db.Create(&product).Error
	if err != nil {
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, err.Error(), err)
	}

	return services.GetAllProducts(c, fiber.StatusCreated, "Success create product", 1, 10, 0)
}

// UpdateUser update user
func UpdateProduct(c *fiber.Ctx) error {
	input := model.Product{}
	if err := c.BodyParser(&input); err != nil {
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, reviewYourInput, err)
	}

	code := c.Params("code")
	if code == "" {
		return CreateProduct(c)
	}

	userLogin, err := helper.ConvertJWT(c)
	if err != nil {
		return helper.HandleErrorResponse(c, fiber.StatusUnauthorized, "Invalid token", err)
	}

	if userLogin.Role != "admin" {
		return helper.HandleErrorResponse(c, fiber.StatusForbidden, forbidden, nil)
	}

	db := database.DB
	err = db.Model(&model.Product{}).Where(deletedNull).Where("code =?", code).Updates(input).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.HandleErrorResponse(c, fiber.StatusNotFound, err.Error(), err)
		}
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, err.Error(), err)
	}

	return services.GetAllProducts(c, fiber.StatusOK, "Success update product", 1, 10, 0)
}

// DeleteProduct delete product
func DeleteProduct(c *fiber.Ctx) error {
	code := c.Params("code")
	db := database.DB

	var product model.Product

	userLogin, err := helper.ConvertJWT(c)
	if err != nil {
		return helper.HandleErrorResponse(c, fiber.StatusUnauthorized, "Invalid user", err)
	}

	if userLogin.Role != "admin" {
		return helper.HandleErrorResponse(c, fiber.StatusForbidden, forbidden, nil)
	}

	err = db.Where(deletedNull).Delete(&product, model.Product{Code: code}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.HandleErrorResponse(c, fiber.StatusNotFound, err.Error(), err)
		}
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, err.Error(), err)
	}

	return services.GetAllProducts(c, 0, "Success delete product", 1, 10, 0)
}
