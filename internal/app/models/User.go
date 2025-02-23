package models

import "awesomeProject/internal/app/ds"

type CreateUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	isAdmin  bool   `json:"is_admin"`
}

type CreateUserResponse struct {
	User *ds.User `json:"user"`
}

type UpdateUserRequest struct {
	ID       uint   `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthUserResponse struct {
	Token string `json:"token"`
}

type LogoutUserRequest struct {
	Login string `json:"login"`
}

type LoginUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}
