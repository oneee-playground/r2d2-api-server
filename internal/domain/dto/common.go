package dto

import "github.com/google/uuid"

type UserInfo struct {
	ID         uuid.UUID `json:"id"`
	Username   string    `json:"username"`
	ProfileURL string    `json:"profileURL"`
	Role       string    `json:"role"`
}

type IDInput struct {
	ID uuid.UUID
}

type IDOutput struct {
	ID uuid.UUID `json:"id"`
}
