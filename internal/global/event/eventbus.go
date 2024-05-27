package event

import (
	"context"
)

//go:generate mockgen -source=eventbus.go -destination=../../../test/mocks/eventbus.go -package=mocks

type Handlerfunc func(ctx context.Context, topic Topic, e any) error

type Publisher interface {
	Publish(ctx context.Context, topic Topic, e any) error
}

type Subscriber interface {
	Subscribe(ctx context.Context, topic Topic, handlers ...Handlerfunc) error
}
