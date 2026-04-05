package userhandler

import (
	"fmt"
	"strconv"
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
    INSERT INTO tickets (title, description, link, priority, user_id, created_by, updated_by)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING id, title, description, COALESCE(link, ''), priority, created_by, created_at, updated_by, updated_at

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
		userEmail,
		userEmail,
	).Scan(
		&inserted.ID,
		&inserted.Title,
		&inserted.Description,
		&inserted.Link,
		&inserted.Priority,
		&inserted.CreatedBy,
		&inserted.CreatedAt,
		&inserted.UpdatedBy,
		&inserted.UpdatedAt,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(inserted)

}

func GetMyTickets(c *fiber.Ctx) error {

	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user ID not found",
		})
	}

	sqlStatement := `SELECT id,user_id, title, description, COALESCE(link,''), priority, created_by, created_at, updated_by, updated_at
	FROM tickets 
	WHERE user_id = $1
	AND deleted_at IS NULL`

	rows, err := database.InitiateDataBase().Query(c.Context(), sqlStatement, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	defer rows.Close()

	tickets := []data.TicketResponse{}

	for rows.Next() {
		var t data.TicketResponse

		err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.Title,
			&t.Description,
			&t.Link,
			&t.Priority,
			&t.CreatedBy,
			&t.CreatedAt,
			&t.UpdatedBy,
			&t.UpdatedAt,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		tickets = append(tickets, t)
	}
	return c.Status(fiber.StatusOK).JSON(tickets)
}

func GetTicket(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user ID not found",
		})
	}

	id := c.Params("id")

	ticketID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	sqlStatement := `SELECT id,user_id, title, description, COALESCE(link,''), priority, created_by, created_at, updated_by, updated_at
	FROM tickets 
	WHERE id = $1 AND user_id = $2
	AND deleted_at IS NULL`

	var t data.TicketResponse

	err = database.InitiateDataBase().QueryRow(c.Context(), sqlStatement, ticketID, userID).Scan(
		&t.ID,
		&t.UserID,
		&t.Title,
		&t.Description,
		&t.Link,
		&t.Priority,
		&t.CreatedBy,
		&t.CreatedAt,
		&t.UpdatedBy,
		&t.UpdatedAt,
	)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "ticket not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(t)
}

func UpdateTicket(c *fiber.Ctx) error {

	var updateTicket data.TicketUpdateRequest
	if err := c.BodyParser(&updateTicket); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user ID not found",
		})
	}

	userEmail := c.Locals("userEmail")
	if userEmail == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user email not found",
		})
	}

	id := c.Params("id")

	ticketID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	sqlStatement := `UPDATE tickets
	SET title = COALESCE(NULLIF($1, ''), title),
		description = COALESCE(NULLIF($2, ''), description),
		link = COALESCE(NULLIF($3, ''), link),
		priority = COALESCE(NULLIF($4, ''), priority),
		updated_at = NOW(),
		updated_by = $5
		WHERE id = $6 AND user_id = $7 AND deleted_at IS NULL
		RETURNING id,user_id,title,description,link,priority,updated_at,updated_by,created_at,created_by
	`

	row := database.InitiateDataBase().QueryRow(
		c.Context(),
		sqlStatement,
		updateTicket.Title,
		updateTicket.Description,
		updateTicket.Link,
		updateTicket.Priority,
		userEmail,
		ticketID,
		userID,
	)

	var updated data.TicketResponse
	err = row.Scan(
		&updated.ID,
		&updated.UserID,
		&updated.Title,
		&updated.Description,
		&updated.Link,
		&updated.Priority,
		&updated.UpdatedAt,
		&updated.UpdatedBy,
		&updated.CreatedAt,
		&updated.CreatedBy,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ticket updated successfully",
		"ticket":  updated,
	})
}

func DeleteTicket(c *fiber.Ctx) error {
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
        UPDATE tickets
        SET deleted_at = NOW(),
            deleted_by = $1
        WHERE id = $2
          AND user_id = $3
          AND deleted_at IS NULL
    `

	result, err := database.InitiateDataBase().Exec(
		c.Context(),
		sqlStatement,
		userEmail,
		ticketID,
		userID,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if result.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "ticket not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ticket deleted successfully",
	})
}
