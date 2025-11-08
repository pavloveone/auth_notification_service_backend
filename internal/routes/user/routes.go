package routes

import (
	"auth_notification_service/internal/auth"
	userhandlers "auth_notification_service/internal/handlers/user"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, h *userhandlers.UserHandler) {
	userGroup := app.Group("/users")
	userGroup.Get("/", h.AllUsers)
	userGroup.Get("/:id", auth.AuthMiddleware, h.UserById,)
	userGroup.Post("", h.AddNewUser)
	userGroup.Post("/login", h.LogInUser)
}
