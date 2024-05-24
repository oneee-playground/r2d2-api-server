package resource_module_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/domain/dto"
	"github.com/oneee-playground/r2d2-api-server/internal/global/status"
	resource_module "github.com/oneee-playground/r2d2-api-server/internal/module/resource"
	"github.com/oneee-playground/r2d2-api-server/test/mocks"
	"github.com/stretchr/testify/suite"
)

func TestResourceUsecaseSuite(t *testing.T) {
	suite.Run(t, new(ResourceUsecaseSuite))
}

type ResourceUsecaseSuite struct {
	suite.Suite

	usecase domain.ResourceUsecase

	ctl  *gomock.Controller
	mock struct {
		taskRepository     *mocks.MockTaskRepository
		resourceRepository *mocks.MockResourceRepository
	}
}

func (s *ResourceUsecaseSuite) SetupTest() {
	s.ctl = gomock.NewController(s.T())
	s.mock.taskRepository = mocks.NewMockTaskRepository(s.ctl)
	s.mock.resourceRepository = mocks.NewMockResourceRepository(s.ctl)

	s.usecase = resource_module.NewResourceUsecase(s.mock.resourceRepository, s.mock.taskRepository)
}

func (s *ResourceUsecaseSuite) TestCreateResource() {
	availableTask := domain.Task{Stage: domain.StageAvailable}
	draftTask := domain.Task{Stage: domain.StageDraft}

	testcases := []struct {
		desc     string
		setup    func()
		checkErr func(err error) bool
	}{
		{
			desc: "success",
			setup: func() {
				s.mock.taskRepository.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).Return(draftTask, nil)
				s.mock.resourceRepository.EXPECT().
					Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			checkErr: func(err error) bool { return err == nil },
		},
		{
			desc: "task not found",
			setup: func() {
				s.mock.taskRepository.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).Return(draftTask, domain.ErrTaskNotFound)
			},
			checkErr: func(err error) bool {
				sErr, ok := err.(status.Error)
				return ok && sErr.StatusCode == http.StatusNotFound
			},
		},
		{
			desc: "available task",
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
			desc: "resource already exists",
			setup: func() {
				s.mock.taskRepository.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).Return(draftTask, nil)
				s.mock.resourceRepository.EXPECT().
					Create(gomock.Any(), gomock.Any()).Return(domain.ErrDuplicateResource)
			},
			checkErr: func(err error) bool {
				sErr, ok := err.(status.Error)
				return ok && sErr.StatusCode == http.StatusConflict
			},
		},
	}

	for _, tc := range testcases {
		s.Run(tc.desc, func() {
			ctx := context.Background()

			tc.setup()

			err := s.usecase.CreateResource(ctx, dto.CreateResourceInput{})
			s.True(tc.checkErr(err), err)
		})
	}
}

func (s *ResourceUsecaseSuite) TestDeleteResource() {
	availableTask := domain.Task{Stage: domain.StageAvailable}
	draftTask := domain.Task{Stage: domain.StageDraft}

	testcases := []struct {
		desc     string
		setup    func()
		checkErr func(err error) bool
	}{
		{
			desc: "success",
			setup: func() {
				s.mock.taskRepository.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).Return(draftTask, nil)
				s.mock.resourceRepository.EXPECT().
					Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			checkErr: func(err error) bool { return err == nil },
		},
		{
			desc: "task not found",
			setup: func() {
				s.mock.taskRepository.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).Return(draftTask, domain.ErrTaskNotFound)
			},
			checkErr: func(err error) bool {
				sErr, ok := err.(status.Error)
				return ok && sErr.StatusCode == http.StatusNotFound
			},
		},
		{
			desc: "available task",
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
			desc: "resource not found",
			setup: func() {
				s.mock.taskRepository.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).Return(draftTask, nil)
				s.mock.resourceRepository.EXPECT().
					Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(domain.ErrResourceNotFound)
			},
			checkErr: func(err error) bool {
				sErr, ok := err.(status.Error)
				return ok && sErr.StatusCode == http.StatusNotFound
			},
		},
	}

	for _, tc := range testcases {
		s.Run(tc.desc, func() {
			ctx := context.Background()

			tc.setup()

			err := s.usecase.DeleteResource(ctx, dto.ResourceIDInput{})
			s.True(tc.checkErr(err), err)
		})
	}
}
