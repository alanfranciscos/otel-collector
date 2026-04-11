package gin

import (
	"os"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func makeFields(c *gin.Context, duration int64) logrus.Fields {
	httpField := logrus.Fields{
		"path":         c.Request.URL.Path,
		"method":       c.Request.Method,
		"query_params": c.Request.URL.RawQuery,
		"ip":           c.ClientIP(),
		"user_agent":   c.Request.UserAgent(),
	}

	responseField := logrus.Fields{
		"status_code": c.Writer.Status(),
	}

	databaseField := logrus.Fields{
		"num_calls":    c.GetInt("num_db_calls"),
		"num_failures": c.GetInt("num_db_failures"),
	}

	errorsJson := c.Errors.JSON()
	typeErrors := reflect.TypeOf(errorsJson)
	var errors []interface{}
	switch {
	case errorsJson == nil:
		errors = []interface{}{}
	case typeErrors.Kind() != reflect.Slice:
		errors = []interface{}{errorsJson}
	default:
		errors = errorsJson.([]interface{})
	}
	errorField := logrus.Fields{
		"error": errors,
	}

	fields := logrus.Fields{
		"duration_ms": duration,
		"request":     httpField,
		"response":    responseField,
		"database":    databaseField,
		"error":       errorField,
	}

	return fields
}

func loggerMiddleware(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()

		// Proceed
		c.Next()

		duration := time.Since(start).Milliseconds()

		logrus.SetFormatter(&logrus.JSONFormatter{})
		fields := makeFields(c, duration)

		entry := logrus.WithContext(c.Request.Context()).WithFields(fields)

		logrus.SetOutput(os.Stdout)
		switch {
		case c.Writer.Status() >= 500:
			entry.Error("[" + serviceName + "] - Request Completed with Server Error")
		case c.Writer.Status() >= 400:
			entry.Warn("[" + serviceName + "] - Request Completed with Client Error")
		default:
			entry.Info("[" + serviceName + "] - Request Completed")
		}

	}
}

func Middleware(serviceName string) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		loggerMiddleware(serviceName),
	}
}
