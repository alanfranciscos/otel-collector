package logger

import (
	"encoding/json"
	"time"

	"github.com/alanfranciscos/otel-collector/internal/pkg/telemetry/schema/logger"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type CustomJSONFormatter struct {
	serviceName string
}

func (f *CustomJSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	span := trace.SpanFromContext(entry.Context)
	span_id := span.SpanContext().SpanID().String()
	trace_id := span.SpanContext().TraceID().String()

	payload := logger.NewLogFields(
		time.Now().Format(time.RFC3339),
	)

	payload.
		SetService(f.serviceName).
		SetTraceID(trace_id).
		SetSpanID(span_id)

	// Build a map for the final JSON to allow merging with entry.Data efficiently
	data := entry.Data
	if data == nil {
		data = make(logrus.Fields)
	}

	data["timestamp"] = payload.Timestamp
	data["service_name"] = payload.Service
	data["trace_id"] = payload.TraceID
	data["span_id"] = payload.SpanID
	data["level"] = entry.Level.String()
	data["message"] = entry.Message

	if payload.Request != nil {
		data["request"] = payload.Request
	}
	if payload.Response != nil {
		data["response"] = payload.Response
	}
	if payload.Database != nil {
		data["database"] = payload.Database
	}
	if len(payload.Errors) > 0 {
		data["errors"] = payload.Errors
	}
	if len(payload.Events) > 0 {
		data["events"] = payload.Events
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return append(bytes, '\n'), nil
}
