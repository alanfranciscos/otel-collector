package logger

import (
	"github.com/sirupsen/logrus"
)

type logFields struct {
	Timestamp  string            `json:"timestamp"`
	Service    string            `json:"service_name,omitempty"`
	TraceID    string            `json:"trace_id,omitempty"`
	SpanID     string            `json:"span_id,omitempty"`
	DurationMs int64             `json:"duration_ms,omitempty"`
	Request    *RequestLogField  `json:"request,omitempty"`
	Response   *ResponseLogField `json:"response,omitempty"`
	Database   *DatabaseLogField `json:"database,omitempty"`
	Errors     []error           `json:"errors,omitempty"`
	Events     []*Event          `json:"events,omitempty"`
}

type LogFields interface {
	ToLogrusFields() logrus.Fields
	SetService(service string) *logFields
	SetEvents(events ...*Event) *logFields
	SetErrors(errors ...error) *logFields
	SetDurationMs(durationMs int64) *logFields
	SetRequest(request *RequestLogField) *logFields
	SetResponse(response *ResponseLogField) *logFields
	SetDatabase(database *DatabaseLogField) *logFields
	SetTraceID(traceID string) *logFields
	SetSpanID(spanID string) *logFields
}

func NewLogFields(timestamp string) *logFields {
	return &logFields{
		Timestamp: timestamp,
	}
}

func (f *logFields) SetService(service string) *logFields {
	f.Service = service
	return f
}

func (f *logFields) SetTraceID(traceID string) *logFields {
	f.TraceID = traceID
	return f
}

func (f *logFields) SetSpanID(spanID string) *logFields {
	f.SpanID = spanID
	return f
}

func (f *logFields) SetEvents(events ...*Event) *logFields {
	f.Events = append(f.Events, events...)
	return f
}

func (f *logFields) SetErrors(errors ...error) *logFields {
	f.Errors = append(f.Errors, errors...)
	return f
}

func (f *logFields) SetDurationMs(durationMs int64) *logFields {
	f.DurationMs = durationMs
	return f
}

func (f *logFields) SetRequest(request *RequestLogField) *logFields {
	f.Request = request
	return f
}

func (f *logFields) SetResponse(response *ResponseLogField) *logFields {
	f.Response = response
	return f
}

func (f *logFields) SetDatabase(database *DatabaseLogField) *logFields {
	f.Database = database
	return f
}

func (f *logFields) ToLogrusFields() logrus.Fields {
	return logrus.Fields{
		"service":     f.Service,
		"duration_ms": f.DurationMs,
		"request":     f.Request,
		"response":    f.Response,
		"database":    f.Database,
		"errors":      f.Errors,
		"events":      f.Events,
	}
}
