package handler

import (
	"fmt"
	"subscription_manager/data"
	"subscription_manager/database"
	"subscription_manager/utils"

	"github.com/gofiber/fiber/v2"
)

func RegisterUser(c *fiber.Ctx) error {
	var req_user data.UserRegisterRequest

	if err := c.BodyParser(&req_user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	if err := utils.Validate(req_user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Validation error: %s", err))
	}

	sqlstatement := `INSERT INTO public.users (first_name,last_name,email,password,is_admin,created_at,updated_at)
	VALUES($1,$2,$3,$4,$5,NOW(),NOW())
	RETURNING first_name, last_name, email, is_admin, created_at, updated_at`

	var inserted data.UserResponse
	err := database.InitiateDataBase().QueryRow(
		c.Context(),
		sqlstatement,
		req_user.FirstName,
		req_user.LastName,
		req_user.Email,
		utils.GeneratePassword(req_user.Password),
		req_user.IsAdmin,
	).Scan(
		&inserted.FirstName,
		&inserted.LastName,
		&inserted.Email,
		&inserted.IsAdmin,
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
