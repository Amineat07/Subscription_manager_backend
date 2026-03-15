package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRouter(app *fiber.App,con *pgxpool.Pool) {
		app.Post("/subscription", AddSubscription)
}