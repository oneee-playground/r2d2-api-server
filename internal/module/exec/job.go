package exec_module

import (
	"github.com/google/uuid"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
)

type Resource struct {
	Image     string  `json:"image"`
	Name      string  `json:"name"`
	Port      uint16  `json:"port"`
	CPU       float64 `json:"cpu"`
	Memory    uint64  `json:"memory"`
	IsPrimary bool    `json:"isPrimary"`
}

type Section struct {
	ID   uuid.UUID          `json:"id"`
	Type domain.SectionType `json:"type"`
}

type Submission struct {
	ID         uuid.UUID `json:"id"`
	Repository string    `json:"repositoy"`
	CommitHash string    `json:"commitHash"`
}

type Job struct {
	TaskID uuid.UUID `json:"taskID"`

	Resources []Resource `json:"resources"`
	Sections  []Section  `json:"sections"`

	Submission Submission `json:"submission"`
}
