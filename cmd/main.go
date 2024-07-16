package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt"
	"github.com/oneee-playground/r2d2-api-server/internal/global/config"
	"github.com/oneee-playground/r2d2-api-server/internal/global/event"
	lambda_module "github.com/oneee-playground/r2d2-api-server/internal/infra/aws/lambda"
	sqs_module "github.com/oneee-playground/r2d2-api-server/internal/infra/aws/sqs"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/data/ent/datasource"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/data/ent/model"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/data/ent/repository"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/data/redis"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/email"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/github"
	httproute "github.com/oneee-playground/r2d2-api-server/internal/infra/http"
	"github.com/oneee-playground/r2d2-api-server/internal/infra/http/handler"
	jwt_token "github.com/oneee-playground/r2d2-api-server/internal/infra/jwt"
	auth_module "github.com/oneee-playground/r2d2-api-server/internal/module/auth"
	event_module "github.com/oneee-playground/r2d2-api-server/internal/module/event"
	exec_module "github.com/oneee-playground/r2d2-api-server/internal/module/exec"
	resource_module "github.com/oneee-playground/r2d2-api-server/internal/module/resource"
	section_module "github.com/oneee-playground/r2d2-api-server/internal/module/section"
	submission_module "github.com/oneee-playground/r2d2-api-server/internal/module/submission"
	task_module "github.com/oneee-playground/r2d2-api-server/internal/module/task"
	user_module "github.com/oneee-playground/r2d2-api-server/internal/module/user"
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidislock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	ctx := context.Background()

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	defer stop()

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(os.Stdout), zap.DebugLevel,
	))

	if err := config.Load(ctx, &config.EnvLoader{}); err != nil {
		logger.Panic("loading config failed", zap.Error(err))
	}

	emailConfig := config.GetEmailConfig()

	// etc.
	var (
		tokenManager = jwt_token.NewManager(jwt.SigningMethodHS256, []byte(config.GetJWTConfig().Secret))
		oauthClient  = github.NewClient(http.DefaultClient, logger, config.GetGitHubConfig().ClientID, config.GetGitHubConfig().ClientSecret)
		emailSender  = email.NewGomailSender(logger, email.GomailOptions{
			Host:     emailConfig.Host,
			Port:     emailConfig.Port,
			Username: emailConfig.Username,
			Password: emailConfig.Password,
			FromAddr: emailConfig.FromAddr,
		})
	)

	awsConfig := config.GetAWSConfig()

	awsConf := aws.Config{
		Region:      awsConfig.Region,
		Credentials: credentials.NewStaticCredentialsProvider(awsConfig.AccessKeyID, awsConfig.SecretAccessKey, ""),
	}

	lambdaClient := lambda.NewFromConfig(awsConf)
	imageBuilder := lambda_module.NewLambdaImageBuilder(lambdaClient, logger)

	sqsClient := sqs.NewFromConfig(awsConf)

	jobQueue := sqs_module.NewSQSJobQueue(sqsClient, logger, awsConfig.SQSConfig.JobQueueURL)
	eventBus := sqs_module.NewSQSEventBus(sqsClient, logger, map[event.Topic]sqs_module.QueueConfig{
		event.TopicBuild: {
			URL:          awsConfig.SQSConfig.BuildEventQueueURL,
			PollInterval: awsConfig.SQSConfig.PollInterval,
		},
		event.TopicSubmission: {
			URL:          awsConfig.SQSConfig.SubmissionEventQueueURL,
			PollInterval: awsConfig.SQSConfig.PollInterval,
		},
		event.TopicTest: {
			URL:          awsConfig.SQSConfig.TestEventQueueURL,
			PollInterval: awsConfig.SQSConfig.PollInterval,
		},
	})

	go eventBus.Listen(ctx)

	mysqlConf := config.GetMYSQLConfig()

	entClient, err := model.Open("mysql", fmt.Sprintf("root:%s@tcp(%s)/r2d2?parseTime=true", mysqlConf.Pass, mysqlConf.Addr))
	if err != nil {
		logger.Panic("failed to open client", zap.Error(err))
	}

	datasource := datasource.New(entClient)
	defer entClient.Close()

	if err := datasource.Migrate(ctx); err != nil {
		logger.Panic("failed to migrate db", zap.Error(err))
	}

	var (
		eventRepo      = repository.NewEventRepository(datasource)
		resourceRepo   = repository.NewResourceRepository(datasource)
		sectionRepo    = repository.NewSectionRepository(datasource)
		submissionRepo = repository.NewSubmissionRepository(datasource)
		taskRepo       = repository.NewTaskRepository(datasource)
		userRepo       = repository.NewUserRepository(datasource)
	)

	rueidisOpts := rueidis.ClientOption{
		InitAddress: []string{config.GetRedisConfig().Addr},
		SelectDB:    config.GetRedisConfig().DBNum,
	}

	redisClient, err := rueidis.NewClient(rueidisOpts)
	if err != nil {
		logger.Panic("failed to initialize redis client", zap.Error(err))
	}
	defer redisClient.Close()

	lock, err := rueidislock.NewLocker(rueidislock.LockerOption{ClientOption: rueidisOpts})
	if err != nil {
		logger.Panic("failed to initizlize redis lock client", zap.Error(err))
	}
	defer lock.Close()

	txLocker := redis.NewLocker(lock)
	execContextStorage := redis.NewExecContextStroage(redisClient)

	var (
		authUsecase       = auth_module.NewAuthUsecase(oauthClient, tokenManager, userRepo, txLocker)
		resourceUsecase   = resource_module.NewResourceUsecase(resourceRepo, taskRepo, txLocker)
		sectionUsecase    = section_module.NewSectionUsecase(sectionRepo, taskRepo, txLocker)
		submissionUsecase = submission_module.NewSubmissionUsecase(taskRepo, submissionRepo, eventRepo, eventBus, txLocker)
		taskUsecase       = task_module.NewTaskUsecase(taskRepo, txLocker)
		userUsecase       = user_module.NewUserUsecase(userRepo)
		eventUsecase      = event_module.NewEventUsecase(eventRepo)

		execEventHandler  = exec_module.NewEventHandler(submissionRepo, sectionRepo, resourceRepo, eventBus, jobQueue, imageBuilder, execContextStorage)
		eventEventHandler = event_module.NewEventHandler(emailSender, userRepo, eventRepo)
	)

	if err := execEventHandler.Register(ctx, eventBus); err != nil {
		logger.Panic("registering exec event handler failed", zap.Error(err))
	}
	if err := eventEventHandler.Register(ctx, eventBus); err != nil {
		logger.Panic("registering event event handler failed", zap.Error(err))
	}

	router := &httproute.Router{
		Engine:            gin.New(),
		TokenDecoder:      tokenManager,
		RequestLogger:     logger,
		ErrorLogger:       logger,
		EventHandler:      handler.NewEventHandler(eventUsecase),
		ResourceHandler:   handler.NewResourceHandler(resourceUsecase),
		SectionHandler:    handler.NewSectionHandler(sectionUsecase),
		SubmissionHandler: handler.NewSubmissionHandler(submissionUsecase),
		TaskHandler:       handler.NewTaskHandler(taskUsecase),
		UserHandler:       handler.NewUserHandler(userUsecase),
		AuthHandler:       handler.NewAuthHandler(authUsecase),
	}
	router.Build()

	router.Serve(config.GetServerConfig().Port)
}
