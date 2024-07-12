package event_module

import (
	"context"
	"sort"

	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/domain/dto"
	"github.com/pkg/errors"
)

type eventUsecase struct {
	eventRepository domain.EventRepository
}

var _ domain.EventUsecase = (*eventUsecase)(nil)

func NewEventUsecase(er domain.EventRepository) *eventUsecase {
	return &eventUsecase{eventRepository: er}
}

func (u *eventUsecase) GetAllFromSubmission(ctx context.Context, in dto.SubmissionIDInput) (out *dto.EventListOutput, err error) {
	events, err := u.eventRepository.FetchAllBySubmissionID(ctx, in.SubmissionID)
	if err != nil {
		return nil, errors.Wrap(err, "fetching events")
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Timestamp.Before(events[j].Timestamp)
	})

	return toEventListOutput(events), nil
}
