package sqs

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/oneee-playground/r2d2-api-server/internal/global/event"
	"github.com/pkg/errors"
)

type QueueConfig struct {
	URL          string
	PollInterval time.Duration
}

type SQSEventBus struct {
	client *sqs.Client

	// topicMapping maps topic into queue.
	topicMapping map[event.Topic]QueueConfig
	handlers     map[event.Topic][]event.HandlerFunc
}

var (
	_ event.Subscriber = (*SQSEventBus)(nil)
	_ event.Publisher  = (*SQSEventBus)(nil)
)

// topicMapping should not be nil.
func NewSQSEventBus(client *sqs.Client, topicMapping map[event.Topic]QueueConfig) *SQSEventBus {
	return &SQSEventBus{
		client:       client,
		topicMapping: topicMapping,
		handlers:     make(map[event.Topic][]event.HandlerFunc),
	}
}

func (b *SQSEventBus) Publish(ctx context.Context, topic event.Topic, e any) error {
	payload, err := json.Marshal(e)
	if err != nil {
		return errors.Wrap(err, "marshalling payload")
	}

	queueURL := b.topicMapping[topic].URL

	input := &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(string(payload)),
	}

	if _, err := b.client.SendMessage(ctx, input); err != nil {
		return errors.Wrap(err, "sending message")
	}

	return nil
}

func (b *SQSEventBus) Subscribe(ctx context.Context, topic event.Topic, handlers ...event.HandlerFunc) error {
	b.handlers[topic] = append(b.handlers[topic], handlers...)
	return nil
}

// Listen subscribes all topics and periodically polls messages.
// It is required to call it within seperate goroutine since it blocks the flow.
func (b *SQSEventBus) Listen(ctx context.Context) {
	var wg sync.WaitGroup

	for topic, queue := range b.topicMapping {
		wg.Add(1)
		go b.listenTopic(ctx, &wg, topic, queue)
	}

	wg.Wait()
}

func (b *SQSEventBus) listenTopic(ctx context.Context, wg *sync.WaitGroup, topic event.Topic, queue QueueConfig) {
	ticker := time.NewTicker(queue.PollInterval)

	input := &sqs.ReceiveMessageInput{
		QueueUrl: &queue.URL,
		// TODO: Are ther any more things to add?
	}

	defer wg.Done()
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// TODO: Maybe log the behavior.
			return
		case <-ticker.C:
			output, err := b.client.ReceiveMessage(ctx, input)
			if err != nil {
				// TODO: log the error.
				continue
			}

			for _, message := range output.Messages {
				payload := []byte(*message.Body)

				for _, f := range b.handlers[topic] {
					err := f(ctx, topic, payload)
					if err == event.NoErrSkipHandler {
						continue
					}

					if err != nil {
						// TODO: log the error.
					}
				}
			}
		}
	}
}
