package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/global/auth"
	"github.com/oneee-playground/r2d2-api-server/internal/global/status"
	auth_module "github.com/oneee-playground/r2d2-api-server/internal/module/auth"
)

var (
	errInvalidAuthToken = status.NewErr(http.StatusUnauthorized, "invalid auth token")
	errNoPermission     = status.NewErr(http.StatusForbidden, "permission denied")
)

type AuthFilter struct {
	tokenDecoder auth_module.TokenDecoder
}

func NewAuthFilter(td auth_module.TokenDecoder) *AuthFilter {
	return &AuthFilter{tokenDecoder: td}
}

func (f *AuthFilter) Required(required bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, err := f.extractPayload(c)
		if err != nil {
			if required {
				c.Error(errInvalidAuthToken)
				c.Abort()
				return
			}
		} else {
			c.Request = c.Request.WithContext(auth.Inject(c, payload))
		}

		c.Next()
	}
}

func (f *AuthFilter) extractPayload(c *gin.Context) (auth.Payload, error) {
	authorization, ok := c.Request.Header["Authorization"]
	if !ok || len(authorization) != 1 {
		return auth.Payload{}, errInvalidAuthToken
	}

	bearerToken, found := strings.CutPrefix(authorization[0], "Bearer ")
	if !found {
		return auth.Payload{}, errInvalidAuthToken
	}

	token, err := f.tokenDecoder.Decode(c, bearerToken)
	if err != nil {
		return auth.Payload{}, errInvalidAuthToken
	}

	return token.Payload, nil
}

func (f *AuthFilter) AtLeast(role domain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		info := auth.MustExtract(c.Request.Context())

		if info.Role < role {
			c.Error(errNoPermission)
			c.Abort()
			return
		}

		c.Next()
	}
}
