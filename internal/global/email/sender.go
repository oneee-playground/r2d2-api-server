package email

import (
	"bytes"
	"context"
)

//go:generate mockgen -source=sender.go -destination=../../../test/mocks/email.go -package=mocks

type Sender interface {
	Send(ctx context.Context, address, subject string, content *bytes.Buffer) error
}
