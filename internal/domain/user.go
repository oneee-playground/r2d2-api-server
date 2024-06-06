package domain

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/oneee-playground/r2d2-api-server/internal/domain/dto"
)

//go:generate mockgen -source=user.go -destination=../../test/mocks/user.go -package=mocks

type UserRole uint8

const (
	RoleMember UserRole = iota + 1
	RoleAdmin
)

func (r UserRole) String() string {
	switch r {
	case RoleAdmin:
		return "ADMIN"
	case RoleMember:
		return "MEMBER"
	}

	panic("unreachable")
}

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

type AuthUsecase interface {
	SignIn(ctx context.Context, in dto.SignInInput) (out *dto.AccessTokenOutput, err error)
}

type UserUsecase interface {
	GetSelfInfo(ctx context.Context) (out *dto.UserInfo, err error)
}

// Defined errors for UserRepository.
var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository interface {
	UsernameExists(ctx context.Context, username string) (bool, error)
	FetchByUsername(ctx context.Context, username string) (User, error)
	FetchByID(ctx context.Context, id uuid.UUID) (User, error)
	Create(ctx context.Context, user User) error
}
