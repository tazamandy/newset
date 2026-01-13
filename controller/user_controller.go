package controller

import (
	"attendance-system/models"

 	"github.com/gofiber/fiber/v2"
)

// GetMyQRCode returns the authenticated user's QR code data
func GetMyQRCode(c *fiber.Ctx) error {
    user, ok := c.Locals("user").(models.User)
    if !ok {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
    }

    if user.QRCodeData == "" {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "QR code not found"})
    }

    return c.JSON(fiber.Map{"qr_code": user.QRCodeData})
}
