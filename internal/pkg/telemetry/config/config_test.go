package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_Defaults(t *testing.T) {
	os.Clearenv()
	cfg := LoadConfig()

	assert.Equal(t, EnvLocal, cfg.Environment)
	assert.False(t, cfg.IsProduction)
	assert.Equal(t, ProtocolHTTP, cfg.Protocol)
	assert.Equal(t, "", cfg.ServiceName)
	assert.Equal(t, "", cfg.Endpoint)
}

func TestLoadConfig_Custom(t *testing.T) {
	os.Clearenv()
	os.Setenv("OTEL_SERVICE_NAME", "my-test-service")
	os.Setenv("ENVIRONMENT", "staging")
	os.Setenv("OTEL_EXPORTER_OTLP_PROTOCOL", "grpc")
	os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")

	cfg := LoadConfig()

	assert.Equal(t, "MY-TEST-SERVICE", cfg.ServiceName)
	assert.Equal(t, EnvStaging, cfg.Environment)
	assert.False(t, cfg.IsProduction)
	assert.Equal(t, ProtocolGRPC, cfg.Protocol)
	assert.Equal(t, "LOCALHOST:4317", cfg.Endpoint)
}

func TestLoadConfig_Production(t *testing.T) {
	os.Clearenv()
	os.Setenv("ENVIRONMENT", "production")

	cfg := LoadConfig()

	assert.Equal(t, EnvProduction, cfg.Environment)
	assert.True(t, cfg.IsProduction)
}

func TestLoadConfig_HttpProtocol(t *testing.T) {
	os.Clearenv()
	os.Setenv("OTEL_EXPORTER_OTLP_PROTOCOL", "http")

	cfg := LoadConfig()

	assert.Equal(t, ProtocolHTTP, cfg.Protocol)
}
