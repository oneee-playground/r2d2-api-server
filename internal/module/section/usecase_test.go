package section_module_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/domain/dto"
	"github.com/oneee-playground/r2d2-api-server/internal/global/status"
	section_module "github.com/oneee-playground/r2d2-api-server/internal/module/section"
	"github.com/oneee-playground/r2d2-api-server/test/mocks"
	"github.com/stretchr/testify/suite"
)

func TestSectionUsecaseSuite(t *testing.T) {
	suite.Run(t, new(SectionUsecaseSuite))
}

type SectionUsecaseSuite struct {
	suite.Suite

	usecase domain.SectionUsecase

	ctl  *gomock.Controller
	mock struct {
		taskRepository    *mocks.MockTaskRepository
		sectionRepository *mocks.MockSectionRepository
	}
}

func (s *SectionUsecaseSuite) SetupTest() {
	s.ctl = gomock.NewController(s.T())

	s.mock.taskRepository = mocks.NewMockTaskRepository(s.ctl)
	s.mock.sectionRepository = mocks.NewMockSectionRepository(s.ctl)

	s.usecase = section_module.NewSectionUsecase(s.mock.sectionRepository, s.mock.taskRepository)
}

func (s *SectionUsecaseSuite) TestUpdateSection() {
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
				s.mock.sectionRepository.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).Return(domain.Section{}, nil)
				s.mock.sectionRepository.EXPECT().
					Update(gomock.Any(), gomock.Any()).Return(nil)
			},
			checkErr: func(err error) bool { return err == nil },
		},
		{
			desc: "task not found",
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
			desc: "section not found",
			setup: func() {
				s.mock.taskRepository.EXPECT().
					ExistsByID(gomock.Any(), gomock.Any()).Return(true, nil)
				s.mock.sectionRepository.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).Return(domain.Section{}, domain.ErrSectionNotFound)
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

			err := s.usecase.UpdateSection(ctx, dto.UpdateSectionInput{})
			s.True(tc.checkErr(err), err)
		})
	}
}

func (s *SectionUsecaseSuite) TestChangeIndex() {
	testSections := []domain.Section{
		{ID: uuid.New(), Index: 0},
		{ID: uuid.New(), Index: 1},
		{ID: uuid.New(), Index: 2},
		{ID: uuid.New(), Index: 3},
		{ID: uuid.New(), Index: 4},
	}

	invalidID := uuid.New()
	for _, section := range testSections {
		s.Require().NotEqual(section.ID, invalidID)
	}

	testcases := []struct {
		desc      string
		sectionID uuid.UUID
		index     int
		setup     func()
		checkErr  func(err error) bool
	}{
		{
			desc:      "success",
			sectionID: testSections[0].ID, index: 1,
			setup: func() {
				s.mock.taskRepository.EXPECT().
					ExistsByID(gomock.Any(), gomock.Any()).Return(true, nil)
				s.mock.sectionRepository.EXPECT().
					FetchAllByTaskID(gomock.Any(), gomock.Any(), gomock.Any()).Return(testSections, nil)
				s.mock.sectionRepository.EXPECT().
					SaveIndexes(gomock.Any(), gomock.Any()).Return(nil)
			},
			checkErr: func(err error) bool { return err == nil },
		},
		{
			desc: "task not found",
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
			desc:      "section not found",
			sectionID: invalidID, index: 0,
			setup: func() {
				s.mock.taskRepository.EXPECT().
					ExistsByID(gomock.Any(), gomock.Any()).Return(true, nil)
				s.mock.sectionRepository.EXPECT().
					FetchAllByTaskID(gomock.Any(), gomock.Any(), gomock.Any()).Return(testSections, nil)
			},
			checkErr: func(err error) bool {
				sErr, ok := err.(status.Error)
				return ok && sErr.StatusCode == http.StatusNotFound
			},
		},
		{
			desc:  "invalid index (-1)",
			index: -1,
			setup: func() {
				s.mock.taskRepository.EXPECT().
					ExistsByID(gomock.Any(), gomock.Any()).Return(true, nil)
				s.mock.sectionRepository.EXPECT().
					FetchAllByTaskID(gomock.Any(), gomock.Any(), gomock.Any()).Return(testSections, nil)
			},
			checkErr: func(err error) bool {
				sErr, ok := err.(status.Error)
				return ok && sErr.StatusCode == http.StatusForbidden
			},
		},
		{
			desc:  "invalid index (len)",
			index: len(testSections),
			setup: func() {
				s.mock.taskRepository.EXPECT().
					ExistsByID(gomock.Any(), gomock.Any()).Return(true, nil)
				s.mock.sectionRepository.EXPECT().
					FetchAllByTaskID(gomock.Any(), gomock.Any(), gomock.Any()).Return(testSections, nil)
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

			in := dto.SectionIndexInput{
				SectionIDInput: dto.SectionIDInput{
					SectionID: tc.sectionID,
				},
				Index: tc.index,
			}

			tc.setup()

			err := s.usecase.ChangeIndex(ctx, in)
			s.True(tc.checkErr(err), err)
		})
	}
}
