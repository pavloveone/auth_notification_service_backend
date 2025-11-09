package userhandlers

import (
	"auth_notification_service/internal/models"
	userrepository "auth_notification_service/internal/repositories/user"
	userservice "auth_notification_service/internal/services/user"
	"context"
	"database/sql"
	"errors"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service *userservice.UserService
	ctx     context.Context
}

func NewUserHandler(service *userservice.UserService, ctx context.Context) *UserHandler {
	return &UserHandler{service: service, ctx: ctx}
}

func (h *UserHandler) AllUsers(c *fiber.Ctx) error {
	users, err := h.service.AllUsers(h.ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Users not found"})
		}
		log.Printf("failed to find all users: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "internal server error"})
	}
	return c.Status(fiber.StatusOK).JSON(users)
}

func (h *UserHandler) UserById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user id"})
	}
	user, err := h.service.UserById(h.ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "user not found"})
		}
		log.Printf("internal error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "internal server error"})

	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *UserHandler) AddNewUser(c *fiber.Ctx) error {
	req := models.UserCreateRequest{}
	if err := c.BodyParser(&req); err != nil {
		log.Printf("failed to parse body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request body"})
	}
	id, err := h.service.AddNewUser(h.ctx, req)
	if err != nil {
		if errors.Is(err, userrepository.ErrUserAlreadyExists) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": "user with this username or email already exists",
			})
		}
		log.Printf("failed to add new user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to create user"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"id": id})
}

func (h *UserHandler) LogInUser(c *fiber.Ctx) error {
	req := models.UserLogInRequest{}
	if err := c.BodyParser(&req); err != nil {
		log.Printf("failed to parse body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request body"})
	}
	user, err := h.service.LogInUser(h.ctx, req)
	if err != nil {
		log.Printf("failed to login: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "current user doesn't exist"})
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *UserHandler) LogoutUser(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"accessToken":  "",
		"refreshToken": "",
		"message":      "logout successful",
	})
}
