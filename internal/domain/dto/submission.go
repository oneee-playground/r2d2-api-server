package dto

import (
	"time"
)

type SubmissionPaginator struct {
	Offset int `form:"offset"`
}

type SubmissionListInput struct {
	IDInput
	SubmissionPaginator
}

type SubmissionListElem struct {
	ID        string    `json:"id" binding:"uuid"`
	Timestamp time.Time `json:"timestamp"`
	SourceURL string    `json:"sourceURL"`
	IsDone    bool      `json:"isDone"`
	User      UserInfo  `json:"user"`
}

type SubmissionListOutput []SubmissionListElem

type SubmissionInput struct {
	IDInput
	Repository string `json:"repository" binding:"required"`
	CommitHash string `json:"commitHash" binding:"required"`
}

type SubmissionIDInput struct {
	TaskID       string `uri:"id" binding:"required,uuid"`
	SubmissionID string `uri:"submissionID" binding:"required,uuid"`
}

type SubmissionDecisionInput struct {
	SubmissionIDInput
	Action string `json:"action" binding:"required" validate:"submission_action"`
	Extra  string `json:"extra" binding:"required"`
}
