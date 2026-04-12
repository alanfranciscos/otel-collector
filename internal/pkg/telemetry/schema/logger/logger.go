package logger

import "github.com/sirupsen/logrus"

type logFields struct {
	Service    string           `json:"service_name"`
	DurationMs int64            `json:"duration_ms"`
	Request    RequestLogField  `json:"request"`
	Response   ResponseLogField `json:"response"`
	Database   DatabaseLogField `json:"database"`
	Errors     []error          `json:"errors,omitempty"`
	Events     []Event          `json:"events,omitempty"`
}

type LogFields interface {
	ToLogrusFields() logrus.Fields
	SetEvents(events ...Event)
	SetErrors(errors ...error)
}

func NewLogFields(serviceName string, durationMs int64, request RequestLogField, response ResponseLogField, database DatabaseLogField) LogFields {
	return &logFields{
		Service:    serviceName,
		DurationMs: durationMs,
		Request:    request,
		Response:   response,
		Database:   database,
	}
}

func (f *logFields) SetEvents(events ...Event) {
	f.Events = append(f.Events, events...)
}

func (f *logFields) SetErrors(errors ...error) {
	f.Errors = append(f.Errors, errors...)
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
