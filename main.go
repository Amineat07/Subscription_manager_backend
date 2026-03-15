package main

import (
	"subscription_manager/database"
	"subscription_manager/handler"

	"github.com/gofiber/fiber/v2"
)

func main() {

	connection := database.InitiateDataBase()
	app := fiber.New()

	handler.SetupRouter(app,connection)

	app.Listen(":3000")

}
