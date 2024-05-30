package loader

import (
	"context"
	"os"

	"github.com/oneee-playground/r2d2-api-server/internal/global/config"
)

type OSEnvLoader struct{}

var _ config.Loader = (*OSEnvLoader)(nil)

func NewOSEnvLoader() *OSEnvLoader {
	return &OSEnvLoader{}
}

func (l *OSEnvLoader) Fill(ctx context.Context, conf *config.Config) error {
	// Fill in the struct here.
	os.Getenv("example")

	return nil
}
