package exec_module

import (
	"context"
	"encoding/json"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/global/event"
	"github.com/pkg/errors"
)

type EventHandler struct {
	submissionRepository domain.SubmissionRepository
	sectionRepository    domain.SectionRepository
	resourceRepository   domain.ResourceRepository

	eventPublisher event.Publisher
	jobQueue       JobQueue
	imageBuilder   ImageBuilder
	contextStorage ExecContextStroage
}

func NewEventHandler(
	sur domain.SubmissionRepository, ser domain.SectionRepository, rr domain.ResourceRepository,
	ep event.Publisher, jq JobQueue, ib ImageBuilder, cs ExecContextStroage,
) *EventHandler {
	return &EventHandler{
		submissionRepository: sur,
		sectionRepository:    ser,
		resourceRepository:   rr,
		eventPublisher:       ep,
		jobQueue:             jq,
		imageBuilder:         ib,
		contextStorage:       cs,
	}
}

func (h *EventHandler) Register(ctx context.Context, subscriber event.Subscriber) error {
	if err := subscriber.Subscribe(ctx, event.TopicSubmission,
		h.StartBuild,
	); err != nil {
		return err
	}

	if err := subscriber.Subscribe(ctx, event.TopicBuild,
		h.EnqueueJob, h.NotifyBuildFailure,
	); err != nil {
		return err
	}

	if err := subscriber.Subscribe(ctx, event.TopicTest,
		h.NotifyTestResult,
	); err != nil {
		return err
	}

	return nil
}

func (h *EventHandler) StartBuild(ctx context.Context, topic event.Topic, payload []byte) error {
	var ev event.SubmissionEvent
	if err := json.Unmarshal(payload, &ev); err != nil {
		return errors.Wrap(err, "unmarshalling payload")
	}

	if ev.Kind != domain.KindApprove {
		return event.NoErrSkipHandler
	}

	err := h.publishSubmissionEvent(ctx, domain.KindBuildStart, "", ev.SubmissionID, ev.UserID)
	if err != nil {
		return err
	}

	submission, err := h.submissionRepository.FetchByID(ctx, ev.SubmissionID)
	if err != nil {
		return errors.Wrap(err, "fetching submission")
	}

	buildOpts := BuildOpts{
		ID:         submission.ID,
		TaskID:     submission.TaskID,
		Repository: submission.Repository,
		CommitHash: submission.CommitHash,
		Platform:   runtime.GOOS + "/" + runtime.GOARCH,
	}

	if err := h.imageBuilder.RequestBuild(ctx, buildOpts); err != nil {
		return errors.Wrap(err, "requesting to build image")
	}

	execCtx := ExecContext{
		TaskID:     submission.TaskID,
		Repository: submission.Repository,
		CommitHash: submission.CommitHash,
		UserID:     submission.UserID,
	}

	if err := h.contextStorage.Set(ctx, submission.ID, execCtx); err != nil {
		return errors.Wrap(err, "setting exec context")
	}

	return nil
}

func (h *EventHandler) EnqueueJob(ctx context.Context, topic event.Topic, payload []byte) error {
	var ev event.ExecEvent
	if err := json.Unmarshal(payload, &ev); err != nil {
		return errors.Wrap(err, "unmarshalling payload")
	}

	if !ev.Success {
		return event.NoErrSkipHandler
	}

	submissionID := ev.ID

	execCtx, err := h.contextStorage.Get(ctx, ev.ID)
	if err != nil {
		return errors.Wrap(err, "fetching exec context")
	}

	err = h.publishSubmissionEvent(ctx, domain.KindBuildSuccess, "", submissionID, execCtx.UserID)
	if err != nil {
		return err
	}

	job := Job{
		TaskID: execCtx.TaskID,
		Submission: Submission{
			ID:         submissionID,
			Repository: execCtx.Repository,
			CommitHash: execCtx.CommitHash,
		},
	}

	fetchSectionOpts := domain.FetchSectionsOption{IncludeContent: false}

	sections, err := h.sectionRepository.FetchAllByTaskID(ctx, execCtx.TaskID, fetchSectionOpts)
	if err != nil {
		return errors.Wrap(err, "fetching sections")
	}

	job.Sections = make([]Section, len(sections))
	for idx, section := range sections {
		job.Sections[idx] = Section{
			ID:   section.ID,
			Type: section.Type,
		}
	}

	resources, err := h.resourceRepository.FetchAllByTaskID(ctx, execCtx.TaskID)
	if err != nil {
		return errors.Wrap(err, "fetching resources")
	}

	job.Resources = make([]Resource, len(resources))
	for idx, resource := range resources {
		job.Resources[idx] = Resource{
			Image:     resource.Image,
			Name:      resource.Name,
			Port:      resource.Port,
			CPU:       resource.CPU,
			Memory:    resource.Memory,
			IsPrimary: resource.IsPrimary,
		}
	}

	if err := h.jobQueue.Append(ctx, &job); err != nil {
		return errors.Wrap(err, "appending job to the queue")
	}

	return h.publishSubmissionEvent(ctx, domain.KindTestStart, "", submissionID, execCtx.UserID)
}

func (h *EventHandler) NotifyBuildFailure(ctx context.Context, topic event.Topic, payload []byte) error {
	var ev event.ExecEvent
	if err := json.Unmarshal(payload, &ev); err != nil {
		return errors.Wrap(err, "unmarshalling payload")
	}

	if ev.Success {
		return event.NoErrSkipHandler
	}

	execCtx, err := h.contextStorage.Get(ctx, ev.ID)
	if err != nil {
		return errors.Wrap(err, "fetching exec context")
	}

	err = h.publishSubmissionEvent(ctx, domain.KindBuildFail, ev.Extra, ev.ID, execCtx.UserID)
	if err != nil {
		return err
	}

	return h.deleteExecContext(ctx, ev.ID)
}

func (h *EventHandler) NotifyTestResult(ctx context.Context, topic event.Topic, payload []byte) error {
	var ev event.ExecEvent
	if err := json.Unmarshal(payload, &ev); err != nil {
		return errors.Wrap(err, "unmarshalling payload")
	}

	execCtx, err := h.contextStorage.Get(ctx, ev.ID)
	if err != nil {
		return errors.Wrap(err, "fetching exec context")
	}

	var eventKind domain.EventKind
	if ev.Success {
		eventKind = domain.KindTestSuccess
	} else {
		eventKind = domain.KindTestFail
	}

	err = h.publishSubmissionEvent(ctx, eventKind, ev.Extra, ev.ID, execCtx.UserID)
	if err != nil {
		return err
	}

	return h.deleteExecContext(ctx, ev.ID)
}

func (h *EventHandler) deleteExecContext(ctx context.Context, id uuid.UUID) error {
	if err := h.contextStorage.Delete(ctx, id); err != nil {
		return errors.Wrap(err, "deleting exec context")
	}

	return nil
}

func (h *EventHandler) publishSubmissionEvent(
	ctx context.Context, kind domain.EventKind, extra string, submissionID, userID uuid.UUID,
) error {
	e := event.SubmissionEvent{
		ID:           uuid.New(),
		Timestamp:    time.Now(),
		Kind:         kind,
		Extra:        extra,
		SubmissionID: submissionID,
		UserID:       userID,
	}

	if err := h.eventPublisher.Publish(ctx, event.TopicSubmission, e); err != nil {
		return errors.Wrap(err, "publishing event")
	}

	return nil
}
