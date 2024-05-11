package auth_module

import (
	"context"

	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/pkg/errors"
)

//go:generate mockgen -source=oauth.go -destination=../../../test/mocks/oauth.go -package=mocks

var (
	ErrInvalidCode    = errors.New("given code is not valid")
	ErrNotEnoughScope = errors.New("given scope is not enough")
)

type OAuthClient interface {
	// IssueAccessToken generates access token with given code.
	IssueAccessToken(ctx context.Context, code string) (string, error)
	// GetUserInfo gets user information with given token.
	// returned user should have its fields filled. except id and role.
	GetUserInfo(ctx context.Context, token string) (domain.User, error)
}
