package event

import (
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/domain/dto"
)

func toEventListOutput(events []domain.Event) *dto.EventListOutput {
	out := make(dto.EventListOutput, len(events))
	for i, event := range events {
		out[i] = dto.EventListElem{
			Kind:      string(event.Kind),
			Extra:     event.Extra,
			Timestamp: event.Timestamp,
		}
	}

	return &out
}
