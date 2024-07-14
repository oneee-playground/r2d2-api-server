package dto

import (
	"time"

	"github.com/google/uuid"
)

type SubmissionPaginator struct {
	Offset int `form:"offset" binding:"required"`
}

type SubmissionListInput struct {
	IDInput
	SubmissionPaginator
}

type SubmissionListElem struct {
	ID        uuid.UUID `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	SourceURL string    `json:"sourceURL"`
	IsDone    bool     `json:"isDone"`
	User      UserInfo  `json:"user"`
}

type SubmissionListOutput []SubmissionListElem

type SubmissionInput struct {
	IDInput
	Repository string `json:"repository" binding:"required"`
	CommitHash string `json:"commitHash" binding:"required"`
}

type SubmissionIDInput struct {
	TaskID       uuid.UUID `uri:"id" binding:"required"`
	SubmissionID uuid.UUID `uri:"submissionID" binding:"required"`
}

type SubmissionDecisionInput struct {
	SubmissionIDInput
	Action string `json:"action" binding:"required" validate:"submission_action"`
	Extra  string `json:"extra" binding:"required"`
}
