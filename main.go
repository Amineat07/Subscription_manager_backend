package main

import (
	"log"
	"os"
	"subscription_manager/database"
	"subscription_manager/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {

	connection := database.InitiateDataBase()
	defer connection.Close()

	app := fiber.New()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("FRONTEND_BASE_URL"),
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))

	handler.SetupRouter(app, connection)

	app.Listen(":3000")

}
