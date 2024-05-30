package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/oneee-playground/r2d2-api-server/internal/domain/dto"
)

//go:generate mockgen -source=submission.go -destination=../../test/mocks/submission.go -package=mocks

type SubmissionAction string

const (
	ActionApprove SubmissionAction = "APPROVE"
	ActionReject  SubmissionAction = "REJECT"
)

type Submission struct {
	ID        uuid.UUID
	Timestamp time.Time
	IsDone    bool

	// Github Repository name, e.g. "oneee-playground/empty"
	Repository string
	// Commit hash of the source.
	CommitHash string

	UserID uuid.UUID
	User   *User

	TaskID uuid.UUID
	Task   *Task
}

type SubmissionUsecase interface {
	GetList(ctx context.Context, in dto.SubmissionListInput) (out *dto.SubmissionListOutput, err error)
	Submit(ctx context.Context, in dto.SubmissionInput) (out *dto.IDOutput, err error)
	DecideApproval(ctx context.Context, in dto.SubmissionDecisionInput) (err error)
	Cancel(ctx context.Context, in dto.SubmissionIDInput) (err error)
}

var (
	ErrSubmissionNotFound = errors.New("submission not found")
)

type SubmissionRepository interface {
	// FetchPaginated returns list of submissions with given offset and limit.
	// It is ordered by timestamp desc.
	// Submissions will include User field.
	FetchPaginated(ctx context.Context, taskID uuid.UUID, offset, limit int) ([]Submission, error)
	Create(ctx context.Context, submission Submission) error
	Update(ctx context.Context, submission Submission) error
	UndoneExists(ctx context.Context, taskID, userID uuid.UUID) (bool, error)
	FetchByID(ctx context.Context, id uuid.UUID) (Submission, error)
}
