package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/domain/dto"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/http/util"
)

type TaskHandler struct {
	usecase domain.TaskUsecase
}

func (h *TaskHandler) HandleGetList(c *gin.Context) {
	out, err := h.usecase.GetList(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, out)
}

func (h *TaskHandler) HandleGetTask(c *gin.Context) {
	var in dto.IDInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.Error(util.WrapWithBadRequest(err))
		return
	}

	out, err := h.usecase.GetTask(c.Request.Context(), in)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, out)
}

func (h *TaskHandler) HandleCreateTask(c *gin.Context) {
	var in dto.TaskInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.Error(util.WrapWithBadRequest(err))
		return
	}

	out, err := h.usecase.CreateTask(c.Request.Context(), in)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, out)
}

func (h *TaskHandler) HandleUpdateTask(c *gin.Context) {
	var in dto.UpdateTaskInput

	if err := c.ShouldBindUri(&in.IDInput); err != nil {
		c.Error(util.WrapWithBadRequest(err))
		return
	}

	if err := c.ShouldBindJSON(&in.TaskInput); err != nil {
		c.Error(util.WrapWithBadRequest(err))
		return
	}

	err := h.usecase.UpdateTask(c.Request.Context(), in)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusCreated)
}

func (h *TaskHandler) HandleChangeStage(c *gin.Context) {
	var in dto.TaskStageInput

	if err := c.ShouldBindUri(&in.IDInput); err != nil {
		c.Error(util.WrapWithBadRequest(err))
		return
	}

	if err := c.ShouldBindJSON(&in); err != nil {
		c.Error(util.WrapWithBadRequest(err))
		return
	}

	err := h.usecase.ChangeStage(c.Request.Context(), in)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}
