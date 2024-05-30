package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/data/ent/model"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/data/ent/model/task"
)

type TaskRepository struct {
	client *model.TaskClient
}

var _ domain.TaskRepository = (*TaskRepository)(nil)

func NewTaskRepository(client *model.TaskClient) *TaskRepository {
	return &TaskRepository{client: client}
}

func (r *TaskRepository) Create(ctx context.Context, task domain.Task) error {
	return r.client.Create().
		SetID(task.ID).
		SetTitle(task.Title).
		SetDescription(task.Description).
		SetStage(string(task.Stage)).
		Exec(ctx)
}

func (r *TaskRepository) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	return r.client.Query().
		Where(task.ID(id)).
		Exist(ctx)
}

func (r *TaskRepository) FetchAll(ctx context.Context) ([]domain.Task, error) {
	models, err := r.client.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	tasks := make([]domain.Task, len(models))
	for idx, model := range models {
		tasks[idx] = domain.Task{
			ID:          model.ID,
			Title:       model.Title,
			Description: model.Description,
			Stage:       domain.TaskStage(model.Stage),
		}
	}

	return tasks, nil
}

func (r *TaskRepository) FetchByID(ctx context.Context, id uuid.UUID) (domain.Task, error) {
	entity, err := r.client.Get(ctx, id)
	if err != nil {
		if model.IsNotFound(err) {
			return domain.Task{}, domain.ErrTaskNotFound
		}
		return domain.Task{}, err
	}

	task := domain.Task{
		ID:          entity.ID,
		Title:       entity.Title,
		Description: entity.Description,
		Stage:       domain.TaskStage(entity.Stage),
	}

	return task, nil
}

func (r *TaskRepository) Update(ctx context.Context, task domain.Task) error {
	return r.client.UpdateOneID(task.ID).
		SetTitle(task.Title).
		SetDescription(task.Description).
		SetStage(string(task.Stage)).
		Exec(ctx)
}
