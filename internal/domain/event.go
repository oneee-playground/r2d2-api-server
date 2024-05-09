package domain

import "github.com/google/uuid"

type EventKind uint8

const (
	KindSubmit EventKind = iota
	KindApprove
	KindReject
	KindBuildStart
	KindBuildFail
	KindBuildSuccess
	KindQueue
	KindTestStart
	KindTestFail
	KindTestSuccess
	KindCancel
)

type Event struct {
	ID uuid.UUID

	Kind  EventKind
	Extra string

	SubmissionID uuid.UUID
	Submission   *Submission
}
