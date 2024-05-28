package event

import (
	"context"

	"github.com/pkg/errors"
)

//go:generate mockgen -source=eventbus.go -destination=../../../test/mocks/eventbus.go -package=mocks

type HandlerFunc func(ctx context.Context, topic Topic, e any) error

type Publisher interface {
	Publish(ctx context.Context, topic Topic, e any) error
}

var NoErrSkipHandler = errors.New("this isn't event for the handler. skip")

type Subscriber interface {
	Subscribe(ctx context.Context, topic Topic, handlers ...HandlerFunc) error
}
