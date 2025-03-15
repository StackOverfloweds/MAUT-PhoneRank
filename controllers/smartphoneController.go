package controllers

import (
	"github.com/StackOverfloweds/MAUT-PhoneRank/database"
	"github.com/StackOverfloweds/MAUT-PhoneRank/models"
	"github.com/gofiber/fiber/v2"
)

func GetSmartphone(c *fiber.Ctx) error {
	var smartphone []models.Smartphone
	database.DB.Find(&smartphone)
	return c.JSON(smartphone)
}

func CreateSmartphone(c *fiber.Ctx) error {
	var smartphone models.Smartphone
	if err := c.BodyParser(&smartphone); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	database.DB.Create(&smartphone)
	return c.JSON(smartphone)
}
