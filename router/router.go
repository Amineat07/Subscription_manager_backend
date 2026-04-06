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

	user.Get("/newsfeed", sharedhandler.GetNewsFeed)
	user.Get("/newsfeed/:id", sharedhandler.GetNewsFeed)

	user.Post("/tickets", userhandler.CreateTicket)
	user.Get("/tickets", userhandler.GetMyTickets)
	user.Get("/tickets/:id", userhandler.GetTicket)
	user.Patch("/tickets/:id", userhandler.UpdateTicket)
	user.Delete("/tickets/:id", userhandler.DeleteTicket)
	user.Post("/tickets/:id/reply", sharedhandler.ReplyToTicket)

	admin := app.Group("/admin/api/v1")

	admin.Use(utils.JWTMiddleware(jwtSecret))
	admin.Use(utils.RequireRole("admin"))

	admin.Get("/users", adminhandler.AdminListUsers)
	admin.Get("/users/:id", adminhandler.AdminGetUser)
	admin.Patch("/users/:id", adminhandler.AdminUpdateUserById)
	admin.Delete("/users/:id", adminhandler.AdminDeleteUser)

	admin.Get("/subscriptions", adminhandler.AdminGetAllSubscriptions)
	admin.Get("/subscriptions/:id", adminhandler.AdminGetAllSubscriptionByUserId)
	admin.Get("/subscription/:id", adminhandler.AdminGetSubscription)
	admin.Get("/users/:user_id/subscriptions/:subscription_id", adminhandler.AdminGetSubscriptionByUserId)

	admin.Post("/newsfeed", adminhandler.AdminCreateNewsFeed)
	admin.Get("/newsfeed", sharedhandler.GetNewsFeed)
	admin.Get("/newsfeed/:id", sharedhandler.GetNewsFeed)
	admin.Patch("/newsfeed/:id", adminhandler.AdminUpdateNewsFeed)
	admin.Delete("/newsfeed/:id", adminhandler.AdminDeleteNewsFeed)

	admin.Get("/tickets", adminhandler.AdminGetTickets)
	admin.Get("/tickets/:id", adminhandler.AdminGetTicketsByUserID)
	admin.Get("/ticket/:id", adminhandler.AdminGetTicket)
	admin.Get("/users/:user_id/tickets/:ticket_id", adminhandler.AdminGetTicketByUserID)
	admin.Post("/tickets/:id/reply", sharedhandler.ReplyToTicket)
	admin.Put("/tickets/:id/status", adminhandler.UpdateTicketStatus)

	//API Publishing
	// publish := app.Group("/publish/api/v1")
	// publish.Get("/subscriptions", GetSubscriptions) //should have params like company id and ApiKey
	// publish.Get("/subscriptions/:id", GetUserSubscriptions)

}
