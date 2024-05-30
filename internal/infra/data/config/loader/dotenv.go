package loader

import (
	"bufio"
	"bytes"
	"context"
	"os"

	"github.com/oneee-playground/r2d2-api-server/internal/global/config"
	"github.com/pkg/errors"
)

type DotEnvLoader struct {
	path string
}

var _ config.Loader = (*DotEnvLoader)(nil)

func NewDotEnvLoader(path string) *DotEnvLoader {
	return &DotEnvLoader{
		path: path,
	}
}

func (l *DotEnvLoader) Fill(ctx context.Context, conf *config.Config) error {
	file, err := os.Open(l.path)
	if err != nil {
		return errors.Wrap(err, "opening file from path")
	}

	m, err := l.parseToMap(ctx, bufio.NewScanner(file))
	if err != nil {
		return errors.Wrap(err, "parsing env file to map")
	}

	// Fill in the struct here.
	_ = m["exmaple"]

	return nil
}

func (l *DotEnvLoader) parseToMap(ctx context.Context, s *bufio.Scanner) (map[string]string, error) {
	m := make(map[string]string)

	for s.Scan() {
		select {
		case <-ctx.Done():
			return nil, errors.Wrap(ctx.Err(), "context canceled")
		default:
		}

		b := s.Bytes()
		key, val, found := bytes.Cut(b, []byte("="))
		if !found {
			return nil, errors.Errorf("malformed environment variable: %s", b)
		}

		m[string(key)] = string(val)
	}

	if s.Err() != nil {
		return nil, errors.Wrap(s.Err(), "error while scanniing file")
	}

	return m, nil
}
