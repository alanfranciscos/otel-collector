package tracer

import (
	"context"

	"github.com/alanfranciscos/otel-collector/internal/pkg/telemetry/config"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Shutdowner interface {
	Shutdown(context.Context) error
}

type TracerProviderInterface interface {
	InitTracerProvider(ctx context.Context) (*sdktrace.TracerProvider, error)
}

type TracerProvider struct {
	res *resource.Resource
	cfg *config.Config
}

func NewTracerProvider(res *resource.Resource, cfg *config.Config) *TracerProvider {
	return &TracerProvider{
		res: res,
		cfg: cfg,
	}
}

func (t *TracerProvider) InitTracerProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	var exporter sdktrace.SpanExporter
	var err error

	if t.cfg.IsProduction {
		if t.cfg.Protocol == config.ProtocolGRPC {
			exporter, err = otlptracegrpc.New(ctx)
		} else {
			exporter, err = otlptracehttp.New(ctx)
		}

		if err != nil {
			return nil, err
		}

		bsp := sdktrace.NewBatchSpanProcessor(exporter)
		tp := sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithResource(t.res),
			sdktrace.WithSpanProcessor(bsp),
		)
		return tp, nil
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(t.res),
	)
	return tp, nil
}
