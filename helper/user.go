package helper

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Define a struct for the return value
type UserData struct {
	Username string
	Role     string
	ID       uint
}

// ConvertJWT extracts and converts JWT claims to a custom UserData struct
func ConvertJWT(c *fiber.Ctx) (*UserData, error) {
	// Ambil nilai dari Locals
	user, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token data")
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok || !user.Valid {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token claims")
	}

	// Ambil user_id dan Username dari klaim
	cUsername, ok := claims["username"].(string)
	if !ok {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid Username claim")
	}

	// Ambil user_id sebagai float64
	cIdFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid user_id claim")
	}

	// Ambil role sebagai string
	cRole, ok := claims["role"].(string)
	if !ok {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid role claim")
	}

	// Konversi float64 ke uint
	referId := uint(cIdFloat)

	// Return the UserData struct
	return &UserData{
		Username: cUsername,
		Role:     cRole,
		ID:       referId,
	}, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func ValidToken(c *fiber.Ctx, username string) bool {
	user, err := ConvertJWT(c)
	if err != nil {
		return false
	}

	return user.Username == username
}
