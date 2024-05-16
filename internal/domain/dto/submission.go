package dto

import (
	"time"

	"github.com/google/uuid"
)

type SubmissionPaginator struct {
	Offset int
}

type SubmissionListInput struct {
	IDInput
	SubmissionPaginator
}

type SubmissionListElem struct {
	ID        uuid.UUID `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	SourceURL string    `json:"sourceURL"`
	IsDone    bool      `json:"isDone"`
	User      UserInfo  `json:"user"`
}

type SubmissionListOutput []SubmissionListElem

type SubmissionInput struct {
	IDInput
	Repository string `json:"repository"`
	CommitHash string `json:"commitHash"`
}

type SubmissionIDInput struct {
	TaskID       uuid.UUID
	SubmissionID uuid.UUID
}

type SubmissionDecisionInput struct {
	SubmissionIDInput
	Action string `json:"action" validate:"submission_action"`
	Extra  string `json:"extra"`
}
