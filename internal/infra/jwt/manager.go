package jwt

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/oneee-playground/r2d2-api-server/internal/global/auth"
	auth_module "github.com/oneee-playground/r2d2-api-server/internal/module/auth"
	"github.com/pkg/errors"
)

type jwtClaims struct {
	jwt.StandardClaims
	auth.Payload
}

type Manager struct {
	signingMethod jwt.SigningMethod
	secret        any

	parser *jwt.Parser
}

var _ auth_module.TokenIssuer = (*Manager)(nil)
var _ auth_module.TokenDecoder = (*Manager)(nil)

func NewManager(method jwt.SigningMethod, secret any) *Manager {
	return &Manager{
		signingMethod: method,
		secret:        secret,
		parser:        &jwt.Parser{SkipClaimsValidation: true},
	}
}

func (m *Manager) Issue(ctx context.Context, payload auth.Payload, exp time.Time) (auth_module.Token, error) {
	claims := jwtClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp.Unix(),
		},
		Payload: payload,
	}

	token := jwt.NewWithClaims(m.signingMethod, claims)

	raw, err := token.SignedString(m.secret)
	if err != nil {
		return auth_module.Token{}, errors.Wrap(err, "creating token")
	}

	authToken := auth_module.Token{
		Payload:   payload,
		ExpiresAt: exp,
		Raw:       raw,
	}

	return authToken, nil
}

func (m *Manager) Decode(_ context.Context, raw string) (auth_module.Token, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		if t.Method != m.signingMethod {
			return nil, jwt.NewValidationError("method does not match", jwt.ValidationErrorMalformed)
		}

		return m.secret, nil
	}

	var claims jwtClaims

	_, err := m.parser.ParseWithClaims(raw, &claims, keyFunc)
	if err != nil {
		var vErr *jwt.ValidationError
		if errors.As(err, &vErr) {
			const tokenInvalid = jwt.ValidationErrorSignatureInvalid | jwt.ValidationErrorMalformed
			if vErr.Errors&tokenInvalid != 0 {
				return auth_module.Token{}, auth_module.ErrTokenInvalid
			}
		}

		return auth_module.Token{}, errors.Wrap(err, "parsing token")
	}

	if isExpired(claims, time.Now()) {
		return auth_module.Token{}, auth_module.ErrTokenExpired
	}

	authToken := auth_module.Token{
		Payload:   claims.Payload,
		ExpiresAt: time.Unix(claims.ExpiresAt, 0),
		Raw:       raw,
	}

	return authToken, nil
}

func isExpired(claims jwtClaims, cmp time.Time) bool {
	return cmp.After(time.Unix(claims.ExpiresAt, 0))
}
