package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/domain/dto"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/http/util"
)

type SectionHandler struct {
	usecase domain.SectionUsecase
}

func (h *SectionHandler) HandleGetList(c *gin.Context) {
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

func (h *SectionHandler) HandleCreateSection(c *gin.Context) {
	var in dto.CreateSectionInput

	if err := c.ShouldBindUri(&in.IDInput); err != nil {
		c.Error(util.WrapWithBadRequest(err))
		return
	}

	if err := c.ShouldBindJSON(&in.SectionInput); err != nil {
		c.Error(util.WrapWithBadRequest(err))
		return
	}

	err := h.usecase.CreateSection(c.Request.Context(), in)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusCreated)
}

func (h *SectionHandler) HandleUpdateSeciton(c *gin.Context) {
	var in dto.UpdateSectionInput

	if err := c.ShouldBindUri(&in.SectionIDInput); err != nil {
		c.Error(util.WrapWithBadRequest(err))
		return
	}

	if err := c.ShouldBindJSON(&in.SectionInput); err != nil {
		c.Error(util.WrapWithBadRequest(err))
		return
	}

	err := h.usecase.UpdateSection(c.Request.Context(), in)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *SectionHandler) HandleChangeIndex(c *gin.Context) {
	var in dto.SectionIndexInput

	if err := c.ShouldBindUri(&in.SectionIDInput); err != nil {
		c.Error(util.WrapWithBadRequest(err))
		return
	}

	if err := c.ShouldBindJSON(&in); err != nil {
		c.Error(util.WrapWithBadRequest(err))
		return
	}

	err := h.usecase.ChangeIndex(c.Request.Context(), in)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}
