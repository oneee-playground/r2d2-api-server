package exec_module

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

//go:generate mockgen -source=context.go -destination=../../../test/mocks/context.go -package=mocks

type ExecContext struct {
	TaskID uuid.UUID `json:"taskID"`
	UserID uuid.UUID `json:"userID"`

	Repository string `json:"repositoy"`
	CommitHash string `json:"commitHash"`
}

var (
	ErrContextNotFound = errors.New("context not found")
)

type ExecContextStroage interface {
	Get(ctx context.Context, submissionID uuid.UUID) (ExecContext, error)
	Set(ctx context.Context, submissionID uuid.UUID, execCtx ExecContext) error
	Delete(ctx context.Context, submissionID uuid.UUID) error
}
