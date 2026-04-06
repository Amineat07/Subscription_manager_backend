package main

import (
	"log"
	"os"
	"subscription_manager/cron"
	"subscription_manager/database"
	"subscription_manager/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connection := database.InitiateDataBase()
	defer connection.Close()

	cron.StartCronJobs()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("FRONTEND_BASE_URL"),
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))

	router.SetupRouter(app, connection)

	app.Listen(":3000")

}
