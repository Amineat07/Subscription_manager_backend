package adminhandler

import (
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
