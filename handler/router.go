package handler

import (
	"os"
	"subscription_manager/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRouter(app *fiber.App, con *pgxpool.Pool) {

	app.Use(logger.New())

	user := app.Group("/user")
	user.Post("/register", UserRegister)
	user.Post("/login", UserLogin)
	user.Use(utils.JWTMiddleware([]byte(os.Getenv("JWT_SECRET"))))
	user.Get("/",GetUser)
	user.Get("/me",GetMyAuthInfo)
	user.Delete("/",DeleteMyAccount)
	user.Post("/logout", UserLogout)

	subscription := app.Group("/subscription")
	subscription.Use(utils.JWTMiddleware([]byte(os.Getenv("JWT_SECRET"))))
	subscription.Post("/", AddSubscription)
	subscription.Get("/", GetSubscriptions)
	subscription.Get("/user/:id",GetSubscriptionByUserID)
	subscription.Get("/:id", GetSubscription)
	subscription.Patch("/:id", UpdateSubscription)
	subscription.Delete("/:id", DeleteSubscription)
}
