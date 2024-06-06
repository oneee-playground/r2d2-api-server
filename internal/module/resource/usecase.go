package resource_module

import (
	"context"
	"net/http"

	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/domain/dto"
	"github.com/oneee-playground/r2d2-api-server/internal/global/status"
	"github.com/pkg/errors"
)

type resourceUsecase struct {
	taskRepository     domain.TaskRepository
	resourceRepository domain.ResourceRepository
}

var _ domain.ResourceUsecase = (*resourceUsecase)(nil)

func NewResourceUsecase(rr domain.ResourceRepository, tr domain.TaskRepository) *resourceUsecase {
	return &resourceUsecase{
		resourceRepository: rr,
		taskRepository:     tr,
	}
}

func (u *resourceUsecase) GetList(ctx context.Context, in dto.IDInput) (out *dto.ResourceListOutput, err error) {
	exists, err := u.taskRepository.ExistsByID(ctx, in.ID)
	if err != nil {
		return nil, errors.Wrap(err, "checking task exists")
	}

	if !exists {
		return nil, status.NewErr(http.StatusNotFound, "task not found")
	}

	resources, err := u.resourceRepository.FetchAllByTaskID(ctx, in.ID)
	if err != nil {
		return nil, errors.Wrap(err, "fetching resources")
	}

	return toResourceListOutput(resources), nil
}

func (u *resourceUsecase) CreateResource(ctx context.Context, in dto.CreateResourceInput) (err error) {
	// TODO: transactional
	task, err := u.taskRepository.FetchByID(ctx, in.ID)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			return status.NewErr(http.StatusNotFound, err.Error())
		}

		return err
	}

	if task.Stage == domain.StageAvailable {
		return status.NewErr(http.StatusForbidden, "cannot create resource on available task")
	}

	resource := domain.Resource{
		Image:     in.Image,
		Name:      in.Name,
		Port:      in.Port,
		CPU:       in.CPU,
		Memory:    in.Memory,
		IsPrimary: *in.IsPrimary,
		TaskID:    task.ID,
	}

	if err := u.resourceRepository.Create(ctx, resource); err != nil {
		if errors.Is(err, domain.ErrDuplicateResource) {
			return status.NewErr(http.StatusConflict, err.Error())
		}

		return errors.Wrap(err, "creating resource")
	}

	return nil
}

func (u *resourceUsecase) DeleteResource(ctx context.Context, in dto.ResourceIDInput) (err error) {
	// TODO: transactional
	task, err := u.taskRepository.FetchByID(ctx, in.ID)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			return status.NewErr(http.StatusNotFound, err.Error())
		}

		return err
	}

	if task.Stage == domain.StageAvailable {
		return status.NewErr(http.StatusForbidden, "cannot delete resource on available task")
	}

	if err := u.resourceRepository.Delete(ctx, task.ID, in.Name); err != nil {
		if errors.Is(err, domain.ErrResourceNotFound) {
			return status.NewErr(http.StatusNotFound, err.Error())
		}

		return errors.Wrap(err, "deleting resource")
	}

	return nil
}
