package models

import "awesomeProject/internal/app/ds"

type CreateUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	Registration string   `json:"registration"`
	User         *ds.User `json:"user"`
}

type UpdateUserRequest struct {
	ID       uint   `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthUserResponse struct {
	Auth string   `json:"auth"`
	User *ds.User `json:"user"`
}

type LogoutUserRequest struct {
	ID uint `json:"id"`
}
