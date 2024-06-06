package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/domain/dto"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/http/util"
)

type AuthHandler struct {
	usecase domain.AuthUsecase
}

func (h *AuthHandler) HandleSignIn(c *gin.Context) {
	var in dto.SignInInput

	if err := c.ShouldBindJSON(&in); err != nil {
		c.Error(util.WrapWithBadRequest(err))
		return
	}

	out, err := h.usecase.SignIn(c.Request.Context(), in)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, out)
}
