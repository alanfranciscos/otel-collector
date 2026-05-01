package logger

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestCustomJSONFormatter_Format_NoContext(t *testing.T) {
	formatter := &CustomJSONFormatter{serviceName: "test-service"}
	entry := &logrus.Entry{
		Logger:  logrus.New(),
		Level:   logrus.InfoLevel,
		Message: "test message",
		Context: context.Background(),
	}

	bytes, err := formatter.Format(entry)

	assert.NoError(t, err)
	assert.NotNil(t, bytes)

	var result map[string]interface{}
	err = json.Unmarshal(bytes, &result)
	assert.NoError(t, err)

	assert.Equal(t, "test-service", result["service_name"])
	assert.Equal(t, "info", result["level"])
	assert.Equal(t, "test message", result["message"])
	assert.Equal(t, "00000000000000000000000000000000", result["trace_id"])
	assert.Equal(t, "0000000000000000", result["span_id"])
	assert.NotNil(t, result["timestamp"])
}

func TestCustomJSONFormatter_Format_WithContextSpan(t *testing.T) {
	tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
	tracer := tp.Tracer("test-tracer")
	ctx, span := tracer.Start(context.Background(), "test-span")
	defer span.End()

	formatter := &CustomJSONFormatter{serviceName: "test-service"}
	entry := &logrus.Entry{
		Logger:  logrus.New(),
		Level:   logrus.ErrorLevel,
		Message: "error occurred",
		Context: ctx,
	}

	bytes, err := formatter.Format(entry)

	assert.NoError(t, err)

	var result map[string]interface{}
	json.Unmarshal(bytes, &result)

	assert.NotEqual(t, "00000000000000000000000000000000", result["trace_id"])
	assert.NotEqual(t, "0000000000000000", result["span_id"])
	assert.Equal(t, span.SpanContext().TraceID().String(), result["trace_id"])
	assert.Equal(t, span.SpanContext().SpanID().String(), result["span_id"])
}

func TestCustomJSONFormatter_Format_MergedData(t *testing.T) {
	formatter := &CustomJSONFormatter{serviceName: "test-service"}
	entry := &logrus.Entry{
		Logger:  logrus.New(),
		Level:   logrus.DebugLevel,
		Message: "debug info",
		Context: context.Background(),
		Data: logrus.Fields{
			"custom_field": "custom_value",
			"user_id":      123,
		},
	}

	bytes, err := formatter.Format(entry)
	assert.NoError(t, err)

	var result map[string]interface{}
	json.Unmarshal(bytes, &result)

	assert.Equal(t, "custom_value", result["custom_field"])
	assert.Equal(t, float64(123), result["user_id"])
	assert.Equal(t, "debug info", result["message"])
}

func TestCustomJSONFormatter_Format_Error(t *testing.T) {
	formatter := &CustomJSONFormatter{serviceName: "test-service"}
	
	// Inject a channel into Data, which cannot be JSON marshaled
	unmarshalableData := make(chan int)
	entry := &logrus.Entry{
		Logger:  logrus.New(),
		Level:   logrus.InfoLevel,
		Context: context.Background(),
		Data: logrus.Fields{
			"bad_data": unmarshalableData,
		},
	}

	bytes, err := formatter.Format(entry)

	assert.Error(t, err)
	assert.Nil(t, bytes)
}
