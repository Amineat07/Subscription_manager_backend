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

func AdminUpdateUserById(c *fiber.Ctx) error {

	var updateReq data.AdminUpdateUserRequest
	if err := c.BodyParser(&updateReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	id := c.Params("id")

	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	userEmail := c.Locals("userEmail")
	if userEmail == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user email not found",
		})
	}

	sqlStatement := `
        UPDATE users
        SET first_name = COALESCE(NULLIF($1, ''), first_name),
    		last_name  = COALESCE(NULLIF($2, ''), last_name),
    		email      = COALESCE(NULLIF($3, ''), email),
    		role       = COALESCE(NULLIF($4, ''), role),
    		updated_at = NOW(),
    		updated_by = $5   
			WHERE id = $6        
          AND deleted_at IS NULL
        RETURNING id, first_name, last_name, email,role,updated_by, updated_at
    `

	row := database.InitiateDataBase().QueryRow(
		c.Context(),
		sqlStatement,
		updateReq.FirstName,
		updateReq.LastName,
		updateReq.Email,
		updateReq.Role,
		userEmail,
		userID,
	)

	var updated data.AdminUpdateUserResponse
	err = row.Scan(
		&updated.ID,
		&updated.FirstName,
		&updated.LastName,
		&updated.Email,
		&updated.Role,
		&updated.Updatedby,
		&updated.Updated_at,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "account updated successfully",
		"user":    updated,
	})

}

func AdminDeleteUser(c *fiber.Ctx) error {

	id := c.Params("id")

	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	userEmail := c.Locals("userEmail")
	if userEmail == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user email not found",
		})
	}

	sqlStatement := `
		UPDATE public.users
		SET deleted_at = NOW(),
    		deleted_by = $1
			WHERE id = $2 
	`
	result, err := database.InitiateDataBase().Exec(c.Context(), sqlStatement, userEmail, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if result.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "user deleted successfully",
	})
}
