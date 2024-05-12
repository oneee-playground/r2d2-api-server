package user_module_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/domain/dto"
	"github.com/oneee-playground/r2d2-api-server/internal/global/auth"
	user_module "github.com/oneee-playground/r2d2-api-server/internal/module/user"
	"github.com/oneee-playground/r2d2-api-server/test/mocks"
	"github.com/stretchr/testify/suite"
)

func TestUserUsecaseSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseSuite))
}

type UserUsecaseSuite struct {
	suite.Suite

	usecase domain.UserUsecase

	ctl  *gomock.Controller
	mock struct {
		userRepository *mocks.MockUserRepository
	}
}

func (s *UserUsecaseSuite) SetupTest() {
	s.ctl = gomock.NewController(s.T())
	s.mock.userRepository = mocks.NewMockUserRepository(s.ctl)

	s.usecase = user_module.NewUserUsecase(s.mock.userRepository)
}

func (s *UserUsecaseSuite) TestGetSelfInfo() {
	testUser := domain.User{
		ID:         uuid.New(),
		Username:   "user",
		Email:      "email@example.com",
		ProfileURL: "profile.com",
		Role:       domain.RoleAdmin,
	}

	testUserInfo := dto.UserInfo{
		ID:         testUser.ID,
		Username:   testUser.Username,
		ProfileURL: testUser.ProfileURL,
		Role:       testUser.Role.String(),
	}

	testAuthPayload := auth.Payload{
		UserID: testUser.ID,
		Role:   testUser.Role,
	}

	testcases := []struct {
		desc    string
		payload auth.Payload
		setup   func()
		wantErr bool
	}{
		{
			desc:    "success",
			payload: testAuthPayload,
			setup: func() {
				s.mock.userRepository.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).Return(testUser, nil)
			},
			wantErr: false,
		},
	}

	for _, tc := range testcases {
		s.Run(tc.desc, func() {
			ctx := context.Background()

			tc.setup()

			ctx = auth.Inject(ctx, tc.payload)

			out, err := s.usecase.GetSelfInfo(ctx)
			if tc.wantErr {
				s.Error(err)

				return
			}

			if !s.NotNil(out) {
				return
			}

			s.Equal(testUserInfo, *out)
		})
	}
}
