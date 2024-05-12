package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
)

type Payload struct {
	UserID uuid.UUID       `json:"userId"`
	Role   domain.UserRole `json:"role"`
}

type _payloadKey struct{}

// Inject injects payload into given context.
func Inject(ctx context.Context, payload Payload) context.Context {
	return context.WithValue(ctx, _payloadKey{}, &payload)
}

// Extract extracts payload from given context.
func Extract(ctx context.Context) (payload *Payload, found bool) {
	v, ok := ctx.Value(_payloadKey{}).(*Payload)
	return v, ok
}

// MustExtract is like Extract. But it will panic if payload is not found.
func MustExtract(ctx context.Context) (payload *Payload) {
	v, ok := ctx.Value(_payloadKey{}).(*Payload)
	if !ok {
		panic("payload not found from context")
	}
	return v
}
