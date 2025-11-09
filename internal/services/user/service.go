package userservice

import (
	"auth_notification_service/internal/models"
	userrepository "auth_notification_service/internal/repositories/user"
	"context"
)

type UserService struct {
	repo *userrepository.UserRepository
}

func NewUserService(repo *userrepository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s UserService) AllUsers(ctx context.Context) ([]models.UserResponse, error) {
	return s.repo.AllUsers(ctx)
}

func (s UserService) UserById(ctx context.Context, id int) (models.UserResponse, error) {
	return s.repo.UserById(ctx, id)
}

func (s UserService) AddNewUser(ctx context.Context, request models.UserCreateRequest) (int, error) {
	return s.repo.AddNewUser(ctx, request)
}

func (s UserService) LogInUser(ctx context.Context, request models.UserLogInRequest) (models.UserLogInResponse, error) {
	return s.repo.LogIn(ctx, request)
}

func (s UserService) LogOutUser(ctx context.Context) (bool, error) {
	return s.repo.LogOut(ctx)
}
