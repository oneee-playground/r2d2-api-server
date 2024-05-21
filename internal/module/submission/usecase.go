package submission_module

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/domain/dto"
	"github.com/oneee-playground/r2d2-api-server/internal/global/auth"
	"github.com/oneee-playground/r2d2-api-server/internal/global/event"
	"github.com/oneee-playground/r2d2-api-server/internal/global/status"
	"github.com/pkg/errors"
)

type submissionUsecase struct {
	taskRepository       domain.TaskRepository
	submissionRepository domain.SubmissionRepository
	eventRepository      domain.EventRepository
	eventPublisher       event.Publisher
}

var _ domain.SubmissionUsecase = (*submissionUsecase)(nil)

func NewSubmissionUsecase(
	tr domain.TaskRepository, sr domain.SubmissionRepository,
	er domain.EventRepository, ep event.Publisher,
) *submissionUsecase {
	return &submissionUsecase{
		taskRepository:       tr,
		submissionRepository: sr,
		eventRepository:      er,
		eventPublisher:       ep,
	}
}

func (u *submissionUsecase) GetList(ctx context.Context, in dto.SubmissionListInput) (out *dto.SubmissionListOutput, err error) {
	// TODO: Change this to actual value.
	const submissionLimit = 20

	submissions, err := u.submissionRepository.FetchPaginated(ctx, in.ID, in.Offset, submissionLimit)
	if err != nil {
		return nil, errors.Wrap(err, "fetching submissions")
	}

	return toSubmissionListOutput(submissions), nil
}

func (u *submissionUsecase) Submit(ctx context.Context, in dto.SubmissionInput) (out *dto.IDOutput, err error) {
	// TODO: Make the logic transactional
	info := auth.MustExtract(ctx)

	taskID := in.ID

	if err := u.assureTaskExists(ctx, taskID); err != nil {
		return nil, err
	}

	exists, err := u.submissionRepository.UndoneExists(ctx, taskID, info.UserID)
	if err != nil {
		return nil, errors.Wrap(err, "checking if unfinished submission exists")
	}

	if exists {
		return nil, status.NewErr(http.StatusConflict, "unfinished submission exists")
	}

	submission := domain.Submission{
		ID:         uuid.New(),
		Timestamp:  time.Now(),
		UserID:     info.UserID,
		TaskID:     taskID,
		Repository: in.Repository,
		CommitHash: in.CommitHash,
	}

	if err := u.submissionRepository.Create(ctx, submission); err != nil {
		return nil, errors.Wrap(err, "creating submission")
	}

	// We are creating event manually because we don't want users to be notified on this.
	// A better approach might be determining what to do on event handler.
	event := domain.Event{
		ID:           uuid.New(),
		Kind:         domain.KindSubmit,
		Timestamp:    submission.Timestamp,
		SubmissionID: submission.ID,
	}

	if err := u.eventRepository.Create(ctx, event); err != nil {
		return nil, errors.Wrap(err, "creating event")
	}

	return toIDOutput(submission), nil
}

func (u *submissionUsecase) DecideApproval(ctx context.Context, in dto.SubmissionDecisionInput) (err error) {
	// TODO: Make the logic transactional

	if err := u.assureTaskExists(ctx, in.TaskID); err != nil {
		return err
	}

	submission, err := u.submissionRepository.FetchByID(ctx, in.SubmissionID)
	if err != nil {
		if errors.Is(err, domain.ErrSubmissionNotFound) {
			return status.NewErr(http.StatusNotFound, err.Error())
		}

		return errors.Wrap(err, "fetching submission")
	}

	if submission.IsDone {
		return status.NewErr(http.StatusForbidden, "submission is already done")
	}

	submission.IsDone = true
	if err := u.submissionRepository.Update(ctx, submission); err != nil {
		return errors.Wrap(err, "updatnig submission")
	}

	action := domain.SubmissionAction(in.Action)

	var eventKind domain.EventKind
	switch action {
	case domain.ActionApprove:
		eventKind = domain.KindApprove
	case domain.ActionReject:
		eventKind = domain.KindReject
	default:
		return errors.New("invalid action given")
	}

	if err := u.publishSubmissionEvent(ctx, eventKind, in.Extra, submission); err != nil {
		return err
	}

	return nil
}

func (u *submissionUsecase) Cancel(ctx context.Context, in dto.SubmissionIDInput) (err error) {
	// TODO: Make the logic transactional
	info := auth.MustExtract(ctx)

	if err := u.assureTaskExists(ctx, in.TaskID); err != nil {
		return err
	}

	submission, err := u.submissionRepository.FetchByID(ctx, in.SubmissionID)
	if err != nil {
		if errors.Is(err, domain.ErrSubmissionNotFound) {
			return status.NewErr(http.StatusNotFound, err.Error())
		}

		return errors.Wrap(err, "fetching submission")
	}

	if submission.IsDone {
		return status.NewErr(http.StatusForbidden, "submission is already done")
	}

	hasPermission := submission.UserID == info.UserID || info.Role == domain.RoleAdmin
	if !hasPermission {
		return status.NewErr(http.StatusForbidden, "no permission to the submission")
	}

	submission.IsDone = true
	if err := u.submissionRepository.Update(ctx, submission); err != nil {
		return errors.Wrap(err, "updatnig submission")
	}

	if err := u.publishSubmissionEvent(ctx, domain.KindCancel, "", submission); err != nil {
		return err
	}

	return nil
}

// assureTaskExists checks if task exists. if not exists, it will return an error.
func (u *submissionUsecase) assureTaskExists(ctx context.Context, taskID uuid.UUID) error {
	exists, err := u.taskRepository.ExistsByID(ctx, taskID)
	if err != nil {
		return errors.Wrap(err, "checking task exists")
	}

	if !exists {
		return status.NewErr(http.StatusNotFound, "task not found")
	}

	return nil
}

func (u *submissionUsecase) publishSubmissionEvent(
	ctx context.Context, kind domain.EventKind, extra string, submission domain.Submission,
) error {
	e := domain.Event{
		ID:         uuid.New(),
		Timestamp:  time.Now(),
		Kind:       kind,
		Extra:      extra,
		Submission: &submission,
	}

	if err := u.eventPublisher.Publish(ctx, event.TopicSubmission, e); err != nil {
		return errors.Wrap(err, "publishing event")
	}

	return nil
}
