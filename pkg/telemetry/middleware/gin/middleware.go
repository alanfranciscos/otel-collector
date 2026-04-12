package gin

import (
	"os"
	"time"

	"github.com/alanfranciscos/otel-collector/internal/pkg/telemetry/schema/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type GinMiddlewareConfig struct {
	ServiceName string
}

func NewGinMiddlewareConfig(serviceName string) *GinMiddlewareConfig {
	return &GinMiddlewareConfig{
		ServiceName: serviceName,
	}
}

func makeFields(c *gin.Context, duration int64, serviceName string) logrus.Fields {
	statusCode := c.Writer.Status()
	logFields := logger.NewLogFields(
		serviceName,
		duration,
		logger.RequestLogField{
			Path:        c.Request.URL.Path,
			Method:      c.Request.Method,
			QueryParams: c.Request.URL.Query(),
			IP:          c.ClientIP(),
			UserAgent:   c.Request.UserAgent(),
		},
		logger.ResponseLogField{
			StatusCode: statusCode,
		},
		logger.DatabaseLogField{
			NumberOfCalls:    c.GetInt("num_db_calls"),
			NumberOfFailures: c.GetInt("num_db_failures"),
		},
	)

	if statusCode >= 400 {
		logFields.SetEvents()
	}

	for _, err := range c.Errors {
		logFields.SetErrors(err)
	}

	return logFields.ToLogrusFields()
}

func loggerMiddleware(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()

		// Proceed
		c.Next()

		duration := time.Since(start).Milliseconds()

		// fix this because break in race conditions and this is inneficient to set this every request
		logrus.SetFormatter(&logrus.JSONFormatter{})
		fields := makeFields(c, duration, serviceName)

		entry := logrus.WithContext(c.Request.Context()).WithFields(fields)

		logrus.SetOutput(os.Stdout)
		switch {
		case c.Writer.Status() >= 500:
			entry.Error("Request Completed with Server Error")
		case c.Writer.Status() >= 400:
			entry.Warn("Request Completed with Client Error")
		default:
			entry.Info("Request Completed")
		}
	}
}

func Middleware(serviceName string) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		loggerMiddleware(serviceName),
	}
}
