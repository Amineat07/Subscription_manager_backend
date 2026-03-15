package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRouter(app *fiber.App, con *pgxpool.Pool) {

	app.Use(logger.New())

	subscription := app.Group("/subscription")
	subscription.Post("/", AddSubscription)
	subscription.Get("/", GetSubscriptions)
	subscription.Get("/:id", GetSubscription)
	subscription.Patch("/:id", UpdateSubscription)
	subscription.Delete("/:id", DeleteSubscription)
}
