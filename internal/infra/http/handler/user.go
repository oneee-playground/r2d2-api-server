package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
)

type UserHandler struct {
	usecase domain.UserUsecase
}

func NewUserHandler(usecase domain.UserUsecase) *UserHandler {
	return &UserHandler{usecase: usecase}
}

func (h *UserHandler) HandleSelfInfo(c *gin.Context) {
	out, err := h.usecase.GetSelfInfo(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, out)
}
