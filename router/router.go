package router

import (
	"api-fiber-gorm/handler"
	"api-fiber-gorm/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {
	// Middleware
	app.Use(middleware.Security)

	api := app.Group("/api", logger.New())
	api.Get("/", handler.Hello)

	// Auth
	auth := api.Group("/auth")
	auth.Post("/login", handler.Login)

	// User
	user := api.Group("/user")
	user.Get("/", handler.GetAllUsers)
	user.Get("/:id", handler.GetUser)
	user.Post("/", handler.CreateUser)
	user.Put("/:id", middleware.Protected(), handler.UpdateUser)
	user.Delete("/:id", middleware.Protected(), handler.DeleteUser)

	// Product
	product := api.Group("/product")
	product.Get("/", handler.GetAllProducts)
	product.Get("/:id", handler.GetProduct)
	product.Post("/", middleware.Protected(), handler.CreateProduct)
	product.Delete("/:id", middleware.Protected(), handler.DeleteProduct)

	// Order
	order := api.Group("/order")
	order.Get("/", handler.GetOrders)
	order.Post("/", handler.CreateOrder)
	order.Get("/:id", handler.GetOrderByID)
	order.Put("/:id", handler.UpdateOrder)
	order.Delete("/:id", handler.DeleteOrder)
}
