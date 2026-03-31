package handler

import (
	"fmt"
	"subscription_manager/data"
	"subscription_manager/database"
	"subscription_manager/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

func UserRegister(c *fiber.Ctx) error {
	var req_user data.UserRegisterRequest

	if err := c.BodyParser(&req_user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	if err := utils.Validate(req_user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Validation error: %s", err))
	}

	if !utils.EmailValidation(req_user.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Please enter valid email",
		})
	}

	if !utils.PasswordValidation(req_user.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Please enter valid password",
		})
	}

	sqlstatement := `INSERT INTO public.users (first_name,last_name,email,password,role,created_at,updated_at)
	VALUES($1,$2,$3,$4,$5,NOW(),NOW())
	RETURNING id,first_name, last_name, email, role, created_at, updated_at`

	role := "user"
	var inserted data.UserResponse
	err := database.InitiateDataBase().QueryRow(
		c.Context(),
		sqlstatement,
		req_user.FirstName,
		req_user.LastName,
		req_user.Email,
		utils.GeneratePassword(req_user.Password),
		role,
	).Scan(
		&inserted.ID,
		&inserted.FirstName,
		&inserted.LastName,
		&inserted.Email,
		&inserted.Role,
		&inserted.Created_at,
		&inserted.Updated_at,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(inserted)
}

func UserLogin(c *fiber.Ctx) error {
	var reqUser data.UserLogin

	if err := c.BodyParser(&reqUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Please provide valid login data",
		})
	}

	if reqUser.Email == "" || reqUser.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email and password are required",
		})
	}

	if err := utils.Validate(reqUser); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Validation error: %s", err))
	}

	sqlStatement := `SELECT id, first_name, last_name, email, password, role, created_at, updated_at 
	                 FROM public.users WHERE email = $1 AND deleted_at IS NULL`

	user := data.UserResponse{}
	var hashedPassword string
	row := database.InitiateDataBase().QueryRow(
		c.Context(),
		sqlStatement,
		reqUser.Email,
	)
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&hashedPassword,
		&user.Role,
		&user.Created_at,
		&user.Updated_at,
	)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid email or password",
		})
	}

	if !utils.ComparePassword(hashedPassword, reqUser.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid email or password",
		})
	}

	role := "user"
	if user.Role == "admin" {
		role = "admin"
	}

	token, err := utils.GenerateToken(uint(user.ID), role, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate token",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})

	return c.Status(fiber.StatusOK).JSON(user)
}

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
