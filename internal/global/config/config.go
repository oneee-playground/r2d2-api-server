package config

import "time"

type Config struct {
	ServerConfig ServerConfig
	JWTConfig    JWTConfig
	GitHubConfig GitHubConfig
	AWSConfig    AWSConfig
	RedisConfig  RedisConfig
	MYSQLConfig  MYSQLConfig
	EmailConfig  EmailConfig
}

type ServerConfig struct {
	Port int
}

type JWTConfig struct {
	Secret string
}

type GitHubConfig struct {
	ClientID     string
	ClientSecret string
}

type AWSConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string

	SQSConfig SQSConfig
}

type SQSConfig struct {
	JobQueueURL             string
	SubmissionEventQueueURL string
	BuildEventQueueURL      string
	TestEventQueueURL       string

	PollInterval time.Duration
}

type RedisConfig struct {
	Addr  string
	DBNum int
}

type MYSQLConfig struct {
	Addr string
	Pass string
}

type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string

	FromAddr string
}

func GetServerConfig() ServerConfig { return loaded.ServerConfig }
func GetJWTConfig() JWTConfig       { return loaded.JWTConfig }
func GetGitHubConfig() GitHubConfig { return loaded.GitHubConfig }
func GetAWSConfig() AWSConfig       { return loaded.AWSConfig }
func GetRedisConfig() RedisConfig   { return loaded.RedisConfig }
func GetMYSQLConfig() MYSQLConfig   { return loaded.MYSQLConfig }
func GetEmailConfig() EmailConfig   { return loaded.EmailConfig }
