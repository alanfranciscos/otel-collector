package logger

import (
	"context"

	"github.com/alanfranciscos/otel-collector/internal/pkg/telemetry/config"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
)

type Shutdowner interface {
	Shutdown(context.Context) error
}

type LogProvider struct {
	ServiceName *string
	cfg         *config.Config
	res         *resource.Resource
}

func NewLogProvider(serviceName *string, config *config.Config, res *resource.Resource) *LogProvider {
	return &LogProvider{
		ServiceName: serviceName,
		cfg:         config,
		res:         res,
	}
}

type LoggerProviderInterface interface {
	InitLoggerProvider(ctx context.Context) (*sdklog.LoggerProvider, error)
	SetupGlobalLogger()
}

func (l *LogProvider) InitLoggerProvider(ctx context.Context) (*sdklog.LoggerProvider, error) {
	var exporter sdklog.Exporter
	var err error

	// Local: Disable OTLP Logging to stdout for cleanliness.
	// Our Logrus formatter takes care of printing rigid JSON correctly anyway.
	isProduction := l.cfg.IsProduction
	if !isProduction {
		loggerProvider := sdklog.NewLoggerProvider(
			sdklog.WithResource(l.res),
		)

		return loggerProvider, nil
	}

	switch l.cfg.Protocol {
	case config.ProtocolGRPC:
		exporter, err = otlploggrpc.New(ctx)
	case config.ProtocolHTTP:
		exporter, err = otlploghttp.New(ctx)
	}

	if err != nil {
		return nil, err
	}

	processor := sdklog.NewBatchProcessor(exporter)
	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithResource(l.res),
		sdklog.WithProcessor(processor),
	)
	return loggerProvider, nil
}

func (l *LogProvider) SetupGlobalLogger() {
	logrus.SetFormatter(&CustomJSONFormatter{
		serviceName: *l.ServiceName,
	})
}
