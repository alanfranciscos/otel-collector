package meter

import (
	"context"
	"testing"

	"github.com/alanfranciscos/otel-collector/internal/pkg/telemetry/config"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/resource"
)

func TestNewMeterProvider(t *testing.T) {
	cfg := &config.Config{ServiceName: "test-svc"}
	res := resource.Default()

	provider := NewMeterProvider(res, cfg)

	assert.NotNil(t, provider)
	assert.Equal(t, cfg, provider.cfg)
	assert.Equal(t, res, provider.res)
}

func TestInitMeterProvider_Local(t *testing.T) {
	cfg := &config.Config{
		ServiceName:  "test-svc",
		IsProduction: false,
	}
	res := resource.Default()

	provider := NewMeterProvider(res, cfg)
	meterProvider, err := provider.InitMeterProvider(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, meterProvider)
}

func TestInitMeterProvider_ProductionHTTP(t *testing.T) {
	cfg := &config.Config{
		ServiceName:  "test-svc",
		IsProduction: true,
		Protocol:     config.ProtocolHTTP,
		Endpoint:     "localhost:4318",
	}
	res := resource.Default()

	provider := NewMeterProvider(res, cfg)
	meterProvider, err := provider.InitMeterProvider(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, meterProvider)
}

func TestInitMeterProvider_ProductionGRPC(t *testing.T) {
	cfg := &config.Config{
		ServiceName:  "test-svc",
		IsProduction: true,
		Protocol:     config.ProtocolGRPC,
		Endpoint:     "localhost:4317",
	}
	res := resource.Default()

	provider := NewMeterProvider(res, cfg)
	meterProvider, err := provider.InitMeterProvider(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, meterProvider)
}
