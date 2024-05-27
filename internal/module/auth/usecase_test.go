package auth_module_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/domain/dto"
	"github.com/oneee-playground/r2d2-api-server/internal/global/status"
	auth_module "github.com/oneee-playground/r2d2-api-server/internal/module/auth"
	"github.com/oneee-playground/r2d2-api-server/test/mocks"
	"github.com/stretchr/testify/suite"
)

func TestAuthUsecaseSuite(t *testing.T) {
	suite.Run(t, new(AuthUsecaseSuite))
}

type AuthUsecaseSuite struct {
	suite.Suite

	usecase domain.AuthUsecase

	ctl  *gomock.Controller
	mock struct {
		oauth          *mocks.MockOAuthClient
		tokenIssuer    *mocks.MockTokenIssuer
		userRepository *mocks.MockUserRepository
	}
}

func (s *AuthUsecaseSuite) SetupTest() {
	s.ctl = gomock.NewController(s.T())
	s.mock.oauth = mocks.NewMockOAuthClient(s.ctl)
	s.mock.tokenIssuer = mocks.NewMockTokenIssuer(s.ctl)
	s.mock.userRepository = mocks.NewMockUserRepository(s.ctl)

	s.usecase = auth_module.NewAuthUsecase(s.mock.oauth, s.mock.tokenIssuer, s.mock.userRepository)
}

func (s *AuthUsecaseSuite) TestSignIn() {
	testcases := []struct {
		desc     string
		setup    func()
		checkerr func(err error) bool
	}{
		{
			desc: "user exists",
			setup: func() {
				s.mock.oauth.EXPECT().
					IssueAccessToken(gomock.Any(), gomock.Any()).Return("code", nil)
				s.mock.oauth.EXPECT().
					GetUserInfo(gomock.Any(), gomock.Any()).Return(domain.User{}, nil)
				s.mock.userRepository.EXPECT().
					UsernameExists(gomock.Any(), gomock.Any()).Return(true, nil)
				s.mock.userRepository.EXPECT().
					FetchByUsername(gomock.Any(), gomock.Any()).Return(domain.User{}, nil)
				s.mock.tokenIssuer.EXPECT().
					Issue(gomock.Any(), gomock.Any(), gomock.Any()).Return(auth_module.Token{}, nil)
			},
			checkerr: func(err error) bool { return err == nil },
		},
		{
			desc: "user does not exist",
			setup: func() {
				s.mock.oauth.EXPECT().
					IssueAccessToken(gomock.Any(), gomock.Any()).Return("code", nil)
				s.mock.oauth.EXPECT().
					GetUserInfo(gomock.Any(), gomock.Any()).Return(domain.User{}, nil)
				s.mock.userRepository.EXPECT().
					UsernameExists(gomock.Any(), gomock.Any()).Return(false, nil)
				s.mock.userRepository.EXPECT().
					Create(gomock.Any(), gomock.Any()).Return(nil)
				s.mock.tokenIssuer.EXPECT().
					Issue(gomock.Any(), gomock.Any(), gomock.Any()).Return(auth_module.Token{}, nil)
			},
			checkerr: func(err error) bool { return err == nil },
		},
		{
			desc: "inavlid code",
			setup: func() {
				s.mock.oauth.EXPECT().
					IssueAccessToken(gomock.Any(), gomock.Any()).Return("", auth_module.ErrInvalidCode)
			},
			checkerr: func(err error) bool {
				sErr, ok := err.(status.Error)
				return ok && sErr.StatusCode == http.StatusBadRequest
			},
		},
		{
			desc: "not enough scope",
			setup: func() {
				s.mock.oauth.EXPECT().
					IssueAccessToken(gomock.Any(), gomock.Any()).Return("", auth_module.ErrNotEnoughScope)
			},
			checkerr: func(err error) bool {
				sErr, ok := err.(status.Error)
				return ok && sErr.StatusCode == http.StatusBadRequest
			},
		},
	}

	for _, tc := range testcases {
		s.Run(tc.desc, func() {
			tc.setup()

			ctx := context.Background()

			_, err := s.usecase.SignIn(ctx, &dto.SignInInput{})
			s.True(tc.checkerr(err), err)
		})
	}
}
