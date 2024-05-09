package domain

import "github.com/google/uuid"

type SectionType uint8

const (
	TypeScenario SectionType = iota
	TypeLoad
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
