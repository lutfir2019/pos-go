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
	auth.Post("/register", handler.Register)
	auth.Post("/danger-secret-side-create-admin", handler.CreatAdmin)

	// User
	user := api.Group("/user")
	user.Post("/create-user", middleware.Protected(), handler.CreateUser)
	user.Get("/get-users", middleware.Protected(), handler.GetAllUsers)
	user.Get("/get-user", middleware.Protected(), handler.GetUser)
	user.Put("/:username/update-user", middleware.Protected(), handler.UpdateUser)
	user.Delete("/:username/delete-user", middleware.Protected(), handler.DeleteUser)

	// Product
	product := api.Group("/product")
	product.Post("/create-product", middleware.Protected(), handler.CreateProduct)
	product.Get("/get-products", middleware.Protected(), handler.GetAllProducts)
	product.Get("/:code/get-product", middleware.Protected(), handler.GetProduct)
	product.Put("/:code/update-product", middleware.Protected(), handler.UpdateProduct)
	product.Delete("/:code/delete-product", middleware.Protected(), handler.DeleteProduct)

	// Order
	order := api.Group("/order")
	order.Get("/", middleware.Protected(), handler.GetOrders)
	order.Post("/", middleware.Protected(), handler.CreateOrder)
	order.Get("/:id", middleware.Protected(), handler.GetOrderByID)
	order.Put("/:id", middleware.Protected(), handler.UpdateOrder)
	order.Delete("/:id", middleware.Protected(), handler.DeleteOrder)
}
