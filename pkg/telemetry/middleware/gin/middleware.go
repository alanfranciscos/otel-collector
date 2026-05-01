package gin

import (
	"os"
	"time"

	"github.com/alanfranciscos/otel-collector/internal/pkg/telemetry/schema/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type GinMiddleware interface {
	Middleware() []gin.HandlerFunc
}

type GinMiddlewareConfig struct {
}

func NewGinMiddlewareConfig() *GinMiddlewareConfig {
	return &GinMiddlewareConfig{}
}

func makeFields(c *gin.Context, duration int64) logrus.Fields {
	statusCode := c.Writer.Status()
	logFields := logger.NewLogFields(
		time.Now().Format(time.RFC3339),
	)

	requestField := &logger.RequestLogField{
		Path:        c.Request.URL.Path,
		Method:      c.Request.Method,
		QueryParams: c.Request.URL.Query(),
		IP:          c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
	}

	responseField := &logger.ResponseLogField{
		StatusCode: statusCode,
	}

	databaseField := &logger.DatabaseLogField{
		NumberOfCalls:    c.GetInt("num_db_calls"),
		NumberOfFailures: c.GetInt("num_db_failures"),
	}

	logFields.
		SetDurationMs(duration).
		SetRequest(requestField).
		SetResponse(responseField).
		SetDatabase(databaseField)

	if statusCode >= 400 {
		logFields.SetEvents()
	}

	for _, err := range c.Errors {
		logFields.SetErrors(err)
	}

	return logFields.ToLogrusFields()
}

func loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()

		// Proceed
		c.Next()

		duration := time.Since(start).Milliseconds()

		fields := makeFields(c, duration)

		entry := logrus.WithContext(c.Request.Context()).WithFields(fields)

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

func (g GinMiddlewareConfig) Middleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		otelgin.Middleware(os.Getenv("OTEL_SERVICE_NAME")),
		loggerMiddleware(),
	}
}
