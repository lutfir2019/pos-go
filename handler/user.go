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
	user            types.User
	reviewYourInput = "Review your input"
)

// GetAlluser query all users
func GetAllUsers(c *fiber.Ctx) error {
	page, limit, offset, err := helper.GetPaginationParams(c)
	if err == nil {
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, "Failed to get pagination params", err)
	}

	return services.GetAllUsers(c, 0, "", page, limit, offset)
}

// GetUser get a user
func GetUser(c *fiber.Ctx) error {
	db := database.DB

	userLogin, err := helper.ConvertJWT(c)
	if err != nil {
		return helper.HandleErrorResponse(c, fiber.StatusUnauthorized, "Invalid token", err)
	}

	err = db.Where(deletedNull).First(&user, model.User{Username: userLogin.Username}).Count(&total).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.HandleErrorResponse(c, fiber.StatusNotFound, err.Error(), err)
		}
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, err.Error(), err)
	}

	return helper.HandleSuccessResponse(c, fiber.StatusOK, "Get Product by username", user, 1, 10, total)
}

// CreateUser new user
func CreateUser(c *fiber.Ctx) error {
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
		Role:     input.Role,
		Password: hash,
	}

	err = db.Create(&user).Error
	if err != nil {
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, err.Error(), err)
	}

	return services.GetAllUsers(c, fiber.StatusCreated, "Success create user", 1, 10, 0)
}

// UpdateUser update user
func UpdateUser(c *fiber.Ctx) error {
	input := types.UpdateUser{}
	if err := c.BodyParser(&input); err != nil {
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, reviewYourInput, err)
	}

	username := c.Params("username")
	if username == "" {
		return CreateUser(c)
	}

	if !helper.ValidToken(c, username) {
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, "Invalid token", nil)
	}

	userModel, err := services.GetUserByUsername(username)
	if err != nil {
		return helper.HandleErrorResponse(c, fiber.StatusInternalServerError, "Internal Server Error", err)
	} else if userModel == nil {
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, "Invalid username or password", err)
	}

	// untuk change password
	if input.CurrentPassword != "" {
		return services.ChangePassword(c, &input, userModel, username)
	} else {
		return services.ChangeProfile(c, &input, username)
	}
}

// DeleteUser delete user
func DeleteUser(c *fiber.Ctx) error {
	username := c.Params("username")

	db := database.DB
	var user model.User

	userLogin, err := helper.ConvertJWT(c)
	if err != nil {
		return helper.HandleErrorResponse(c, fiber.StatusUnauthorized, "Invalid user", err)
	}

	if userLogin.Role != "admin" {
		return helper.HandleErrorResponse(c, fiber.StatusForbidden, "Access Forbidden", nil)
	}

	err = db.Where(deletedNull).Delete(&user, model.User{Username: username}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.HandleErrorResponse(c, fiber.StatusNotFound, err.Error(), err)
		}
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, err.Error(), err)
	}

	return services.GetAllUsers(c, fiber.StatusOK, "Success delete user", 1, 10, 0)
}
