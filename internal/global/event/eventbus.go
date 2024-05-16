package event

import (
	"context"
)

//go:generate mockgen -source=eventbus.go -destination=../../../test/mocks/eventbus.go -package=mocks

type Topic string

const (
	TopicSubmission Topic = "submission"
)

type Handlerfunc func(event any) error

type Publisher interface {
	Publish(ctx context.Context, topic Topic, event any) error
}

type Subscriber interface {
	Subscribe(ctx context.Context, topic Topic, f Handlerfunc) error
}
