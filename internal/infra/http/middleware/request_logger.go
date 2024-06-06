package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RequestLogger struct {
	logger *zap.Logger
}

func NewRequestLogger(logger *zap.Logger) *RequestLogger {
	return &RequestLogger{logger: logger}
}

func (l *RequestLogger) Log(c *gin.Context) {
	start := time.Now()

	// Process request.
	c.Next()

	end := time.Now()
	latency := end.Sub(start)

	l.logger.Info("handled request",
		zap.Time("timestamp", end),
		zap.Int("status", c.Writer.Status()),
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("query", c.Request.URL.RawQuery),
		zap.Duration("latency", latency),
		zap.String("ip", c.ClientIP()),
	)
}
