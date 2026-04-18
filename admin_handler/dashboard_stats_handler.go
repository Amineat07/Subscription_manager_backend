package adminhandler

import (
	"subscription_manager/database"

	"github.com/gofiber/fiber/v2"
)

func AdminDashboardStats(c *fiber.Ctx) error {

	var users int
	var subscriptions int
	var newsfeed int
	var open, inProgress, closed int

	db := database.InitiateDataBase()

	db.QueryRow(c.Context(), `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`).Scan(&users)

	db.QueryRow(c.Context(), `SELECT COUNT(*) FROM subscriptions WHERE deleted_at IS NULL`).Scan(&subscriptions)

	db.QueryRow(c.Context(), `SELECT COUNT(*) FROM news_feed WHERE deleted_at IS NULL`).Scan(&newsfeed)

	db.QueryRow(c.Context(), `SELECT COUNT(*) FROM tickets WHERE status = 'open' AND deleted_at IS NULL`).Scan(&open)
	db.QueryRow(c.Context(), `SELECT COUNT(*) FROM tickets WHERE status = 'in_progress' AND deleted_at IS NULL`).Scan(&inProgress)
	db.QueryRow(c.Context(), `SELECT COUNT(*) FROM tickets WHERE status = 'closed' AND deleted_at IS NULL`).Scan(&closed)

	return c.JSON(fiber.Map{
		"users":         users,
		"subscriptions": subscriptions,
		"newsfeed":      newsfeed,
		"tickets": fiber.Map{
			"open":        open,
			"in_progress": inProgress,
			"closed":      closed,
		},
	})
}
