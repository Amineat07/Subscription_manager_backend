package userhandler

import (
	"fmt"
	"subscription_manager/data"
	"subscription_manager/database"
	"subscription_manager/utils"

	"github.com/gofiber/fiber/v2"
)

func CreateTicket(c *fiber.Ctx) error {

	var req data.TicketRequest

	userID, ok := c.Locals("user_id").(int64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	body := c.Body()
	fmt.Println(string(body))
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	if err := utils.Validate(req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(
			fmt.Sprintf("Validation error: %s", err),
		)
	}

	userEmail := c.Locals("userEmail")
	if userEmail == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user email not found",
		})
	}

	sqlstatement := `
    INSERT INTO tickets (title, description, link, priority, user_id, created_by)
    VALUES ($1, $2, $3, $4, $5, $5)
    RETURNING id, title, description, link, priority, created_at
`

	var inserted data.TicketResponse
	err := database.InitiateDataBase().QueryRow(
		c.Context(),
		sqlstatement,
		req.Title,
		req.Description,
		req.Link,
		req.Priority,
		userID,
	).Scan(
		&inserted.ID,
		&inserted.Title,
		&inserted.Description,
		&inserted.Link,
		&inserted.Priority,
		&inserted.Created_at,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(inserted)

}
