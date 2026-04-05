package adminhandler

import (
	"strconv"
	"subscription_manager/data"
	"subscription_manager/database"

	"github.com/gofiber/fiber/v2"
)

func AdminGetTickets(c *fiber.Ctx) error {

	sqlStatement := `
        SELECT id, user_id, title, description, COALESCE(link,''), priority, status, created_by, created_at, updated_by, updated_at
        FROM tickets 
        WHERE deleted_at IS NULL
    `

	rows, err := database.InitiateDataBase().Query(c.Context(), sqlStatement)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

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
			&t.Status,
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
	rows.Close()

	for i, t := range tickets {
		replyRows, err := database.InitiateDataBase().Query(c.Context(), `
            SELECT r.id, r.ticket_id, r.user_id, r.message, r.created_at, r.created_by, u.role
            FROM ticket_replies r
            JOIN users u ON r.user_id = u.id
            WHERE r.ticket_id = $1
            ORDER BY r.created_at ASC
        `, t.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		replies := []data.TicketReplyResponse{}
		for replyRows.Next() {
			var r data.TicketReplyResponse
			err := replyRows.Scan(
				&r.ID,
				&r.TicketID,
				&r.UserID,
				&r.Message,
				&r.CreatedAt,
				&r.CreatedBy,
				&r.Role,
			)
			if err != nil {
				replyRows.Close()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			replies = append(replies, r)
		}
		replyRows.Close()
		tickets[i].Replies = replies
	}

	return c.Status(fiber.StatusOK).JSON(tickets)
}

func AdminGetTicketsByUserID(c *fiber.Ctx) error {

	id := c.Params("id")
	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	var exists bool
	err = database.InitiateDataBase().QueryRow(c.Context(), `
        SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)
    `, userID).Scan(&exists)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	sqlStatement := `
        SELECT id, user_id, title, description, COALESCE(link,''), priority, status, created_by, created_at, updated_by, updated_at
        FROM tickets 
        WHERE user_id = $1 AND deleted_at IS NULL
    `

	rows, err := database.InitiateDataBase().Query(c.Context(), sqlStatement, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

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
			&t.Status,
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
	rows.Close()

	for i, t := range tickets {
		replyRows, err := database.InitiateDataBase().Query(c.Context(), `
            SELECT r.id, r.ticket_id, r.user_id, r.message, r.created_at, r.created_by, u.role
            FROM ticket_replies r
            JOIN users u ON r.user_id = u.id
            WHERE r.ticket_id = $1
            ORDER BY r.created_at ASC
        `, t.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		replies := []data.TicketReplyResponse{}
		for replyRows.Next() {
			var r data.TicketReplyResponse
			err := replyRows.Scan(
				&r.ID,
				&r.TicketID,
				&r.UserID,
				&r.Message,
				&r.CreatedAt,
				&r.CreatedBy,
				&r.Role,
			)
			if err != nil {
				replyRows.Close()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			replies = append(replies, r)
		}
		replyRows.Close()
		tickets[i].Replies = replies
	}

	return c.Status(fiber.StatusOK).JSON(tickets)
}

func AdminGetTicket(c *fiber.Ctx) error {

	id := c.Params("id")
	ticketID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	sqlStatement := `SELECT id,user_id, title, description, COALESCE(link,''), priority, status, created_by, created_at, updated_by, updated_at
	FROM tickets 
	WHERE id = $1
	AND deleted_at IS NULL`

	var t data.TicketResponse
	err = database.InitiateDataBase().QueryRow(c.Context(), sqlStatement, ticketID).Scan(
		&t.ID,
		&t.UserID,
		&t.Title,
		&t.Description,
		&t.Link,
		&t.Priority,
		&t.Status,
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

	replyRows, err := database.InitiateDataBase().Query(c.Context(), `
        SELECT r.id, r.ticket_id, r.user_id, r.message, r.created_at, r.created_by, u.role
        FROM ticket_replies r
        JOIN users u ON r.user_id = u.id
        WHERE r.ticket_id = $1
        ORDER BY r.created_at ASC
    `, t.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	defer replyRows.Close()

	replies := []data.TicketReplyResponse{}
	for replyRows.Next() {
		var r data.TicketReplyResponse
		err := replyRows.Scan(
			&r.ID,
			&r.TicketID,
			&r.UserID,
			&r.Message,
			&r.CreatedAt,
			&r.CreatedBy,
			&r.Role,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		replies = append(replies, r)
	}

	t.Replies = replies

	return c.Status(fiber.StatusOK).JSON(t)
}

func AdminGetTicketByUserID(c *fiber.Ctx) error {
	user_id := c.Params("user_id")
	userID, err := strconv.ParseInt(user_id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	ticket_id := c.Params("ticket_id")
	ticketID, err := strconv.ParseInt(ticket_id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	var exists bool
	err = database.InitiateDataBase().QueryRow(c.Context(), `
        SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)
    `, userID).Scan(&exists)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	sqlStatement := `SELECT id,user_id, title, description, COALESCE(link,''), priority, status, created_by, created_at, updated_by, updated_at
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
		&t.Status,
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

	replyRows, err := database.InitiateDataBase().Query(c.Context(), `
        SELECT r.id, r.ticket_id, r.user_id, r.message, r.created_at, r.created_by, u.role
        FROM ticket_replies r
        JOIN users u ON r.user_id = u.id
        WHERE r.ticket_id = $1
        ORDER BY r.created_at ASC
    `, t.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	defer replyRows.Close()

	replies := []data.TicketReplyResponse{}
	for replyRows.Next() {
		var r data.TicketReplyResponse
		err := replyRows.Scan(
			&r.ID,
			&r.TicketID,
			&r.UserID,
			&r.Message,
			&r.CreatedAt,
			&r.CreatedBy,
			&r.Role,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		replies = append(replies, r)
	}

	t.Replies = replies

	return c.Status(fiber.StatusOK).JSON(t)
}

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
