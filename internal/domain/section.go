package domain

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/oneee-playground/r2d2-api-server/internal/domain/dto"
)

//go:generate mockgen -source=section.go -destination=../../test/mocks/section.go -package=mocks

type SectionType string

const (
	TypeScenario SectionType = "SCENARIO"
	TypeLoad     SectionType = "LOAD"
)

type Section struct {
	ID          uuid.UUID
	Title       string
	Description string
	// Index holds the order of section in the task.
	Index uint8
	// RPM stands for 'requests per minute'.
	// It will be non-zero when section's type is load testing.
	RPM  uint64
	Type SectionType
	// Example format of request-response.
	// It is json-formatted string.
	// TODO: Find a better way to represent this.
	Example string

	TaskID uuid.UUID
	Task   *Task
}

type SectionUsecase interface {
	GetList(ctx context.Context, in dto.IDInput) (out *dto.SectionListOutput, err error)
	CreateSection(ctx context.Context, in dto.CreateSectionInput) (err error)
	UpdateSection(ctx context.Context, in dto.UpdateSectionInput) (err error)
	ChangeIndex(ctx context.Context, in dto.SectionIndexInput) (err error)
}

type FetchSectionsOption struct {
	IncludeContent bool
}

var (
	ErrSectionNotFound = errors.New("section not found")
)

type SectionRepository interface {
	// FetchAllByTaskID fetches all the sections associated with task.
	// It always orders sections with its index.
	FetchAllByTaskID(ctx context.Context, taskID uuid.UUID, opt FetchSectionsOption) ([]Section, error)
	CountByTaskID(ctx context.Context, taskID uuid.UUID) (uint8, error)
	FetchByID(ctx context.Context, id uuid.UUID) (Section, error)
	Create(ctx context.Context, section Section) error
	Update(ctx context.Context, section Section) error
	// SaveIndexes saves index of given sections.
	// It only saves its field index. It is not affected by order of sections.
	SaveIndexes(ctx context.Context, sections []Section) error
}
