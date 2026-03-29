package main

import (
	"subscription_manager/database"
	"subscription_manager/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	connection := database.InitiateDataBase()
	defer connection.Close()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))

	handler.SetupRouter(app, connection)

	app.Listen(":3000")

}
