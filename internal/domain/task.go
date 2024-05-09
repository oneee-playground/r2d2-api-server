package domain

import "github.com/google/uuid"

type TaskStage uint8

const (
	StageDraft TaskStage = iota
	StageAvailable
	StageFixing
)

type Task struct {
	ID          uuid.UUID
	Title       string
	Description string
	Stage       TaskStage

	Sections []*Section
}
