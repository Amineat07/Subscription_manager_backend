package publichandler

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
