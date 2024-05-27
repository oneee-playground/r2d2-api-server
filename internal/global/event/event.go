package event

import (
	"time"

	"github.com/google/uuid"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
)

type Topic string

const (
	TopicSubmission Topic = "submission"
)

// Event schema for TopicSubmission
type SubmissionEvent struct {
	ID           uuid.UUID
	SubmissionID uuid.UUID
	UserID       uuid.UUID

	Kind      domain.EventKind
	Extra     string
	Timestamp time.Time
}
