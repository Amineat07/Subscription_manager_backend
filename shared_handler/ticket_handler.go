package sharedhandler

import (
	"fmt"
	"strconv"
	"subscription_manager/data"
	"subscription_manager/database"
	"subscription_manager/utils"

	"github.com/gofiber/fiber/v2"
)

func ReplyToTicket(c *fiber.Ctx) error {
	var req_ticket_reply data.TicketReplyRequest
	if err := c.BodyParser(&req_ticket_reply); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}
	if err := utils.Validate(req_ticket_reply); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Validation error: %s", err))
	}

	userID, ok := c.Locals("user_id").(int64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

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

	sqlStatement := `
		WITH inserted AS (
			INSERT INTO ticket_replies (ticket_id, user_id, message, created_by)
			VALUES ($1, $2, $3, $4)
			RETURNING id, ticket_id, user_id, message, created_at, created_by
		)
		SELECT 
			i.id,
			i.ticket_id,
			i.user_id,
			i.message,
			i.created_at,
			i.created_by,
			u.role
		FROM inserted i
		JOIN users u ON i.user_id = u.id
	`

	var inserted data.TicketReplyResponse
	err = database.InitiateDataBase().QueryRow(
		c.Context(),
		sqlStatement,
		ticketID,
		userID,
		req_ticket_reply.Message,
		userEmail,
	).Scan(
		&inserted.ID,
		&inserted.TicketID,
		&inserted.UserID,
		&inserted.Message,
		&inserted.CreatedAt,
		&inserted.CreatedBy,
		&inserted.Role,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(inserted)
}
