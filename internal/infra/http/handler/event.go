package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/domain/dto"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/http/util"
)

type EventHandler struct {
	usecase domain.EventUsecase
}

func NewEventHandler(usecase domain.EventUsecase) *EventHandler {
	return &EventHandler{usecase: usecase}
}

func (h *EventHandler) HandleGetAll(c *gin.Context) {
	var in dto.SubmissionIDInput

	if err := c.ShouldBindUri(&in); err != nil {
		c.Error(util.WrapWithBadRequest(err))
		return
	}

	out, err := h.usecase.GetAllFromSubmission(c.Request.Context(), in)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, out)
}
