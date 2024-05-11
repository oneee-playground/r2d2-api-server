package jwt

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	auth_module "github.com/oneee-playground/r2d2-api-server/internal/module/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestJWTManagerSuite(t *testing.T) {
	suite.Run(t, new(JWTManagerSuite))
}

type JWTManagerSuite struct {
	suite.Suite

	signingMethod jwt.SigningMethod
	secret        any

	jwtManager *Manager
}

func (s *JWTManagerSuite) SetupTest() {
	s.signingMethod = jwt.SigningMethodHS256
	s.secret = []byte("secret")

	s.jwtManager = NewManager(s.signingMethod, s.secret)
}

func (s *JWTManagerSuite) TestIssue() {
	testcases := []struct {
		desc    string
		wantErr bool
	}{
		{
			desc:    "success",
			wantErr: false,
		},
	}

	payload := auth_module.TokenPayload{
		UserID: uuid.New(),
		Role:   domain.RoleMember,
	}

	exp := time.Now()

	for _, tc := range testcases {
		s.Run(tc.desc, func() {
			ctx := context.Background()

			token, err := s.jwtManager.Issue(ctx, payload, exp)
			if tc.wantErr {
				s.Error(err)
				return
			}

			s.NoError(err)
			s.Equal(payload, token.Payload)
			s.Equal(exp, token.ExpiresAt)
		})
	}
}

func (s *JWTManagerSuite) TestDecode() {
	invalidSignMethod := jwt.SigningMethodNone
	s.Require().NotEqual(s.signingMethod, invalidSignMethod)

	createToken := func(method jwt.SigningMethod, secret any, claims jwt.Claims) string {
		raw, err := jwt.NewWithClaims(method, claims).SignedString(secret)
		s.Require().NoError(err)

		return raw
	}

	defaultClaims := jwtClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(100 * time.Hour).Unix(),
		},
		TokenPayload: auth_module.TokenPayload{
			UserID: uuid.New(),
			Role:   domain.RoleAdmin,
		},
	}

	expiredClaims := defaultClaims
	expiredClaims.ExpiresAt = time.Now().Add(-1 * time.Hour).Unix()

	testcases := []struct {
		desc    string
		token   string
		wantErr bool
		err     error
	}{
		{
			desc:    "success",
			token:   createToken(s.signingMethod, s.secret, defaultClaims),
			wantErr: false,
		},
		{
			desc:    "expired token",
			token:   createToken(s.signingMethod, s.secret, expiredClaims),
			wantErr: true,
			err:     auth_module.ErrTokenExpired,
		},
		{
			desc:    "invalid signature",
			token:   createToken(s.signingMethod, []byte("invalid"), defaultClaims),
			wantErr: true,
			err:     auth_module.ErrTokenInvalid,
		},
		{
			desc:    "invalid signing method",
			token:   createToken(invalidSignMethod, jwt.UnsafeAllowNoneSignatureType, defaultClaims),
			wantErr: true,
			err:     auth_module.ErrTokenInvalid,
		},
		{
			desc:    "malformed token",
			token:   "aeflakefjhdl.aefaekfjah.sefe",
			wantErr: true,
			err:     auth_module.ErrTokenInvalid,
		},
	}

	for _, tc := range testcases {
		s.Run(tc.desc, func() {
			ctx := context.Background()

			token, err := s.jwtManager.Decode(ctx, tc.token)
			if tc.wantErr {
				s.Error(err)

				if tc.err != nil {
					s.ErrorIs(err, tc.err)
				}

				return
			}

			s.NoError(err)
			s.Equal(defaultClaims.TokenPayload, token.Payload)
		})
	}
}

func TestIsExpired(t *testing.T) {
	now := time.Now()

	claims := jwtClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Unix(),
		},
	}

	assert.False(t, isExpired(claims, now.Add(-1*time.Hour)))
	assert.True(t, isExpired(claims, now.Add(time.Hour)))
}
