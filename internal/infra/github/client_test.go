package github

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	auth_module "github.com/oneee-playground/r2d2-api-server/internal/module/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestGitHubClientSuote(t *testing.T) {
	suite.Run(t, new(GitHubClientSuote))
}

type GitHubClientSuote struct {
	suite.Suite

	mockTransport *httpmock.MockTransport

	client *Client
}

func (s *GitHubClientSuote) SetupTest() {
	s.mockTransport = httpmock.NewMockTransport()

	httpClient := &http.Client{Transport: s.mockTransport}

	s.client = NewClient(httpClient, "", "")
}

func (s *GitHubClientSuote) TestIssueAccessToken() {
	defer s.mockTransport.Reset()

	validCode := "code"

	s.mockTransport.RegisterResponder(http.MethodPost,
		"https://github.com/login/oauth/access_token",
		func(r *http.Request) (*http.Response, error) {
			if r.URL.Query().Get("code") == validCode {
				return httpmock.NewStringResponse(http.StatusOK, `
				{
					"access_token": "token",
					"scope": "read:user,user:email"
				}
				`), nil
			}

			return httpmock.NewStringResponse(http.StatusOK, `
			{
				"error": "bad_verification_code"
			}
			`), nil
		},
	)

	testcases := []struct {
		desc string
		code string
		err  error
	}{
		{
			desc: "valid code",
			code: validCode,
			err:  nil,
		},
		{
			desc: "invalid code",
			code: "invalid",
			err:  auth_module.ErrInvalidCode,
		},
	}

	for _, tc := range testcases {
		s.Run(tc.desc, func() {
			ctx := context.Background()

			_, err := s.client.IssueAccessToken(ctx, tc.code)
			s.Equal(tc.err, err)
		})
	}
}

func TestScopeValid(t *testing.T) {
	required := []string{"hi", "there", "test"}

	assert.True(t, scopeValid(required, strings.Join(required, ",")))
	assert.True(t, scopeValid(required, strings.Join(append(required, "aef"), ",")))
	assert.False(t, scopeValid(required, strings.Join(required, ",")[:2]))
}
