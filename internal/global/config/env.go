package config

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type EnvLoader struct{}

var _ Loader = (*EnvLoader)(nil)

func (el *EnvLoader) Fill(ctx context.Context, conf *Config) error {
	confFuncs := []func(conf *Config) error{
		el.serverConfig, el.jwtConfig, el.gitHubConfig,
		el.awsConfig, el.redisConfig, el.emailConfig, el.mysqlConfig,
	}

	for _, f := range confFuncs {
		if err := f(conf); err != nil {
			return err
		}
	}

	return nil
}

func (el *EnvLoader) serverConfig(conf *Config) error {
	serverConf := ServerConfig{}

	port, err := strconv.ParseInt(os.Getenv("SERVER_PORT"), 10, 64)
	if err != nil {
		return errors.Wrap(err, "could not parse port")
	}

	serverConf.Port = int(port)

	conf.ServerConfig = serverConf
	return nil
}

func (el *EnvLoader) jwtConfig(conf *Config) error {
	jwtConf := JWTConfig{}

	jwtConf.Secret = os.Getenv("JWT_SECRET")

	conf.JWTConfig = jwtConf
	return nil
}

func (el *EnvLoader) gitHubConfig(conf *Config) error {
	githubConf := GitHubConfig{}

	githubConf.ClientID = os.Getenv("GITHUB_CLIENT_ID")
	githubConf.ClientSecret = os.Getenv("GITHUB_CLIENT_SECRET")

	conf.GitHubConfig = githubConf
	return nil
}

func (el *EnvLoader) awsConfig(conf *Config) error {
	awsConf := AWSConfig{}

	awsConf.Region = os.Getenv("AWS_REGION")
	awsConf.SQSConfig = SQSConfig{
		JobQueueURL:             os.Getenv("AWS_SQS_JOB_QUEUE_URL"),
		SubmissionEventQueueURL: os.Getenv("AWS_SQS_SUBMISSION_EVENT_QUEUE_URL"),
		BuildEventQueueURL:      os.Getenv("AWS_SQS_BUILD_EVENT_QUEUE_URL"),
		TestEventQueueURL:       os.Getenv("AWS_SQS_TEST_EVENT_QUEUE_URL"),
	}

	pollIntervalRaw := os.Getenv("AWS_SQS_POLL_INTERVAL_SECOND")

	pollInterval, err := strconv.ParseInt(pollIntervalRaw, 10, 64)
	if err != nil {
		return errors.Wrap(err, "parsing poll interval")
	}

	awsConf.SQSConfig.PollInterval = time.Duration(pollInterval) * time.Second

	conf.AWSConfig = awsConf
	return nil
}

func (el *EnvLoader) redisConfig(conf *Config) error {
	redisConf := RedisConfig{}

	redisConf.Addr = os.Getenv("REDIS_ADDR")

	redisDBNumRaw := os.Getenv("REDIS_DB_NUM")

	redisDBNum, err := strconv.ParseInt(redisDBNumRaw, 10, 64)
	if err != nil {
		return errors.Wrap(err, "parsing redis db num")
	}

	redisConf.DBNum = int(redisDBNum)

	conf.RedisConfig = redisConf
	return nil
}

func (el *EnvLoader) mysqlConfig(conf *Config) error {
	mysqlConf := MYSQLConfig{}

	mysqlConf.Addr = os.Getenv("MYSQL_ADDR")
	mysqlConf.Pass = os.Getenv("MYSQL_PASSWORD")

	conf.MYSQLConfig = mysqlConf
	return nil
}

func (el *EnvLoader) emailConfig(conf *Config) error {
	emailConf := EmailConfig{}

	emailConf.FromAddr = os.Getenv("EMAIL_FROM_ADDR")
	emailConf.Host = os.Getenv("EMAIL_HOST")
	emailConf.Password = os.Getenv("EMAIL_PASSWORD")
	emailConf.Username = os.Getenv("EMAIL_USERNAME")

	emailPortRaw := os.Getenv("EMAIL_PORT")

	emailPort, err := strconv.ParseInt(emailPortRaw, 10, 64)
	if err != nil {
		return errors.Wrap(err, "parsing redis db num")
	}

	emailConf.Port = int(emailPort)

	conf.EmailConfig = emailConf
	return nil
}
