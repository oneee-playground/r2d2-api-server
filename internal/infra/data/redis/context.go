package redis

import (
	"context"

	"github.com/google/uuid"
	exec_module "github.com/oneee-playground/r2d2-api-server/internal/module/exec"
	"github.com/redis/rueidis"
)

const _execContextKey = "exec-context"

type RedisExecContextStroage struct {
	client rueidis.Client
}

var _ exec_module.ExecContextStroage = (*RedisExecContextStroage)(nil)

func NewExecContextStroage(client rueidis.Client) *RedisExecContextStroage {
	return &RedisExecContextStroage{client: client}
}

func (s *RedisExecContextStroage) Get(ctx context.Context, submissionID uuid.UUID) (exec_module.ExecContext, error) {
	cmd := s.client.B().
		Get().
		Key(s.buildExecCtxKey(submissionID)).
		Build()

	var decoded exec_module.ExecContext
	if err := s.client.Do(ctx, cmd).DecodeJSON(&decoded); err != nil {
		return exec_module.ExecContext{}, err
	}

	return decoded, nil
}

func (s *RedisExecContextStroage) Set(ctx context.Context, submissionID uuid.UUID, execCtx exec_module.ExecContext) error {
	cmd := s.client.B().
		Set().
		Key(s.buildExecCtxKey(submissionID)).
		Value(rueidis.JSON(execCtx)).
		Build()

	return s.client.Do(ctx, cmd).Error()
}

func (s *RedisExecContextStroage) Delete(ctx context.Context, submissionID uuid.UUID) error {
	cmd := s.client.B().
		Del().
		Key(s.buildExecCtxKey(submissionID)).
		Build()

	return s.client.Do(ctx, cmd).Error()
}

func (s *RedisExecContextStroage) buildExecCtxKey(submissionID uuid.UUID) string {
	return buildKey(_execContextKey, submissionID.String())
}
