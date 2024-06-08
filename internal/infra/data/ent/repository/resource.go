package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/global/tx"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/data/ent/datasource"
	resource "github.com/oneee-playground/r2d2-api-server/internal/infra/data/ent/model/resource"
)

type ResourceRepository struct {
	*datasource.DataSource
}

var (
	_ domain.ResourceRepository = (*ResourceRepository)(nil)
	_ tx.DataSource             = (*ResourceRepository)(nil)
)

func NewResourceRepository(ds *datasource.DataSource) *ResourceRepository {
	return &ResourceRepository{DataSource: ds}
}

func (r *ResourceRepository) Create(ctx context.Context, res domain.Resource) error {
	client := r.DataSource.TxOrPlain(ctx)

	exists, err := client.Resource.
		Query().
		Where(
			resource.And(
				resource.Name(res.Name),
				resource.TaskID(res.TaskID),
			),
		).Exist(ctx)
	if err != nil {
		return err
	}

	if exists {
		return domain.ErrDuplicateResource
	}

	return client.Resource.
		Create().
		SetImage(res.Image).
		SetName(res.Name).
		SetPort(res.Port).
		SetIsPrimary(res.IsPrimary).
		SetCPU(res.CPU).
		SetMemory(res.Memory).
		SetTaskID(res.TaskID).
		Exec(ctx)
}

func (r *ResourceRepository) Delete(ctx context.Context, taskID uuid.UUID, name string) error {
	deleted, err := r.DataSource.TxOrPlain(ctx).Resource.
		Delete().
		Where(
			resource.And(
				resource.Name(name),
				resource.TaskID(taskID),
			),
		).Exec(ctx)
	if err != nil {
		return err
	}

	if deleted < 1 {
		return domain.ErrResourceNotFound
	}

	return nil
}

func (r *ResourceRepository) FetchAllByTaskID(ctx context.Context, taskID uuid.UUID) ([]domain.Resource, error) {
	models, err := r.DataSource.TxOrPlain(ctx).Resource.
		Query().
		Where(resource.TaskID(taskID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	resources := make([]domain.Resource, len(models))
	for idx, model := range models {
		resources[idx] = domain.Resource{
			Image:     model.Image,
			Name:      model.Name,
			Port:      model.Port,
			CPU:       model.CPU,
			Memory:    model.Memory,
			IsPrimary: model.IsPrimary,
			TaskID:    model.TaskID,
		}
	}

	return resources, nil
}
