package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oneee-playground/r2d2-api-server/internal/global/status"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	errInternalError = status.NewErr(http.StatusInternalServerError, "internal server error")
)

type ErrorHandler struct {
	logger *zap.Logger
}

func NewErrorHandler(logger *zap.Logger) *ErrorHandler {
	return &ErrorHandler{logger: logger}
}

func (h *ErrorHandler) CatchWithStatusCode(c *gin.Context) {
	c.Next()

	switch len(c.Errors) {
	case 0:
		return
	case 1:
		err := c.Errors[0].Err
		sErr, ok := err.(status.Error)
		if !ok {
			h.logger.Error("unhandled error", zap.Error(err))

			sErr = errInternalError
		}

		c.AbortWithStatusJSON(sErr.StatusCode,
			gin.H{
				"message": sErr.Message,
			},
		)
		return
	}

	panic(errors.Errorf("handler resulted in multiple errors: %v", c.Errors))
}

func (h *ErrorHandler) RecoverPanic(c *gin.Context) {
	c.Next()

	if err := recover(); err != nil {
		h.logger.Error("server panicked", zap.Any("recovered", err))

		sErr := errInternalError
		c.AbortWithStatusJSON(sErr.StatusCode,
			gin.H{
				"message": sErr.Message,
			},
		)
		return
	}
}
