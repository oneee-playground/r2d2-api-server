package submission_module_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/domain/dto"
	"github.com/oneee-playground/r2d2-api-server/internal/global/auth"
	"github.com/oneee-playground/r2d2-api-server/internal/global/status"
	submission_module "github.com/oneee-playground/r2d2-api-server/internal/module/submission"
	"github.com/oneee-playground/r2d2-api-server/test/mocks"
	"github.com/stretchr/testify/suite"
)

func TestSubmissionUsecaseSuite(t *testing.T) {
	suite.Run(t, new(SubmissionUsecaseSuite))
}

type SubmissionUsecaseSuite struct {
	suite.Suite

	usecase domain.SubmissionUsecase

	ctl  *gomock.Controller
	mock struct {
		taskRepository       *mocks.MockTaskRepository
		submissionRepository *mocks.MockSubmissionRepository
		eventRepository      *mocks.MockEventRepository
		eventPublisher       *mocks.MockPublisher
	}
}

func (s *SubmissionUsecaseSuite) SetupTest() {
	s.ctl = gomock.NewController(s.T())
	s.mock.taskRepository = mocks.NewMockTaskRepository(s.ctl)
	s.mock.submissionRepository = mocks.NewMockSubmissionRepository(s.ctl)
	s.mock.eventRepository = mocks.NewMockEventRepository(s.ctl)
	s.mock.eventPublisher = mocks.NewMockPublisher(s.ctl)

	s.usecase = submission_module.NewSubmissionUsecase(
		s.mock.taskRepository, s.mock.submissionRepository,
		s.mock.eventRepository, s.mock.eventPublisher,
	)
}

func (s *SubmissionUsecaseSuite) TestSubmit() {
	testcases := []struct {
		desc     string
		setup    func()
		checkErr func(err error) bool
	}{
		{
			desc: "success",
			setup: func() {
				s.mock.taskRepository.EXPECT().
					ExistsByID(gomock.Any(), gomock.Any()).Return(true, nil)
				s.mock.submissionRepository.EXPECT().
					UndoneExists(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)
				s.mock.submissionRepository.EXPECT().
					Create(gomock.Any(), gomock.Any()).Return(nil)
				s.mock.eventRepository.EXPECT().
					Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			checkErr: func(err error) bool { return err == nil },
		},
		{
			desc: "task does not exist",
			setup: func() {
				s.mock.taskRepository.EXPECT().
					ExistsByID(gomock.Any(), gomock.Any()).Return(false, nil)
			},
			checkErr: func(err error) bool {
				sErr, ok := err.(status.Error)
				return ok && sErr.StatusCode == http.StatusNotFound
			},
		},
		{
			desc: "duplicate submission",
			setup: func() {
				s.mock.taskRepository.EXPECT().
					ExistsByID(gomock.Any(), gomock.Any()).Return(true, nil)
				s.mock.submissionRepository.EXPECT().
					UndoneExists(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
			},
			checkErr: func(err error) bool {
				sErr, ok := err.(status.Error)
				return ok && sErr.StatusCode == http.StatusConflict
			},
		},
	}

	ctx := auth.Inject(context.Background(), auth.Payload{})
	for _, tc := range testcases {
		s.Run(tc.desc, func() {
			tc.setup()

			_, err := s.usecase.Submit(ctx, dto.SubmissionInput{})
			s.True(tc.checkErr(err), err)
		})
	}
}

func (s *SubmissionUsecaseSuite) TestDecideApproval() {
	undoneSubmission := domain.Submission{IsDone: false}
	doneSubmission := domain.Submission{IsDone: true}

	testcases := []struct {
		desc     string
		action   domain.SubmissionAction
		setup    func()
		checkErr func(err error) bool
	}{
		{
			desc:   "success (approve)",
			action: domain.ActionApprove,
			setup: func() {
				s.mock.taskRepository.EXPECT().
					ExistsByID(gomock.Any(), gomock.Any()).Return(true, nil)
				s.mock.submissionRepository.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).Return(undoneSubmission, nil)
				s.mock.submissionRepository.EXPECT().
					Update(gomock.Any(), gomock.Any()).Return(nil)
				s.mock.eventPublisher.EXPECT().
					Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			checkErr: func(err error) bool { return err == nil },
		},
		{
			desc:   "success (reject)",
			action: domain.ActionReject,
			setup: func() {
				s.mock.taskRepository.EXPECT().
					ExistsByID(gomock.Any(), gomock.Any()).Return(true, nil)
				s.mock.submissionRepository.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).Return(undoneSubmission, nil)
				s.mock.submissionRepository.EXPECT().
					Update(gomock.Any(), gomock.Any()).Return(nil)
				s.mock.eventPublisher.EXPECT().
					Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			checkErr: func(err error) bool { return err == nil },
		},
		{
			desc: "task does not exist",
			setup: func() {
				s.mock.taskRepository.EXPECT().
					ExistsByID(gomock.Any(), gomock.Any()).Return(false, nil)
			},
			checkErr: func(err error) bool {
				sErr, ok := err.(status.Error)
				return ok && sErr.StatusCode == http.StatusNotFound
			},
		},
		{
			desc: "submission not found",
			setup: func() {
				s.mock.taskRepository.EXPECT().
					ExistsByID(gomock.Any(), gomock.Any()).Return(true, nil)
				s.mock.submissionRepository.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).Return(undoneSubmission, domain.ErrSubmissionNotFound)
			},
			checkErr: func(err error) bool {
				sErr, ok := err.(status.Error)
				return ok && sErr.StatusCode == http.StatusNotFound
			},
		},
		{
			desc: "submission done",
			setup: func() {
				s.mock.taskRepository.EXPECT().
					ExistsByID(gomock.Any(), gomock.Any()).Return(true, nil)
				s.mock.submissionRepository.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).Return(doneSubmission, nil)
			},
			checkErr: func(err error) bool {
				sErr, ok := err.(status.Error)
				return ok && sErr.StatusCode == http.StatusForbidden
			},
		},
	}

	ctx := auth.Inject(context.Background(), auth.Payload{})
	for _, tc := range testcases {
		s.Run(tc.desc, func() {
			tc.setup()

			in := dto.SubmissionDecisionInput{
				Action: string(tc.action),
			}

			err := s.usecase.DecideApproval(ctx, in)
			s.True(tc.checkErr(err), err)
		})
	}
}

func (s *SubmissionUsecaseSuite) TestCancel() {
	testUser := auth.Payload{
		UserID: uuid.New(),
		Role:   domain.RoleMember,
	}

	adminUser := auth.Payload{
		UserID: uuid.New(),
		Role:   domain.RoleAdmin,
	}

	otherUser := auth.Payload{
		UserID: uuid.New(),
	}

	s.Require().NotEqual(testUser.UserID, adminUser.UserID)
	s.Require().NotEqual(testUser.UserID, otherUser.UserID)

	doneSubmission := domain.Submission{IsDone: true}
	userSubmission := domain.Submission{
		UserID: testUser.UserID,
		IsDone: false,
	}

	testcases := []struct {
		desc     string
		userInfo auth.Payload
		setup    func()
		checkErr func(err error) bool
	}{
		{
			desc:     "success (user self)",
			userInfo: testUser,
			setup: func() {
				s.mock.taskRepository.EXPECT().
					ExistsByID(gomock.Any(), gomock.Any()).Return(true, nil)
				s.mock.submissionRepository.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).Return(userSubmission, nil)
				s.mock.submissionRepository.EXPECT().
					Update(gomock.Any(), gomock.Any()).Return(nil)
				s.mock.eventPublisher.EXPECT().
					Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			checkErr: func(err error) bool { return err == nil },
		},
		{
			desc:     "success (admin)",
			userInfo: adminUser,
			setup: func() {
				s.mock.taskRepository.EXPECT().
					ExistsByID(gomock.Any(), gomock.Any()).Return(true, nil)
				s.mock.submissionRepository.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).Return(userSubmission, nil)
				s.mock.submissionRepository.EXPECT().
					Update(gomock.Any(), gomock.Any()).Return(nil)
				s.mock.eventPublisher.EXPECT().
					Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			checkErr: func(err error) bool { return err == nil },
		},
		{
			desc: "task does not exist",
			setup: func() {
				s.mock.taskRepository.EXPECT().
					ExistsByID(gomock.Any(), gomock.Any()).Return(false, nil)
			},
			checkErr: func(err error) bool {
				sErr, ok := err.(status.Error)
				return ok && sErr.StatusCode == http.StatusNotFound
			},
		},
		{
			desc: "submission not found",
			setup: func() {
				s.mock.taskRepository.EXPECT().
					ExistsByID(gomock.Any(), gomock.Any()).Return(true, nil)
				s.mock.submissionRepository.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).Return(userSubmission, domain.ErrSubmissionNotFound)
			},
			checkErr: func(err error) bool {
				sErr, ok := err.(status.Error)
				return ok && sErr.StatusCode == http.StatusNotFound
			},
		},
		{
			desc: "submission done",
			setup: func() {
				s.mock.taskRepository.EXPECT().
					ExistsByID(gomock.Any(), gomock.Any()).Return(true, nil)
				s.mock.submissionRepository.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).Return(doneSubmission, nil)
			},
			checkErr: func(err error) bool {
				sErr, ok := err.(status.Error)
				return ok && sErr.StatusCode == http.StatusForbidden
			},
		},
		{
			desc:     "no permission",
			userInfo: otherUser,
			setup: func() {
				s.mock.taskRepository.EXPECT().
					ExistsByID(gomock.Any(), gomock.Any()).Return(true, nil)
				s.mock.submissionRepository.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).Return(userSubmission, nil)
			},
			checkErr: func(err error) bool {
				sErr, ok := err.(status.Error)
				return ok && sErr.StatusCode == http.StatusForbidden
			},
		},
	}

	ctx := context.Background()

	for _, tc := range testcases {
		s.Run(tc.desc, func() {
			tc.setup()

			ctx = auth.Inject(ctx, tc.userInfo)

			err := s.usecase.Cancel(ctx, dto.SubmissionIDInput{})
			s.True(tc.checkErr(err), err)
		})
	}
}
