package auth_module

import (
	"context"
	"errors"
	"time"

	"github.com/oneee-playground/r2d2-api-server/internal/global/auth"
)

//go:generate mockgen -source=token.go -destination=../../../test/mocks/token.go -package=mocks

var (
	ErrTokenInvalid = errors.New("token is not valid")
	ErrTokenExpired = errors.New("token is expired")
)

type Token struct {
	Payload   auth.Payload
	ExpiresAt time.Time

	Raw string
}

type TokenIssuer interface {
	Issue(ctx context.Context, payload auth.Payload, exp time.Time) (Token, error)
}

type TokenDecoder interface {
	Decode(ctx context.Context, raw string) (Token, error)
}
