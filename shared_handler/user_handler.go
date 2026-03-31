package sharedhandler

import (
	"subscription_manager/data"
	"subscription_manager/database"
	"time"

	"github.com/gofiber/fiber/v2"
)



func GetUser(c *fiber.Ctx) error {

	userIDVal := c.Locals("user_id")
	if userIDVal == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	userID, ok := userIDVal.(int64)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "invalid user id type",
		})
	}

	sqlStatement := `
		SELECT 
			id,
			first_name,
			last_name,
			email,
			created_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	var u data.UserResponse

	err := database.InitiateDataBase().
		QueryRow(c.Context(), sqlStatement, userID).
		Scan(
			&u.ID,
			&u.FirstName,
			&u.LastName,
			&u.Email,
			&u.Created_at,
		)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(u)
}

func GetMyAuthInfo(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user ID not found",
		})
	}

	sqlStatement := `
		SELECT 
			id,
			first_name,
			last_name,
			email,
			role,
			created_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	var u data.UserResponse

	err := database.InitiateDataBase().
		QueryRow(c.Context(), sqlStatement, userID).
		Scan(
			&u.ID,
			&u.FirstName,
			&u.LastName,
			&u.Email,
			&u.Role,
			&u.Created_at,
		)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(u)
}

func UserLogout(c *fiber.Ctx) error {

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HTTPOnly: true,
		Secure:   false,
		SameSite: fiber.CookieSameSiteStrictMode,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logged out successfully",
	})

}

func UpdateMyAccount(c *fiber.Ctx) error {
	var updateReq data.UpdateUserRequest
	if err := c.BodyParser(&updateReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	userID, ok := c.Locals("user_id").(int64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	sqlStatement := `
        UPDATE users
        SET first_name = COALESCE(NULLIF($1, ''), first_name),
            last_name  = COALESCE(NULLIF($2, ''), last_name),
            email      = COALESCE(NULLIF($3, ''), email),
            updated_at = NOW(),
            updated_by = $4
        WHERE id = $4
          AND deleted_at IS NULL
        RETURNING id, first_name, last_name, email, updated_at, role
    `

	row := database.InitiateDataBase().QueryRow(
		c.Context(),
		sqlStatement,
		updateReq.FirstName,
		updateReq.LastName,
		updateReq.Email,
		userID,
	)

	var updated data.UserResponse
	err := row.Scan(
		&updated.ID,
		&updated.FirstName,
		&updated.LastName,
		&updated.Email,
		&updated.Updated_at,
		&updated.Role,
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

func DeleteMyAccount(c *fiber.Ctx) error {

	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user ID not found",
		})
	}

	sqlStatement := `
		UPDATE public.users
		SET deleted_at = NOW(),
			deleted_by = $1
		WHERE id = $1
	`
	result, err := database.InitiateDataBase().Exec(c.Context(), sqlStatement, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if result.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "subscription not found",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HTTPOnly: true,
		Secure:   false,
		SameSite: fiber.CookieSameSiteStrictMode,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "user deleted successfully",
	})
}
