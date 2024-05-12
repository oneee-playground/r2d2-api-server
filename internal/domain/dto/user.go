package dto

import "github.com/google/uuid"

type UserInfo struct {
	ID         uuid.UUID `json:"id"`
	Username   string    `json:"username"`
	ProfileURL string    `json:"profileURL"`
	Role       string    `json:"role"`
}

type SignInInput struct {
	Code string `json:"code"`
}

type AccessTokenOutput struct {
	Token string `json:"token"`
}
