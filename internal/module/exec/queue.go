package exec_module

import (
	"context"
)

//go:generate mockgen -source=queue.go -destination=../../../test/mocks/queue.go -package=mocks

type JobQueue interface {
	Append(ctx context.Context, job *Job) error
}
