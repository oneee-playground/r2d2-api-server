package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	auth_module "github.com/oneee-playground/r2d2-api-server/internal/module/auth"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var requiredScopes = []string{"read:user", "user:email"}

// Client is github-specific oauth client.
type Client struct {
	httpClient *http.Client
	logger     *zap.Logger

	clientID, clientSecret string
}

var _ auth_module.OAuthClient = (*Client)(nil)

func NewClient(client *http.Client, logger *zap.Logger, clientID, clientSecret string) *Client {
	return &Client{
		httpClient:   client,
		logger:       logger,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

type accessTokenResponse struct {
	Token string `json:"access_token"`
	Scope string `json:"scope"`
	// Error will be filled when github responses with error.
	Error string `json:"error"`
}

func (c *Client) IssueAccessToken(ctx context.Context, code string) (string, error) {
	url := fmt.Sprintf(
		"https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s",
		c.clientID, c.clientSecret, code,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return "", errors.Wrap(err, "creating request")
	}

	var decoded accessTokenResponse
	if err := sendRequest(c.httpClient, req, &decoded); err != nil {
		return "", err
	}

	if decoded.Error != "" {
		if decoded.Error == "bad_verification_code" {
			return "", auth_module.ErrInvalidCode
		}

		return "", errors.Errorf("unexpected error response from github: %s", decoded.Error)
	}

	if !scopeValid(requiredScopes, decoded.Scope) {
		return "", auth_module.ErrNotEnoughScope
	}

	return decoded.Token, nil
}

type githubUserInfoResponse struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"`
}

func (c *Client) GetUserInfo(ctx context.Context, token string) (domain.User, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return domain.User{}, errors.Wrap(err, "creating request")
	}

	req.Header.Set("Authorization", "Bearer "+token)

	var decoded githubUserInfoResponse
	if err := sendRequest(c.httpClient, req, &decoded); err != nil {
		return domain.User{}, err
	}

	//nolint:exhaustruct
	user := domain.User{
		Username:   decoded.Login,
		Email:      decoded.Email,
		ProfileURL: decoded.AvatarURL,
	}

	return user, nil
}

// sendRequest sends request via httpClient and decodes response body with specified type T.
func sendRequest[T any](httpClient *http.Client, req *http.Request, val *T) error {
	req.Header.Set("Accept", "application/json")

	res, err := httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "performing request")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.Errorf("status code is not 200, given: %d", res.StatusCode)
	}

	if err := json.NewDecoder(res.Body).Decode(val); err != nil {
		return errors.Wrap(err, "could not decode response body")
	}

	return nil
}

// scopeValid checks if splitted scopes contains all the required scopes.
func scopeValid(required []string, scope string) bool {
	scopes := strings.Split(scope, ",")

	i := 0

	for _, scope := range scopes {
		if slices.Contains(required, scope) {
			i++
		}
	}

	return i >= len(required)
}
