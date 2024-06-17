package handler

import (
	"strconv"

	"api-fiber-gorm/database"
	"api-fiber-gorm/model"
	"api-fiber-gorm/types"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	reviewYourInput = "Review your input"
	users           []model.User
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func validToken(t *jwt.Token, id string) bool {
	n, err := strconv.Atoi(id)
	if err != nil {
		return false
	}

	claims := t.Claims.(jwt.MapClaims)
	uid := int(claims["user_id"].(float64))

	return uid == n
}

func validUser(id string, p string) bool {
	db := database.DB
	var user model.User
	db.First(&user, id)
	if user.Username == "" {
		return false
	}
	if !CheckPasswordHash(p, user.Password) {
		return false
	}
	return true
}

// GetUser get a user
func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var user model.User
	db.Find(&user, id)
	if user.Username == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "No user found with ID", "data": nil})
	}
	showUser := types.UserData{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Username: user.Username,
	}
	return c.JSON(fiber.Map{"status": "success", "message": "User found", "data": showUser})
}

// GetAlluser query all users
func GetAllUsers(c *fiber.Ctx) error {
	db := database.DB
	db.Find(&users)
	return c.JSON(fiber.Map{"status": "success", "message": "All users", "data": users})
}

// CreateUser new user
func CreateUser(c *fiber.Ctx) error {
	db := database.DB
	input := new(types.UserData)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": reviewYourInput, "data": err})

	}

	hash, err := hashPassword(input.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Couldn't hash password", "data": err})

	}

	user := model.User{
		Model:    gorm.Model{},
		Name:     input.Name,
		Email:    input.Email,
		Username: input.Username,
		Password: hash,
	}
	if err := db.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Couldn't create user", "data": err})
	}

	db.Find(&users)
	return c.JSON(fiber.Map{"status": "success", "message": "Created user", "data": users})
}

// UpdateUser update user
func UpdateUser(c *fiber.Ctx) error {
	input := types.UserData{}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": reviewYourInput, "data": err})
	}
	id := c.Params("id")
	if id == "" {
		return CreateUser(c)
	}

	token := c.Locals("user").(*jwt.Token)

	if !validToken(token, id) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "data": nil})
	}

	db := database.DB
	err := db.Model(&model.User{}).Where("id =?", id).Updates(input).Error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err, "data": nil})
	}
	db.Find(&users)
	return c.JSON(fiber.Map{"status": "success", "message": "User successfully updated", "data": users})
}

// DeleteUser delete user
func DeleteUser(c *fiber.Ctx) error {
	type PasswordInput struct {
		Password string `json:"password"`
	}
	var pi PasswordInput
	if err := c.BodyParser(&pi); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": reviewYourInput, "data": err})
	}
	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	if !validToken(token, id) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "data": nil})

	}

	if !validUser(id, pi.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Not valid user", "data": nil})

	}

	db := database.DB
	var user model.User

	db.First(&user, id)
	if user.Email == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "No user found with ID", "data": nil})
	}

	db.Delete(&user)
	db.Find(&users)
	return c.JSON(fiber.Map{"status": "success", "message": "User successfully deleted", "data": users})
}
