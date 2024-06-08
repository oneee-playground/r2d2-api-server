package task_module

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/domain/dto"
	"github.com/oneee-playground/r2d2-api-server/internal/global/status"
	"github.com/oneee-playground/r2d2-api-server/internal/global/tx"
	"github.com/pkg/errors"
)

type taskUsecase struct {
	lock tx.Locker

	taskRepository domain.TaskRepository
}

var _ domain.TaskUsecase = (*taskUsecase)(nil)

func NewTaskUsecase(tr domain.TaskRepository, l tx.Locker) *taskUsecase {
	return &taskUsecase{
		taskRepository: tr,
		lock:           l,
	}
}

func (u *taskUsecase) GetList(ctx context.Context) (out *dto.TaskListOutput, err error) {
	tasks, err := u.taskRepository.FetchAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "fetching all tasks")
	}

	return toTaskListOutput(tasks), nil
}

func (u *taskUsecase) GetTask(ctx context.Context, in dto.IDInput) (out *dto.TaskOutput, err error) {
	task, err := u.taskRepository.FetchByID(ctx, in.ID)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			return nil, status.NewErr(http.StatusNotFound, err.Error())
		}

		return nil, errors.Wrap(err, "fetching task by id")
	}

	return toTaskOutput(task), nil
}

func (u *taskUsecase) CreateTask(ctx context.Context, in dto.TaskInput) (out *dto.IDOutput, err error) {
	task := domain.Task{
		ID:          uuid.New(),
		Title:       in.Title,
		Description: in.Description,
		Stage:       domain.StageDraft,
	}

	if err := u.taskRepository.Create(ctx, task); err != nil {
		return nil, errors.Wrap(err, "creating task")
	}

	return toIDOutput(task), nil
}

func (u *taskUsecase) UpdateTask(ctx context.Context, in dto.UpdateTaskInput) (err error) {
	task, err := u.taskRepository.FetchByID(ctx, in.ID)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			return status.NewErr(http.StatusNotFound, err.Error())
		}

		return errors.Wrap(err, "fetching task by id")
	}

	task.Title = in.Title
	task.Description = in.Description

	if err := u.taskRepository.Update(ctx, task); err != nil {
		return errors.Wrap(err, "updating task")
	}

	return nil
}

func (u *taskUsecase) ChangeStage(ctx context.Context, in dto.TaskStageInput) (err error) {
	ctx, err = tx.NewAtomic(ctx, tx.AtomicOpts{
		ReadOnly:    false,
		DataSources: []any{u.taskRepository},
	})
	if err != nil {
		return errors.Wrap(err, "starting atomic transaction")
	}
	defer tx.Evaluate(ctx, &err)

	ctx, release, err := u.lock.Acquire(ctx, "task", in.ID.String())
	if err != nil {
		return errors.Wrap(err, "acquiring lock")
	}
	defer release()

	task, err := u.taskRepository.FetchByID(ctx, in.ID)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			return status.NewErr(http.StatusNotFound, err.Error())
		}

		return errors.Wrap(err, "fetching task by id")
	}

	stage := domain.TaskStage(in.Stage)

	switch stage {
	case task.Stage:
		return status.NewErr(http.StatusForbidden, "redundant change")
	case domain.StageDraft:
		return status.NewErr(http.StatusForbidden, "cannot go back to draft")
	case domain.StageAvailable:
		// TODO: need to validate that task has at least and at most one primary resource.
	case domain.StageFixing:
		// TODO: Publish a event to cancel all running or queued submissions.
		_ = 1 + 1
	}

	task.Stage = stage

	if err := u.taskRepository.Update(ctx, task); err != nil {
		return errors.Wrap(err, "updating task")
	}

	return nil
}
