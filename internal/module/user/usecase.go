package user_module

import (
	"context"

	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/domain/dto"
	"github.com/oneee-playground/r2d2-api-server/internal/global/auth"
	"github.com/pkg/errors"
)

type userUsecase struct {
	userRepository domain.UserRepository
}

var _ domain.UserUsecase = (*userUsecase)(nil)

func NewUserUsecase(ur domain.UserRepository) *userUsecase {
	return &userUsecase{
		userRepository: ur,
	}
}

func (u *userUsecase) GetSelfInfo(ctx context.Context) (out *dto.UserInfo, err error) {
	// Assume anonymous users are filtered with middleware.
	info := auth.MustExtract(ctx)

	user, err := u.userRepository.FetchByID(ctx, info.UserID)
	if err != nil {
		return nil, errors.Wrap(err, "fetching user by id")
	}

	return toUserInfo(user), nil
}
