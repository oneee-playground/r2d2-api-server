package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/oneee-playground/r2d2-api-server/internal/domain/dto"
	"github.com/pkg/errors"
)

//go:generate mockgen -source=resource.go -destination=../../test/mocks/resource.go -package=mocks

type Resource struct {
	Image  string
	Name   string
	Port   uint16
	CPU    float64
	Memory uint64
	// IsPrimary sets if the given resource is primary (user will implement their process on it).
	// Only one of the resources of same task should be primary.
	IsPrimary bool

	TaskID uuid.UUID
	Task   *Task
}

type ResourceUsecase interface {
	GetList(ctx context.Context, in dto.IDInput) (out *dto.ResourceListOutput, err error)
	CreateResource(ctx context.Context, in dto.CreateResourceInput) (err error)
	DeleteResource(ctx context.Context, in dto.ResourceIDInput) (err error)
}

var (
	ErrResourceNotFound  = errors.New("resource not found")
	ErrDuplicateResource = errors.New("duplicate resource")
)

type ResourceRepository interface {
	FetchAllByTaskID(ctx context.Context, taskID uuid.UUID) ([]Resource, error)
	Create(ctx context.Context, resource Resource) error
	Delete(ctx context.Context, taskID uuid.UUID, name string) error
}
