package handler

import (
	"time"

	"api-fiber-gorm/config"
	"api-fiber-gorm/database"
	"api-fiber-gorm/helper"
	"api-fiber-gorm/model"
	"api-fiber-gorm/services"
	"api-fiber-gorm/types"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// Login get user and password
func Login(c *fiber.Ctx) error {
	input := new(types.LoginInput)
	if err := c.BodyParser(&input); err != nil {
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, "Error on login request", err)
	}

	Username := input.Username
	pass := input.Password
	userModel, err := new(model.User), *new(error)

	userModel, err = services.GetUserByUsername(Username)

	if err != nil {
		return helper.HandleErrorResponse(c, fiber.StatusInternalServerError, "Internal Server Error", err)
	} else if userModel == nil {
		helper.CheckPasswordHash(pass, "")
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, "Invalid username or password", err)
	}

	if !helper.CheckPasswordHash(pass, userModel.Password) {
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, "Invalid username or password", nil)
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = userModel.Username
	claims["role"] = userModel.Role
	claims["user_id"] = userModel.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(config.Config("SECRET_KEY")))
	if err != nil {
		return helper.HandleErrorResponse(c, fiber.StatusInternalServerError, err.Error(), err)
	}

	user := types.User{
		ID:       userModel.ID,
		Username: userModel.Username,
		Name:     userModel.Name,
		Role:     userModel.Role,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Success login", "token": t, "data": user})
}

// CreateUser new user
func Register(c *fiber.Ctx) error {
	db := database.DB
	input := new(types.Register)
	if err := c.BodyParser(input); err != nil {
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, reviewYourInput, err)
	}

	if input.Password != input.ConfirmPassword {
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, "Password not match", nil)
	}

	hash, err := helper.HashPassword(input.Password)
	if err != nil {
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, "Couldn't hash password", err)
	}

	user := model.User{
		Model:    gorm.Model{},
		Name:     input.Name,
		Username: input.Username,
		Role:     "user",
		Password: hash,
	}

	if err := db.Create(&user).Error; err != nil {
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, err.Error(), err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Register successfully", "data": nil})
}

func CreatAdmin(c *fiber.Ctx) error {
	db := database.DB
	input := new(types.Register)
	if err := c.BodyParser(input); err != nil {
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, reviewYourInput, err)
	}

	if input.Password != input.ConfirmPassword {
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, "Password not match", nil)
	}

	hash, err := helper.HashPassword(input.Password)
	if err != nil {
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, "Couldn't hash password", err)
	}

	user := model.User{
		Model:    gorm.Model{},
		Name:     input.Name,
		Username: input.Username,
		Role:     "admin",
		Password: hash,
	}

	if err := db.Create(&user).Error; err != nil {
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, err.Error(), err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Register successfully", "data": nil})
}
