package adminhandler

import (
	"strconv"
	"subscription_manager/data"
	"subscription_manager/database"

	"github.com/gofiber/fiber/v2"
)

func UpdateTicketStatus(c *fiber.Ctx) error {
	userEmail, ok := c.Locals("userEmail").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ticketID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	var req data.TicketStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	var updated data.TicketResponse
	err = database.InitiateDataBase().QueryRow(c.Context(), `
        UPDATE tickets
        SET status = $1,
            updated_at = NOW(),
            updated_by = $2
        WHERE id = $3 AND deleted_at IS NULL
        RETURNING id, user_id, title, description, COALESCE(link,''), priority, status, created_by, created_at, updated_by, updated_at
    `,
		req.Status,
		userEmail,
		ticketID,
	).Scan(
		&updated.ID,
		&updated.UserID,
		&updated.Title,
		&updated.Description,
		&updated.Link,
		&updated.Priority,
		&updated.Status,
		&updated.CreatedBy,
		&updated.CreatedAt,
		&updated.UpdatedBy,
		&updated.UpdatedAt,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ticket status updated successfully",
		"ticket":  updated,
	})
}
