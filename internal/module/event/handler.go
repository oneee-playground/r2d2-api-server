package event_module

import (
	"bytes"
	"context"
	"encoding/json"
	"text/template"

	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/global/email"
	"github.com/oneee-playground/r2d2-api-server/internal/global/event"
	"github.com/pkg/errors"
)

type EventHandler struct {
	userRepository  domain.UserRepository
	eventRepository domain.EventRepository
	emailSender     email.Sender
}

func NewEventHandler(
	es email.Sender, ur domain.UserRepository,
	er domain.EventRepository) *EventHandler {
	return &EventHandler{
		userRepository:  ur,
		eventRepository: er,
		emailSender:     es,
	}
}

func (h *EventHandler) Register(ctx context.Context, subscriber event.Subscriber) error {
	if err := subscriber.Subscribe(ctx, event.TopicSubmission,
		h.StoreEvent,
		// h.SendNotificationEmail,
	); err != nil {
		return err
	}

	return nil
}

func (h *EventHandler) StoreEvent(ctx context.Context, topic event.Topic, payload []byte) error {
	var event event.SubmissionEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return errors.Wrap(err, "unmarshalling payload")
	}

	domainEvent := domain.Event{
		ID:           event.ID,
		Kind:         event.Kind,
		Extra:        event.Extra,
		Timestamp:    event.Timestamp,
		SubmissionID: event.SubmissionID,
	}

	if err := h.eventRepository.Create(ctx, domainEvent); err != nil {
		return errors.Wrap(err, "creating event")
	}

	return nil
}

var emailTemplate = template.Must(template.New("event-mail-template").Parse(`
안녕하세요, {{.Username}}님.

제출하신 답안의 상태가 바뀌었습니다: {{.EventKind}}
`))

type _emailData struct {
	Username  string
	EventKind string
}

func (h *EventHandler) SendNotificationEmail(ctx context.Context, topic event.Topic, payload []byte) error {
	var event event.SubmissionEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return errors.Wrap(err, "unmarshalling payload")
	}

	user, err := h.userRepository.FetchByID(ctx, event.UserID)
	if err != nil {
		return errors.Wrap(err, "fetching user")
	}

	buf := bytes.NewBuffer(nil)

	data := _emailData{
		Username:  user.Username,
		EventKind: string(event.Kind),
	}

	if err := emailTemplate.Execute(buf, data); err != nil {
		return errors.Wrap(err, "executing template")
	}

	if err := h.emailSender.Send(ctx, user.Email, "r2d2 제출 알림", buf); err != nil {
		return errors.Wrap(err, "sending email")
	}

	return err
}
