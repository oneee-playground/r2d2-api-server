package util

import (
	"net/http"

	"github.com/oneee-playground/r2d2-api-server/internal/global/status"
)

func WrapWithBadRequest(err error) error {
	return status.NewErr(http.StatusBadRequest, err.Error())
}
