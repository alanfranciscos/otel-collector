package logger

import (
	"context"
	"testing"

	"github.com/alanfranciscos/otel-collector/internal/pkg/telemetry/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/resource"
)

func TestNewLogProvider(t *testing.T) {
	serviceName := "test-svc"
	cfg := &config.Config{ServiceName: serviceName}
	res := resource.Default()

	provider := NewLogProvider(&serviceName, cfg, res)

	assert.NotNil(t, provider)
	assert.Equal(t, &serviceName, provider.ServiceName)
	assert.Equal(t, cfg, provider.cfg)
	assert.Equal(t, res, provider.res)
}

func TestInitLoggerProvider_Local(t *testing.T) {
	serviceName := "test-svc"
	cfg := &config.Config{
		ServiceName:  serviceName,
		IsProduction: false,
	}
	res := resource.Default()

	provider := NewLogProvider(&serviceName, cfg, res)
	loggerProvider, err := provider.InitLoggerProvider(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, loggerProvider)
}

func TestInitLoggerProvider_ProductionHTTP(t *testing.T) {
	serviceName := "test-svc"
	cfg := &config.Config{
		ServiceName:  serviceName,
		IsProduction: true,
		Protocol:     config.ProtocolHTTP,
		Endpoint:     "localhost:4318",
	}
	res := resource.Default()

	provider := NewLogProvider(&serviceName, cfg, res)
	loggerProvider, err := provider.InitLoggerProvider(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, loggerProvider)
}

func TestInitLoggerProvider_ProductionGRPC(t *testing.T) {
	serviceName := "test-svc"
	cfg := &config.Config{
		ServiceName:  serviceName,
		IsProduction: true,
		Protocol:     config.ProtocolGRPC,
		Endpoint:     "localhost:4317",
	}
	res := resource.Default()

	provider := NewLogProvider(&serviceName, cfg, res)
	loggerProvider, err := provider.InitLoggerProvider(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, loggerProvider)
}

func TestSetupGlobalLogger(t *testing.T) {
	serviceName := "test-svc"
	cfg := &config.Config{ServiceName: serviceName}
	res := resource.Default()

	provider := NewLogProvider(&serviceName, cfg, res)
	provider.SetupGlobalLogger()

	formatter := logrus.StandardLogger().Formatter
	assert.IsType(t, &CustomJSONFormatter{}, formatter)

	customFormatter := formatter.(*CustomJSONFormatter)
	assert.Equal(t, "test-svc", customFormatter.serviceName)
}
