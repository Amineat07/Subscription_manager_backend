package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(id uint, role string, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"role":    role,
		"email":   email,
		"exp":     jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
	})

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return t, nil
}

func JWTMiddleware(secret []byte) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := c.Cookies("jwt")

		if tokenStr == "" {
			authHeader := c.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "missing token",
			})
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid token",
			})
		}

		claims := token.Claims.(jwt.MapClaims)

		var userID int64
		switch id := claims["user_id"].(type) {
		case float64:
			userID = int64(id)
		case string:
			v, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "invalid user_id in token",
				})
			}
			userID = v
		default:
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid user_id in token",
			})
		}
		c.Locals("user_id", userID)

		if role, ok := claims["role"].(string); ok {
			c.Locals("role", role)
		}

		if email, ok := claims["email"].(string); ok {
			c.Locals("userEmail", email)
		}

		return c.Next()
	}
}

func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok || role == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "forbidden: no role assigned",
			})
		}

		for _, r := range roles {
			if role == r {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "forbidden: insufficient permissions",
		})
	}
}
