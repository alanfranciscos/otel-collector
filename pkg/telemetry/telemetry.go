package telemetry

import (
	"context"
	"errors"

	"github.com/alanfranciscos/otel-collector/internal/pkg/telemetry/config"
	"github.com/alanfranciscos/otel-collector/internal/pkg/telemetry/provider/logger"
	"github.com/alanfranciscos/otel-collector/internal/pkg/telemetry/provider/meter"
	"github.com/alanfranciscos/otel-collector/internal/pkg/telemetry/provider/tracer"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// Telemetry manages the initialization of OpenTelemetry signals.
type Telemetry struct {
	serviceName    *string
	loggerProvider logger.LoggerProviderInterface
	tracerProvider tracer.TracerProviderInterface
	meterProvider  meter.MeterProviderInterface
}

// NewTelemetry creates a new Telemetry instance.
func NewTelemetry(serviceName *string) *Telemetry {
	return &Telemetry{
		serviceName: serviceName,
	}
}

// Initialize configures Logs, Traces, and Metrics.
func (t *Telemetry) Initialize(ctx context.Context) (func(context.Context) error, error) {
	cfg := config.LoadConfig()

	res, err := t.buildResource(cfg)
	if err != nil {
		return nil, err
	}

	t.setupPropagation()

	// Use existing providers if injected (for tests), otherwise create defaults
	if t.loggerProvider == nil {
		t.loggerProvider = logger.NewLogProvider(t.serviceName, cfg, res)
	}
	if t.tracerProvider == nil {
		t.tracerProvider = tracer.NewTracerProvider(res, cfg)
	}
	if t.meterProvider == nil {
		t.meterProvider = meter.NewMeterProvider(res, cfg)
	}

	// Initialize signals
	lp, err := t.loggerProvider.InitLoggerProvider(ctx)
	if err != nil {
		return nil, err
	}
	t.loggerProvider.SetupGlobalLogger()

	tp, err := t.tracerProvider.InitTracerProvider(ctx)
	if err != nil {
		return nil, err
	}
	otel.SetTracerProvider(tp)

	mp, err := t.meterProvider.InitMeterProvider(ctx)
	if err != nil {
		return nil, err
	}
	otel.SetMeterProvider(mp)

	logrus.Info("Telemetry successfully initialized")

	return t.buildShutdownFunc(lp, tp, mp), nil
}

func (t *Telemetry) buildResource(cfg *config.Config) (*resource.Resource, error) {
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			resource.Default().SchemaURL(),
			semconv.ServiceName(cfg.ServiceName),
			semconv.DeploymentEnvironment(string(cfg.Environment)),
			semconv.ServiceVersion("1.0.0"),
		),
	)
}

func (t *Telemetry) setupPropagation() {
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
}

func (t *Telemetry) buildShutdownFunc(lp, tp, mp Shutdowner) func(context.Context) error {
	return func(ctx context.Context) error {
		var errs []error
		if err := mp.Shutdown(ctx); err != nil {
			errs = append(errs, err)
		}
		if err := tp.Shutdown(ctx); err != nil {
			errs = append(errs, err)
		}
		if err := lp.Shutdown(ctx); err != nil {
			errs = append(errs, err)
		}
		return errors.Join(errs...)
	}
}

// Shutdowner defines the common interface for shutting down providers.
type Shutdowner interface {
	Shutdown(context.Context) error
}
