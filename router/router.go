package router

import (
	"os"
	adminhandler "subscription_manager/admin_handler"
	publichandler "subscription_manager/public_handler"
	sharedhandler "subscription_manager/shared_handler"
	userhandler "subscription_manager/user_handler"
	"subscription_manager/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRouter(app *fiber.App, con *pgxpool.Pool) {

	app.Use(logger.New())

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	public := app.Group("/api/v1")

	public.Post("/register", publichandler.UserRegister)
	public.Post("/login", publichandler.UserLogin)

	user := app.Group("/api/v1")

	user.Use(utils.JWTMiddleware(jwtSecret))
	user.Use(utils.RequireRole("admin", "user"))

	user.Post("/logout", sharedhandler.UserLogout)

	user.Get("/me", sharedhandler.GetMyAuthInfo)
	user.Patch("/me", sharedhandler.UpdateMyAccount)
	user.Delete("/me", sharedhandler.DeleteMyAccount)

	user.Post("/subscriptions", userhandler.AddSubscription)
	user.Get("/subscriptions", userhandler.GetSubscriptionsByUserID)
	user.Get("/subscriptions/:id", userhandler.GetSubscriptionByUserID)
	user.Patch("/subscriptions/:id", userhandler.UpdateSubscriptionByUserID)
	user.Delete("/subscriptions/:id", userhandler.DeleteSubscriptionByUserID)

	// user.Get("/newsfeed", GetNewsFeed)
	// user.Get("/newsfeed/:id", GetNewsFeedItem)

	user.Post("/tickets", userhandler.CreateTicket)
	// user.Get("/tickets", GetMyTickets)
	// user.Get("/tickets/:id", GetTicket)
	// user.Post("/tickets/:id/reply", ReplyToTicket)
	// user.Post("/tickets/:id/close", CloseTicket)

	admin := app.Group("/admin/api/v1")

	admin.Use(utils.JWTMiddleware(jwtSecret))
	admin.Use(utils.RequireRole("admin"))

	admin.Get("/users", adminhandler.AdminListUsers)
	admin.Get("/users/:id", adminhandler.AdminGetUser)
	// admin.Patch("/users/:id", AdminUpdateUser)
	// admin.Delete("/users/:id", AdminDeleteUser)
	// admin.Post("/users/:id/ban", AdminBanUser)
	// admin.Post("/users/:id/unban", AdminUnbanUser)

	// admin.Get("/subscriptions", AdminGetAllSubscriptions)
	// admin.Get("/subscriptions/:id", AdminGetAllSubscriptionByUserId)
	// admin.Get("/subscriptions/:id", AdminGetSubscription)
	// admin.Get("/subscriptions/:id", AdminGetSubscriptionByUserId)
	// admin.Post("/subscriptions", AdminCreateSubscriptionByUserId)
	// admin.Patch("/subscriptions/:id", AdminUpdateSubscriptionByUserId)
	// admin.Delete("/subscriptions/:id", AdminDeleteSubscriptionByUserId)

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

	//API Publishing
	// publish := app.Group("/publish/api/v1")
	// publish.Get("/subscriptions", GetSubscriptions) //should have query params like company id and ApiKey
	// publish.Get("/subscriptions/:id", GetUserSubscriptions)

}
