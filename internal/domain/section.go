package domain

import (
	"github.com/google/uuid"
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
	Type  SectionType

	TaskID uuid.UUID
	Task   *Task
}
