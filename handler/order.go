package handler

import (
	"api-fiber-gorm/database"
	"api-fiber-gorm/model"

	"github.com/gofiber/fiber/v2"
)

var (
	orderNotFound = "Order not found"
)

// GetOrders mengembalikan daftar pesanan
func GetOrders(c *fiber.Ctx) error {
	var orders []model.Order
	database.DB.Preload("Items").Find(&orders)
	return c.JSON(fiber.Map{"status": "success", "message": "All orders", "data": orders})
}

// CreateOrder membuat pesanan baru
func CreateOrder(c *fiber.Ctx) error {
	order := new(model.Order)
	if err := c.BodyParser(order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
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
