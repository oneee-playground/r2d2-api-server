package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/domain/dto"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/http/util"
)

type SubmissionHandler struct {
	usecase domain.SubmissionUsecase
}

func NewSubmissionHandler(usecase domain.SubmissionUsecase) *SubmissionHandler {
	return &SubmissionHandler{usecase: usecase}
}

func (h *SubmissionHandler) HandleGetList(c *gin.Context) {
	var in dto.SubmissionListInput

	if err := c.ShouldBindUri(&in.IDInput); err != nil {
		c.Error(util.WrapWithBadRequest(err))
		return
	}

	if err := c.ShouldBindQuery(&in); err != nil {
		c.Error(util.WrapWithBadRequest(err))
		return
	}

	out, err := h.usecase.GetList(c.Request.Context(), in)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, out)
}

func (h *SubmissionHandler) HandleSubmit(c *gin.Context) {
	var in dto.SubmissionInput

	if err := c.ShouldBindUri(&in.IDInput); err != nil {
		c.Error(util.WrapWithBadRequest(err))
		return
	}

	if err := c.ShouldBindJSON(&in); err != nil {
		c.Error(util.WrapWithBadRequest(err))
		return
	}

	out, err := h.usecase.Submit(c.Request.Context(), in)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, out)
}

func (h *SubmissionHandler) HandleDecideApproval(c *gin.Context) {
	var in dto.SubmissionDecisionInput

	if err := c.ShouldBindUri(&in.SubmissionIDInput); err != nil {
		c.Error(util.WrapWithBadRequest(err))
		return
	}

	if err := c.ShouldBindJSON(&in); err != nil {
		c.Error(util.WrapWithBadRequest(err))
		return
	}

	err := h.usecase.DecideApproval(c.Request.Context(), in)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}
func (h *SubmissionHandler) HandleCancel(c *gin.Context) {
	var in dto.SubmissionIDInput

	if err := c.ShouldBindUri(&in); err != nil {
		c.Error(util.WrapWithBadRequest(err))
		return
	}

	err := h.usecase.Cancel(c.Request.Context(), in)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}
