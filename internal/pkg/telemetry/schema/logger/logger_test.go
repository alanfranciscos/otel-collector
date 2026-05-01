package logger

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewLogFields(t *testing.T) {
	timestamp := time.Now().Format(time.RFC3339)
	fields := NewLogFields(timestamp)

	assert.NotNil(t, fields)
	assert.Equal(t, timestamp, fields.Timestamp)
	assert.Empty(t, fields.Service)
	assert.Empty(t, fields.TraceID)
}

func TestLogFields_Setters(t *testing.T) {
	fields := NewLogFields("2026-05-01T12:00:00Z")

	fields.SetService("test-service").
		SetTraceID("trace-123").
		SetSpanID("span-123").
		SetDurationMs(150)

	assert.Equal(t, "test-service", fields.Service)
	assert.Equal(t, "trace-123", fields.TraceID)
	assert.Equal(t, "span-123", fields.SpanID)
	assert.Equal(t, int64(150), fields.DurationMs)

	req := &RequestLogField{Method: "GET", Path: "/"}
	res := &ResponseLogField{StatusCode: 200}
	db := &DatabaseLogField{NumberOfCalls: 1}

	fields.SetRequest(req).
		SetResponse(res).
		SetDatabase(db)

	assert.Equal(t, req, fields.Request)
	assert.Equal(t, res, fields.Response)
	assert.Equal(t, db, fields.Database)

	evt := &Event{Message: "test"}
	err := errors.New("test error")

	fields.SetEvents(evt).
		SetErrors(err)

	assert.Len(t, fields.Events, 1)
	assert.Equal(t, evt, fields.Events[0])
	assert.Len(t, fields.Errors, 1)
	assert.Equal(t, err, fields.Errors[0])
}

func TestLogFields_ToLogrusFields(t *testing.T) {
	req := &RequestLogField{Method: "GET", Path: "/"}
	fields := NewLogFields("2026-05-01T12:00:00Z").
		SetService("test-service").
		SetDurationMs(100).
		SetRequest(req)

	logrusFields := fields.ToLogrusFields()

	assert.Equal(t, "test-service", logrusFields["service"])
	assert.Equal(t, int64(100), logrusFields["duration_ms"])
	assert.Equal(t, req, logrusFields["request"])
	assert.Nil(t, logrusFields["response"])
	assert.Nil(t, logrusFields["database"])
}
