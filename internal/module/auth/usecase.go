package auth_module

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/domain/dto"
	"github.com/oneee-playground/r2d2-api-server/internal/global/auth"
	"github.com/oneee-playground/r2d2-api-server/internal/global/status"
	"github.com/oneee-playground/r2d2-api-server/internal/global/tx"
	"github.com/pkg/errors"
)

type authUsecase struct {
	oauth       OAuthClient
	tokenIssuer TokenIssuer
	lock        tx.Locker

	userRepository domain.UserRepository
}

var _ domain.AuthUsecase = (*authUsecase)(nil)

func NewAuthUsecase(oa OAuthClient, ti TokenIssuer, ur domain.UserRepository, l tx.Locker) *authUsecase {
	return &authUsecase{
		oauth:          oa,
		tokenIssuer:    ti,
		userRepository: ur,
		lock:           l,
	}
}

func (uc *authUsecase) SignIn(ctx context.Context, in dto.SignInInput) (out *dto.AccessTokenOutput, err error) {
	token, err := uc.oauth.IssueAccessToken(ctx, in.Code)
	if err != nil {
		if errors.Is(err, ErrInvalidCode) || errors.Is(err, ErrNotEnoughScope) {
			return nil, status.NewErr(http.StatusBadRequest, err.Error())
		}

		return nil, errors.Wrap(err, "issuing access token")
	}

	user, err := uc.oauth.GetUserInfo(ctx, token)
	if err != nil {
		return nil, errors.Wrap(err, "getting user info")
	}

	ctx, err = tx.NewAtomic(ctx, tx.AtomicOpts{
		ReadOnly:    false,
		DataSources: []any{uc.userRepository},
	})
	if err != nil {
		return nil, errors.Wrap(err, "starting atomic transaction")
	}
	defer tx.Evaluate(ctx, &err)

	ctx, release, err := uc.lock.Acquire(ctx, "username", user.Username)
	if err != nil {
		return nil, errors.Wrap(err, "acquiring lock")
	}
	defer release()

	ok, err := uc.userRepository.UsernameExists(ctx, user.Username)
	if err != nil {
		return nil, errors.Wrap(err, "checking if username exists")
	}

	if ok {
		user, err = uc.userRepository.FetchByUsername(ctx, user.Username)
		if err != nil {
			return nil, errors.Wrap(err, "fetching user with username")
		}
	} else {
		user.ID = uuid.New()
		user.Role = domain.RoleMember

		if err := uc.userRepository.Create(ctx, user); err != nil {
			return nil, errors.Wrap(err, "creating user")
		}
	}

	payload := auth.Payload{
		UserID: user.ID,
		Role:   user.Role,
	}

	// TODO: change this exp into constant one
	accessToken, err := uc.tokenIssuer.Issue(ctx, payload, time.Now().Add(7*time.Hour))
	if err != nil {
		return nil, errors.Wrap(err, "issuing token")
	}

	return toAccessTokenOutput(accessToken), nil
}
