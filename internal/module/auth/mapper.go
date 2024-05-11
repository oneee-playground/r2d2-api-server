package auth_module

import "github.com/oneee-playground/r2d2-api-server/internal/domain/dto"

func toAccessTokenOutput(token Token) *dto.AccessTokenOutput {
	return &dto.AccessTokenOutput{Token: token.Raw}
}
