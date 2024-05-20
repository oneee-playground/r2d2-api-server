package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/oneee-playground/r2d2-api-server/internal/domain/dto"
)

//go:generate mockgen -source=event.go -destination=../../test/mocks/event.go -package=mocks

type EventKind string

const (
	KindSubmit       EventKind = "SUBMIT"
	KindApprove      EventKind = "APPROVE"
	KindReject       EventKind = "REJECT"
	KindBuildStart   EventKind = "BUILD_START"
	KindBuildFail    EventKind = "BUILD_FAIL"
	KindBuildSuccess EventKind = "BUILD_SUCCESS"
	KindQueue        EventKind = "QUEUE"
	KindTestStart    EventKind = "TEST_START"
	KindTestFail     EventKind = "TEST_FAIL"
	KindTestSuccess  EventKind = "TEST_SUCCESS"
	KindCancel       EventKind = "CANCEL"
)

type Event struct {
	ID uuid.UUID

	Kind      EventKind
	Extra     string
	Timestamp time.Time

	SubmissionID uuid.UUID
	Submission   *Submission
}

type EventUsecase interface {
	GetAllFromSubmission(ctx context.Context, in dto.SubmissionIDInput) (out *dto.EventListOutput, err error)
}

type EventRepository interface {
	FetchAllBySubmissionID(ctx context.Context, id uuid.UUID) ([]Event, error)
	Create(ctx context.Context, event Event) error
}
