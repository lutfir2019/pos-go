package services

import (
	"api-fiber-gorm/database"
	"api-fiber-gorm/helper"
	"api-fiber-gorm/model"
	"api-fiber-gorm/types"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var (
	total int64
)

func GetAllProducts(c *fiber.Ctx, sts int, msg string, page int, limit int, offset int) error {
	db := database.DB
	var products []types.Products

	message := "Get all products"
	if msg != "" {
		message = msg
	}

	status := fiber.StatusOK
	if sts > 0 {
		status = sts
	}

	err := db.Model(&model.Product{}).
		Select("products.*, users.name as refer_user").
		Joins("left join users on users.id = products.refer_user").
		Where("products.deleted_at IS NULL").Order("ID DESC").
		Offset(offset).
		Limit(limit).
		Scan(&products).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.HandleSuccessResponse(c, fiber.StatusOK, "Product not found", products, page, limit, total)
		}
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, err.Error(), err)
	}

	// Get total count of products
	err = db.Model(&model.Product{}).Where(deletedNull).Count(&total).Error
	if err != nil {
		return helper.HandleErrorResponse(c, fiber.StatusBadRequest, err.Error(), err)
	}

	return helper.HandleSuccessResponse(c, status, message, products, page, limit, total)
}

func GenerateProductCode() (string, error) {
	db := database.DB
	now := time.Now()
	datePrefix := now.Format("20060102")

	var counter model.ProductCounter

	err := db.Where("date = ?", datePrefix).First(&counter).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}

	if err == gorm.ErrRecordNotFound {
		// Insert new record if not exists
		counter = model.ProductCounter{
			Date:    datePrefix,
			Counter: 1,
		}
		err = db.Create(&counter).Error
		if err != nil {
			return "", err
		}
	} else {
		// Update existing record
		counter.Counter++
		err = db.Save(&counter).Error
		if err != nil {
			return "", err
		}
	}

	code := fmt.Sprintf("PR%s%04d", datePrefix, counter.Counter)
	return code, nil
}
