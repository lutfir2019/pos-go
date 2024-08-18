package handler

import (
	"api-fiber-gorm/database"
	"api-fiber-gorm/helper"
	"api-fiber-gorm/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetOrders mengembalikan daftar pesanan
func GetOrders(c *fiber.Ctx) error {
	var orders []model.Order
	database.DB.Preload("Items").Find(&orders)
	return c.JSON(fiber.Map{"status": "success", "message": "All orders", "data": orders})
}

// CreateOrder membuat pesanan baru
func CreateOrder(c *fiber.Ctx) error {
	input := new(model.Order)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := helper.ConvertJWT(c)
	if err != nil {
		return helper.HandleErrorResponse(c, fiber.StatusUnauthorized, "Invalid token", err)
	}

	order := model.Order{
		Model:  gorm.Model{},
		UserID: user.ID,
		Total:  input.Total,
		Items:  input.Items,
	}

	for _, item := range input.Items {
		var prd model.Product
		database.DB.First(&prd, item.ProductID)
		database.DB.Model(&model.Product{}).Where("id =?", item.ProductID).Updates(model.Product{Quantity: prd.Quantity - uint(item.Quantity)})
	}

	database.DB.Create(&order)
	return c.JSON(fiber.Map{"status": "success", "message": "Success create order", "data": order})
}

// GetOrderByID mengembalikan pesanan berdasarkan ID
func GetOrderByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var order model.Order
	if result := database.DB.Preload("Items").First(&order, id); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": orderNotFound})
	}
	return c.JSON(order)
}

// UpdateOrder memperbarui pesanan berdasarkan ID
func UpdateOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	var order model.Order
	if result := database.DB.First(&order, id); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": orderNotFound})
	}
	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	database.DB.Save(&order)
	return c.JSON(order)
}

// DeleteOrder menghapus pesanan berdasarkan ID
func DeleteOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	if result := database.DB.Delete(&model.Order{}, id); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": orderNotFound})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
