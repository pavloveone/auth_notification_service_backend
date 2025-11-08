package main

import (
	userhandlers "auth_notification_service/internal/handlers/user"
	userrepository "auth_notification_service/internal/repositories/user"
	userroutes "auth_notification_service/internal/routes/user"
	userservice "auth_notification_service/internal/services/user"
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

func main() {
	app := fiber.New()

	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))

	if err != nil {
		logrus.Fatal("Failed to connect to database:", err)
	}

	userRepository, err := userrepository.NewUserRepository(ctx, dbpool)
	if err != nil {
		logrus.Fatal("Error initializing user repository:", err)
	}
	userService := userservice.NewUserService(userRepository)
	userHandler := userhandlers.NewUserHandler(userService, ctx)
	userroutes.SetupRoutes(app, userHandler)

	port := "8080"
	logrus.WithFields(logrus.Fields{
		"port": port,
	}).Info("Server starting on port")
	logrus.Fatal(app.Listen(":" + port))
}
