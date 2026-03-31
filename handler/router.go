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

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	public := app.Group("/api/v1")

	public.Post("/register", UserRegister)
	public.Post("/login", UserLogin)

	user := app.Group("/api/v1")

	user.Use(utils.JWTMiddleware(jwtSecret))
	user.Use(utils.RequireRole("admin", "user"))

	user.Post("/logout", UserLogout)

	user.Get("/me", GetMyAuthInfo)
	user.Patch("/me", UpdateMyAccount)
	user.Delete("/me", DeleteMyAccount)

	user.Post("/subscriptions", AddSubscription)
	user.Get("/subscriptions", GetSubscriptionsByUserID)
	user.Get("/subscriptions/:id", GetSubscriptionByUserID)
	user.Patch("/subscriptions/:id", UpdateSubscriptionByUserID)
	user.Delete("/subscriptions/:id", DeleteSubscriptionByUserID)

	// user.Get("/newsfeed", GetNewsFeed)
	// user.Get("/newsfeed/:id", GetNewsFeedItem)

	// user.Post("/tickets", CreateTicket)
	// user.Get("/tickets", GetMyTickets)
	// user.Get("/tickets/:id", GetTicket)
	// user.Post("/tickets/:id/reply", ReplyToTicket)
	// user.Post("/tickets/:id/close", CloseTicket)

	admin := app.Group("/admin/api/v1")

	admin.Use(utils.JWTMiddleware(jwtSecret))
	admin.Use(utils.RequireRole("admin"))

	// admin.Get("/users", AdminListUsers)
	// admin.Get("/users/:id", AdminGetUser)
	// admin.Patch("/users/:id", AdminUpdateUser)
	// admin.Delete("/users/:id", AdminDeleteUser)
	// admin.Post("/users/:id/ban", AdminBanUser)
	// admin.Post("/users/:id/unban", AdminUnbanUser)

	// admin.Get("/subscriptions", AdminGetSubscriptions)
	// admin.Get("/subscriptions/:id", AdminGetSubscription)
	// admin.Post("/subscriptions", AdminCreateSubscription)
	// admin.Patch("/subscriptions/:id", AdminUpdateSubscription)
	// admin.Delete("/subscriptions/:id", AdminDeleteSubscription)
	// admin.Post("/subscriptions/:id/cancel", AdminCancelSubscription)
	// admin.Post("/subscriptions/:id/renew", AdminRenewSubscription)

	// admin.Get("/requests", AdminListUsersRequests)
	// admin.Get("/requests/:id", AdminGetUserRequest)
	// admin.Post("/requests/:id/answer", AdminAnswerUserRequest)

	// admin.Post("/newsfeed", AdminCreateNewsFeed)
	// admin.Get("/newsfeed", AdminListNewsFeed)
	// admin.Get("/newsfeed/:id", AdminGetNewsFeedItem)
	// admin.Patch("/newsfeed/:id", AdminUpdateNewsFeed)
	// admin.Delete("/newsfeed/:id", AdminDeleteNewsFeed)

	// admin.Get("/tickets", AdminListTickets)
	// admin.Get("/tickets/:id", AdminGetTicket)
	// admin.Patch("/tickets/:id", AdminUpdateTicket)
	// admin.Post("/tickets/:id/reply", AdminReplyTicket)
	// admin.Post("/tickets/:id/close", AdminCloseTicket)

}
