package domain

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/oneee-playground/r2d2-api-server/internal/domain/dto"
)

//go:generate mockgen -source=task.go -destination=../../test/mocks/task.go -package=mocks

type TaskStage string

const (
	StageDraft     TaskStage = "DRAFT"
	StageAvailable TaskStage = "AVAILABLE"
	StageFixing    TaskStage = "FIXING"
)

type Task struct {
	ID          uuid.UUID
	Title       string
	Description string
	Stage       TaskStage
}

type TaskUsecase interface {
	GetList(ctx context.Context) (out *dto.TaskListOutput, err error)
	GetTask(ctx context.Context, in dto.IDInput) (out *dto.TaskOutput, err error)
	CreateTask(ctx context.Context, in dto.TaskInput) (out *dto.IDOutput, err error)
	UpdateTask(ctx context.Context, in dto.UpdateTaskInput) (err error)
	ChangeStage(ctx context.Context, in dto.TaskStageInput) (err error)
}

var (
	ErrTaskNotFound = errors.New("task not found")
)

type TaskRepository interface {
	ExistsByID(ctx context.Context, id uuid.UUID) (bool, error)
	FetchAll(ctx context.Context) ([]Task, error)
	FetchByID(ctx context.Context, id uuid.UUID) (Task, error)
	Create(ctx context.Context, task Task) error
	Update(ctx context.Context, task Task) error
}
