package models

import "time"

type User struct {
	ID        int
	Username  string
	Password  string
	Email     string
	CreatedOn time.Time
}

type UserCreateRequest struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=8,max=64"`
	Email    string `json:"email" validate:"required,email"`
}

type UserResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedOn time.Time `json:"createdOn"`
}

type UserLogInRequest struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

type UserLogInResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
}
