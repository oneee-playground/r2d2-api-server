package auth_module

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
)

//go:generate mockgen -source=token.go -destination=../../../test/mocks/token.go -package=mocks

var (
	ErrTokenInvalid = errors.New("token is not valid")
	ErrTokenExpired = errors.New("token is expired")
)

type Token struct {
	Payload   TokenPayload
	ExpiresAt time.Time

	Raw string
}

type TokenPayload struct {
	UserID uuid.UUID       `json:"userId"`
	Role   domain.UserRole `json:"role"`
}

type TokenManager interface {
	Issue(ctx context.Context, payload TokenPayload, exp time.Time) (Token, error)
	Decode(ctx context.Context, raw string) (Token, error)
}
