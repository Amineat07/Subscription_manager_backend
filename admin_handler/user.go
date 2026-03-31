package adminhandler

import (
	"strconv"
	"subscription_manager/data"
	"subscription_manager/database"

	"github.com/gofiber/fiber/v2"
)

func AdminListUsers(c *fiber.Ctx) error {
	rows, err := database.InitiateDataBase().Query(c.Context(), `
        SELECT
            id,
            first_name,
            last_name,
            email,
            role,
            created_at,
            updated_at
        FROM users
        WHERE deleted_at IS NULL
        ORDER BY created_at DESC
    `)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	defer rows.Close()

	users := []data.UserResponse{}
	for rows.Next() {
		var u data.UserResponse
		if err := rows.Scan(
			&u.ID,
			&u.FirstName,
			&u.LastName,
			&u.Email,
			&u.Role,
			&u.Created_at,
			&u.Updated_at,
		); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		users = append(users, u)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"users": users,
	})
}

func AdminGetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	sqlstatement := `
        SELECT
            id,
            first_name,
            last_name,
            email,
            role,
            created_at,
            updated_at
        FROM users
        WHERE id = $1 AND deleted_at IS NULL
    `

	var u data.UserResponse
	err = database.InitiateDataBase().QueryRow(c.Context(), sqlstatement, user_id).Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Role,
		&u.Created_at,
		&u.Updated_at,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(u)

}
