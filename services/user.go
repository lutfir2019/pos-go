package services

import (
	"api-fiber-gorm/database"
	"api-fiber-gorm/helper"
	"api-fiber-gorm/model"
	"api-fiber-gorm/types"

	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var (
	user model.User
)

const deletedNull = "deleted_at IS NULL"

func ValidUser(username string, p string) bool {
	db := database.DB

	db.First(&user, username)
	if user.Username == "" {
		return false
	}
	if !helper.CheckPasswordHash(p, user.Password) {
		return false
	}
	return true
}

func GetUserByUsername(e string) (*model.User, error) {
	db := database.DB

	if err := db.Where(deletedNull).First(&user, model.User{Username: e}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func GetAllUsers(c *fiber.Ctx, sts int, msg string, page int, limit int, offset int) error {
	db := database.DB
	var users []types.User

	message := "Get all users"
	if msg != "" {
		message = msg
	}

	status := fiber.StatusOK
	if sts > 0 {
		status = sts
	}

	err := db.Where(deletedNull).Order("ID DESC").Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.HandleSuccessResponse(c, fiber.StatusOK, message, users, page, limit, total)
		}
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, err.Error(), err)
	}

	// Get total count of users
	err = db.Model(&model.User{}).Where(deletedNull).Count(&total).Error
	if err != nil {
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, err.Error(), err)
	}

	return helper.HandleSuccessResponse(c, status, message, users, page, limit, total)
}

func ChangePassword(c *fiber.Ctx, input *types.UpdateUser, userModel *model.User, username string) error {
	if !helper.CheckPasswordHash(input.CurrentPassword, userModel.Password) {
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, "Password is wrong", nil)
	}

	if helper.CheckPasswordHash(input.Password, userModel.Password) {
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, "Password baru tidak boleh sama dengan password sebelumnya", nil)
	}

	if input.Password != input.ConfirmPassword {
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, "Password not match", nil)
	}

	hash, err := helper.HashPassword(input.Password)
	if err != nil {
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, "Couldn't hash password", err)
	}

	newDataUser := model.User{
		Name:     input.Name,
		Password: hash,
	}

	db := database.DB
	err = db.Model(&model.User{}).Where(deletedNull).Where("username =?", username).Updates(&newDataUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.HandleErrorResponse(c, fiber.StatusNotFound, err.Error(), err)
		}
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, err.Error(), err)
	}

	db.First(&user, model.User{Username: newDataUser.Username})

	return helper.HandleSuccessResponse(c, fiber.StatusOK, "Success update user", user, 1, 10, 0)
}

func ChangeProfile(c *fiber.Ctx, input *types.UpdateUser, username string) error {

	newDataUser := model.User{
		Name: input.Name,
	}

	db := database.DB
	err := db.Model(&model.User{}).Where(deletedNull).Where("username =?", username).Updates(&newDataUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.HandleErrorResponse(c, fiber.StatusNotFound, err.Error(), err)
		}
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, err.Error(), err)
	}

	db.First(&user, model.User{Username: newDataUser.Username})

	return helper.HandleSuccessResponse(c, fiber.StatusOK, "Success update user", user, 1, 10, 0)
}
