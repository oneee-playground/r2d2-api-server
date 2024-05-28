package event

import (
	"time"

	"github.com/google/uuid"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
)

type Topic string

const (
	TopicSubmission Topic = "submission"
	TopicBuild      Topic = "build"
	TopicTest       Topic = "test"
)

// Event schema for TopicSubmission
type SubmissionEvent struct {
	ID           uuid.UUID        `json:"id"`
	SubmissionID uuid.UUID        `json:"submissionID"`
	UserID       uuid.UUID        `json:"userID"`
	Kind         domain.EventKind `json:"kind"`
	Extra        string           `json:"extra"`
	Timestamp    time.Time        `json:"timestamp"`
}

// Event schema for TopicBuild, TopicTest
type ExecEvent struct {
	ID      uuid.UUID     `json:"id"`
	Success bool          `json:"success"`
	Took    time.Duration `json:"took"`
	Extra   string        `json:"extra"`
}
