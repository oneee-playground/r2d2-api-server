package sqs

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	exec_module "github.com/oneee-playground/r2d2-api-server/internal/module/exec"
)

type SQSJobQueue struct {
	client *sqs.Client
	logger *zap.Logger

	queueURL string
}

var _ exec_module.JobQueue = (*SQSJobQueue)(nil)

func NewSQSJobQueue(client *sqs.Client, logger *zap.Logger, queueURL string) *SQSJobQueue {
	return &SQSJobQueue{
		client:   client,
		logger:   logger,
		queueURL: queueURL,
	}
}

func (q *SQSJobQueue) Append(ctx context.Context, job *exec_module.Job) error {
	payload, err := json.Marshal(job)
	if err != nil {
		return errors.Wrap(err, "marshalling payload")
	}

	input := &sqs.SendMessageInput{
		QueueUrl:    aws.String(q.queueURL),
		MessageBody: aws.String(string(payload)),
	}

	if _, err := q.client.SendMessage(ctx, input); err != nil {
		return errors.Wrap(err, "sending message")
	}

	q.logger.Info("job enqueued",
		zap.String("taskID", job.TaskID.String()),
		zap.String("submissionID", job.Submission.ID.String()),
	)

	return nil
}
