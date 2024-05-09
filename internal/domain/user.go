package domain

import "github.com/google/uuid"

type UserRole uint8

const (
	RoleMember UserRole = iota
	RoleAdmin
)

type User struct {
	ID         uuid.UUID
	Username   string
	Email      string
	ProfileURL string
	Role       UserRole
}

func (u User) IsAdmin() bool {
	return u.Role == RoleAdmin
}
