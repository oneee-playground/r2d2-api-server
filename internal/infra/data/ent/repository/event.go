package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/data/ent/model"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/data/ent/model/event"
)

type EventRepository struct {
	client *model.EventClient
}

var _ domain.EventRepository = (*EventRepository)(nil)

func NewEventRepository(client *model.EventClient) *EventRepository {
	return &EventRepository{client: client}
}

func (r *EventRepository) Create(ctx context.Context, event domain.Event) error {
	return r.client.Create().
		SetID(event.ID).
		SetKind(string(event.Kind)).
		SetSubmissionID(event.SubmissionID).
		SetTimestamp(event.Timestamp).
		SetExtra(event.Extra).
		Exec(ctx)
}

func (r *EventRepository) FetchAllBySubmissionID(ctx context.Context, id uuid.UUID) ([]domain.Event, error) {
	models, err := r.client.Query().
		Where(event.SubmissionID(id)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	events := make([]domain.Event, len(models))
	for idx, model := range models {
		events[idx] = domain.Event{
			ID:           model.ID,
			Kind:         domain.EventKind(model.Kind),
			Extra:        model.Extra,
			Timestamp:    model.Timestamp,
			SubmissionID: model.SubmissionID,
		}
	}

	return events, nil
}
