package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/global/tx"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/data/ent/datasource"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/data/ent/model"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/data/ent/model/task"
)

type TaskRepository struct {
	*datasource.DataSource
}

var (
	_ domain.TaskRepository = (*TaskRepository)(nil)
	_ tx.DataSource         = (*TaskRepository)(nil)
)

func NewTaskRepository(ds *datasource.DataSource) *TaskRepository {
	return &TaskRepository{DataSource: ds}
}

func (r *TaskRepository) Create(ctx context.Context, task domain.Task) error {
	return r.DataSource.TxOrPlain(ctx).Task.
		Create().
		SetID(task.ID).
		SetTitle(task.Title).
		SetDescription(task.Description).
		SetStage(string(task.Stage)).
		Exec(ctx)
}

func (r *TaskRepository) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	return r.DataSource.TxOrPlain(ctx).Task.
		Query().
		Where(task.ID(id)).
		Exist(ctx)
}

func (r *TaskRepository) FetchAll(ctx context.Context) ([]domain.Task, error) {
	models, err := r.DataSource.TxOrPlain(ctx).Task.Query().All(ctx)
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
	entity, err := r.DataSource.TxOrPlain(ctx).Task.Get(ctx, id)
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
	return r.DataSource.TxOrPlain(ctx).Task.
		UpdateOneID(task.ID).
		SetTitle(task.Title).
		SetDescription(task.Description).
		SetStage(string(task.Stage)).
		Exec(ctx)
}
