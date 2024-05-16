package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
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

type EventRepository interface {
	Create(ctx context.Context, event Event) error
}
