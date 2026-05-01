package telemetry

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/alanfranciscos/otel-collector/internal/pkg/telemetry/config"
	"github.com/alanfranciscos/otel-collector/internal/pkg/telemetry/provider/logger"
	"github.com/alanfranciscos/otel-collector/internal/pkg/telemetry/provider/meter"
	"github.com/alanfranciscos/otel-collector/internal/pkg/telemetry/provider/tracer"
	"github.com/stretchr/testify/assert"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestNewTelemetry(t *testing.T) {
	serviceName := "test-app"
	tel := NewTelemetry(&serviceName)

	assert.NotNil(t, tel)
	assert.Equal(t, &serviceName, tel.serviceName)
}

type MockLogProvider struct {
	logger.LoggerProviderInterface
}

func (m *MockLogProvider) InitLoggerProvider(ctx context.Context) (*sdklog.LoggerProvider, error) {
	return sdklog.NewLoggerProvider(), nil
}
func (m *MockLogProvider) SetupGlobalLogger() {}

type MockTracerProvider struct {
	tracer.TracerProviderInterface
}

func (m *MockTracerProvider) InitTracerProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	return sdktrace.NewTracerProvider(), nil
}

type MockMeterProvider struct {
	meter.MeterProviderInterface
}

func (m *MockMeterProvider) InitMeterProvider(ctx context.Context) (*sdkmetric.MeterProvider, error) {
	return sdkmetric.NewMeterProvider(), nil
}

func TestTelemetry_Initialize_WithMocks(t *testing.T) {
	os.Clearenv()
	os.Setenv("ENVIRONMENT", "local")
	os.Setenv("OTEL_SERVICE_NAME", "test")
	os.Setenv("OTEL_EXPORTER_OTLP_PROTOCOL", "http")

	serviceName := "test-app"
	tel := NewTelemetry(&serviceName)

	// Inject Mocks
	tel.loggerProvider = &MockLogProvider{}
	tel.tracerProvider = &MockTracerProvider{}
	tel.meterProvider = &MockMeterProvider{}

	ctx := context.Background()
	shutdown, err := tel.Initialize(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, shutdown)

	err = shutdown(ctx)
	assert.NoError(t, err)
}

type MockShutdowner struct {
	shouldErr bool
}

func (m *MockShutdowner) Shutdown(ctx context.Context) error {
	if m.shouldErr {
		return errors.New("mock shutdown error")
	}
	return nil
}

func TestTelemetry_BuildShutdownFunc(t *testing.T) {
	tel := NewTelemetry(nil)

	// Test successful shutdown
	successMock := &MockShutdowner{shouldErr: false}
	shutdownSuccess := tel.buildShutdownFunc(successMock, successMock, successMock)
	err := shutdownSuccess(context.Background())
	assert.NoError(t, err)

	// Test error accumulation
	errorMock := &MockShutdowner{shouldErr: true}
	shutdownError := tel.buildShutdownFunc(errorMock, errorMock, errorMock)
	err = shutdownError(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "mock shutdown error")
}

func TestTelemetry_BuildResource(t *testing.T) {
	tel := NewTelemetry(nil)
	cfg := &config.Config{
		ServiceName: "test",
		Environment: "TEST",
	}

	res, err := tel.buildResource(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}
