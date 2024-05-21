package event

import (
	"bytes"
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/global/event"
	"github.com/oneee-playground/r2d2-api-server/test/mocks"
	"github.com/stretchr/testify/suite"
)

func TestEventHandlerSuite(t *testing.T) {
	suite.Run(t, new(EventHandlerSuite))
}

type EventHandlerSuite struct {
	suite.Suite

	handler *EventHandler

	ctl  *gomock.Controller
	mock struct {
		userRepository  *mocks.MockUserRepository
		eventRepository *mocks.MockEventRepository
		emailSender     *mocks.MockSender
	}
}

func (s *EventHandlerSuite) SetupTest() {
	s.ctl = gomock.NewController(s.T())
	s.mock.userRepository = mocks.NewMockUserRepository(s.ctl)
	s.mock.eventRepository = mocks.NewMockEventRepository(s.ctl)
	s.mock.emailSender = mocks.NewMockSender(s.ctl)

	s.handler = NewEventHandler(
		s.mock.emailSender,
		s.mock.userRepository,
		s.mock.eventRepository,
	)
}

func (s *EventHandlerSuite) TestSendNotificationEmail() {
	testEvent := domain.Event{
		ID:         uuid.New(),
		Kind:       domain.KindApprove,
		Submission: &domain.Submission{},
	}

	testUser := domain.User{
		Username: "test",
		Email:    "test@example.com",
	}

	testcases := []struct {
		desc    string
		setup   func()
		wantErr bool
	}{
		{
			desc: "success",
			setup: func() {
				s.mock.userRepository.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).Return(testUser, nil)
				s.mock.emailSender.EXPECT().
					Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Do(func(_ context.Context, address string, subject string, content *bytes.Buffer) {
						s.Require().Equal(testUser.Email, address)
						s.Require().Contains(content.String(), testUser.Username)
						s.Require().Contains(content.String(), string(testEvent.Kind))
					}).Return(nil)
			},
			wantErr: false,
		},
		{
			desc: "user not found",
			setup: func() {
				s.mock.userRepository.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).Return(testUser, domain.ErrUserNotFound)
			},
			wantErr: true,
		},
	}

	for _, tc := range testcases {
		s.Run(tc.desc, func() {
			ctx := context.Background()

			tc.setup()

			err := s.handler.SendNotificationEmail(ctx, event.TopicSubmission, testEvent)
			if tc.wantErr {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}
