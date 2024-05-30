package config

import (
	"context"

	"github.com/pkg/errors"
)

var loaded Config

type Loader interface {
	Fill(ctx context.Context, conf *Config) error
}

func Load(ctx context.Context, loader Loader) error {
	var conf Config
	if err := loader.Fill(ctx, &conf); err != nil {
		return errors.Wrap(err, "filling configs")
	}

	loaded = conf

	return nil
}
