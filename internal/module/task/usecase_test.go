package task_module_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/domain/dto"
	"github.com/oneee-playground/r2d2-api-server/internal/global/status"
	task_module "github.com/oneee-playground/r2d2-api-server/internal/module/task"
	"github.com/oneee-playground/r2d2-api-server/test/mocks"
	"github.com/stretchr/testify/suite"
)

func TestTaskUsecaseSuite(t *testing.T) {
	suite.Run(t, new(TaskUsecaseSuite))
}

type TaskUsecaseSuite struct {
	suite.Suite

	usecase domain.TaskUsecase

	ctl  *gomock.Controller
	mock struct {
		taskRepository *mocks.MockTaskRepository
	}
}

func (s *TaskUsecaseSuite) SetupTest() {
	s.ctl = gomock.NewController(s.T())
	s.mock.taskRepository = mocks.NewMockTaskRepository(s.ctl)

	s.usecase = task_module.NewTaskUsecase(s.mock.taskRepository)
}

func (s *TaskUsecaseSuite) TestChangeStage() {
	draftTask := domain.Task{
		ID:          uuid.New(),
		Title:       "title",
		Description: "description",
		Stage:       domain.StageDraft,
	}

	availableTask := draftTask
	availableTask.Stage = domain.StageAvailable

	testcases := []struct {
		desc        string
		targetStage domain.TaskStage
		setup       func()
		checkErr    func(err error) bool
	}{
		// TODO: Cover case when input stage is domain.StageFixing
		{
			desc:        "draft -> available",
			targetStage: domain.StageAvailable,
			setup: func() {
				s.mock.taskRepository.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).Return(draftTask, nil)
				s.mock.taskRepository.EXPECT().
					Update(gomock.Any(), gomock.Any()).Return(nil)
			},
			checkErr: func(err error) bool { return err == nil },
		},
		{
			desc:        "available -> fixing",
			targetStage: domain.StageFixing,
			setup: func() {
				s.mock.taskRepository.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).Return(availableTask, nil)
				s.mock.taskRepository.EXPECT().
					Update(gomock.Any(), gomock.Any()).Return(nil)
			},
			checkErr: func(err error) bool { return err == nil },
		},
		{
			desc:        "available -> available",
			targetStage: domain.StageAvailable,
			setup: func() {
				s.mock.taskRepository.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).Return(availableTask, nil)
			},
			checkErr: func(err error) bool {
				sErr, ok := err.(status.Error)
				return ok && sErr.StatusCode == http.StatusForbidden
			},
		},
		{
			desc:        "available -> draft",
			targetStage: domain.StageDraft,
			setup: func() {
				s.mock.taskRepository.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).Return(availableTask, nil)
			},
			checkErr: func(err error) bool {
				sErr, ok := err.(status.Error)
				return ok && sErr.StatusCode == http.StatusForbidden
			},
		},
	}

	for _, tc := range testcases {
		s.Run(tc.desc, func() {
			ctx := context.Background()

			input := dto.TaskStageInput{
				Stage: string(tc.targetStage),
			}

			tc.setup()

			err := s.usecase.ChangeStage(ctx, input)
			s.True(tc.checkErr(err), err)
		})
	}
}
