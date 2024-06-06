package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/domain/dto"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/http/util"
)

type ResourceHandler struct {
	usecase domain.ResourceUsecase
}

func (h *ResourceHandler) HandleGetList(c *gin.Context) {
	var in dto.IDInput

	if err := c.ShouldBindUri(&in); err != nil {
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

func (h *ResourceHandler) HandleCreateResource(c *gin.Context) {
	var in dto.CreateResourceInput

	if err := c.ShouldBindUri(&in.IDInput); err != nil {
		c.Error(util.WrapWithBadRequest(err))
		return
	}

	if err := c.ShouldBindJSON(&in.ResourceInput); err != nil {
		c.Error(util.WrapWithBadRequest(err))
		return
	}

	err := h.usecase.CreateResource(c.Request.Context(), in)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusCreated)
}

func (h *ResourceHandler) HandleDeleteResource(c *gin.Context) {
	var in dto.ResourceIDInput

	if err := c.ShouldBindUri(&in.IDInput); err != nil {
		c.Error(util.WrapWithBadRequest(err))
		return
	}

	if err := c.ShouldBindJSON(&in); err != nil {
		c.Error(util.WrapWithBadRequest(err))
		return
	}

	err := h.usecase.DeleteResource(c.Request.Context(), in)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
