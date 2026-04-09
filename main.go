package main

import (
	"fmt"
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

	origins := fmt.Sprintf("%s,%s",
		os.Getenv("ADMIN_FRONTEND_BASE_URL"),
		os.Getenv("USER_FRONTEND_BASE_URL"),
	)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))

	router.SetupRouter(app, connection)

	app.Listen(":3000")

}
