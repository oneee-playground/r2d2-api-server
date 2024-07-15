package user_module

import (
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/domain/dto"
)

func toUserInfo(user domain.User) *dto.UserInfo {
	return &dto.UserInfo{
		ID:         user.ID.String(),
		Username:   user.Username,
		ProfileURL: user.ProfileURL,
		Role:       user.Role.String(),
	}
}
