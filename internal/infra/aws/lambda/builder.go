package lambda

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	exec_module "github.com/oneee-playground/r2d2-api-server/internal/module/exec"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const builderFuncName = "r2d2-image-builder"

type LambdaImageBuilder struct {
	client *lambda.Client
	logger *zap.Logger
}

var _ exec_module.ImageBuilder = (*LambdaImageBuilder)(nil)

func NewLambdaImageBuilder(client *lambda.Client, logger *zap.Logger) *LambdaImageBuilder {
	return &LambdaImageBuilder{
		client: client,
		logger: logger,
	}
}

func (b *LambdaImageBuilder) RequestBuild(ctx context.Context, opts exec_module.BuildOpts) error {
	payload, err := json.Marshal(opts)
	if err != nil {
		return errors.Wrap(err, "marshalling payload")
	}

	input := &lambda.InvokeInput{
		FunctionName:   aws.String(builderFuncName),
		InvocationType: types.InvocationTypeEvent,
		Payload:        payload,
	}

	output, err := b.client.Invoke(ctx, input)
	if err != nil {
		return errors.Wrap(err, "invoking aws lambda function")
	}

	if output.StatusCode != http.StatusAccepted {
		err := errors.New("invocation is not accepted")
		b.logger.Error(err.Error(), zap.Any("output", output))

		return err
	}

	b.logger.Info("requested image build", zap.String("id", opts.ID.String()))

	return nil
}
