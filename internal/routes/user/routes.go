package routes

import (
	"auth_notification_service/internal/auth"
	userhandlers "auth_notification_service/internal/handlers/user"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, h *userhandlers.UserHandler) {
	app.Post("/register", h.AddNewUser)
	app.Post("/login", h.LogInUser)
	app.Get("/logout", h.LogoutUser)

	userGroup := app.Group("/users", auth.AuthMiddleware)
	userGroup.Get("/", h.AllUsers)
	userGroup.Get("/:id", h.UserById)
}
