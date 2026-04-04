package telemetry

import (
	"os"
	"testing"
)

func TestLoadConfig_Custom(t *testing.T) {
	os.Clearenv()
	os.Setenv("OTEL_SERVICE_NAME", "my-test-service")
	os.Setenv("ENVIRONMENT", "staging")
	os.Setenv("OTEL_EXPORTER_OTLP_PROTOCOL", "grpc")
	os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")

	cfg := LoadConfig()

	if cfg.ServiceName != "MY-TEST-SERVICE" {
		t.Errorf("Expected MY-TEST-SERVICE, got %v", cfg.ServiceName)
	}
	if cfg.Environment != "STAGING" {
		t.Errorf("Expected STAGING, got %v", cfg.Environment)
	}
	if cfg.IsProduction {
		t.Errorf("Expected IsProduction false, got %v", cfg.IsProduction)
	}
	if cfg.Protocol != ProtocolGRPC {
		t.Errorf("Expected protocol grpc, got %v", cfg.Protocol)
	}
}

func TestLoadConfig_WithoutOtelServiceNameEnv(t *testing.T) {
	os.Clearenv()
	os.Unsetenv("OTEL_SERVICE_NAME")
	os.Setenv("ENVIRONMENT", "staging")
	os.Setenv("OTEL_EXPORTER_OTLP_PROTOCOL", "grpc")
	os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic due to missing OTEL_SERVICE_NAME")
		}
	}()

	LoadConfig()
}

func TestLoadConfig_WithoutOtelEndpointEnv(t *testing.T) {
	os.Clearenv()
	os.Setenv("OTEL_SERVICE_NAME", "my-test-service")
	os.Setenv("ENVIRONMENT", "staging")
	os.Setenv("OTEL_EXPORTER_OTLP_PROTOCOL", "grpc")
	os.Unsetenv("OTEL_EXPORTER_OTLP_ENDPOINT")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic due to missing OTEL_EXPORTER_OTLP_ENDPOINT")
		}
	}()

	LoadConfig()
}

func TestLoadConfig_WithoutOtelProtocolEnv(t *testing.T) {
	os.Clearenv()
	os.Setenv("OTEL_SERVICE_NAME", "my-test-service")
	os.Setenv("ENVIRONMENT", "staging")
	os.Unsetenv("OTEL_EXPORTER_OTLP_PROTOCOL")
	os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic due to missing OTEL_EXPORT PROTOCOL")
		}
	}()

	LoadConfig()
}
