package config

import (
	"os"
	"strings"
)

type ExporterProtocol string

const (
	ProtocolHTTP ExporterProtocol = "http"
	ProtocolGRPC ExporterProtocol = "grpc"
)

type EnvironmentEnum string

const (
	EnvProduction EnvironmentEnum = "PRODUCTION"
	EnvStaging    EnvironmentEnum = "STAGING"
	EnvLocal      EnvironmentEnum = "LOCAL"
)

type EnvConfig struct {
	Environment EnvironmentEnum
	ServiceName string
	Protocol    ExporterProtocol
	Endpoint    string
}

func getEnvironment() EnvironmentEnum {
	env := strings.ToUpper(os.Getenv("ENVIRONMENT"))
	switch env {
	case "PRODUCTION":
		return EnvProduction
	case "STAGING":
		return EnvStaging
	case "LOCAL":
		return EnvLocal
	default:
		panic("Invalid ENVIRONMENT value. Must be one of [PRODUCTION, STAGING, LOCAL]")
	}
}

type Config struct {
	ServiceName  string
	Environment  EnvironmentEnum
	IsProduction bool
	Protocol     ExporterProtocol
	Endpoint     string
}

func loadEnvConfig() EnvConfig {
	environment := getEnvironment()
	if os.Getenv("OTEL_SERVICE_NAME") == "" {
		panic("OTEL_SERVICE_NAME environment variable is required")
	}
	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") == "" {
		panic("OTEL_EXPORTER_OTLP_ENDPOINT environment variable is required")
	}
	if os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL") == "" {
		panic("OTEL_EXPORTER_OTLP_PROTOCOL environment variable is required [http, grpc]")
	}

	protocol := ProtocolHTTP
	if strings.ToLower(os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL")) == "grpc" {
		protocol = ProtocolGRPC
	}

	envConfig := EnvConfig{
		Environment: environment,
		ServiceName: strings.ToUpper(os.Getenv("OTEL_SERVICE_NAME")),
		Protocol:    protocol,
		Endpoint:    strings.ToUpper(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")),
	}

	return envConfig
}

func LoadConfig() Config {
	envConfig := loadEnvConfig()

	return Config{
		ServiceName:  envConfig.ServiceName,
		Environment:  envConfig.Environment,
		IsProduction: envConfig.Environment == "PRODUCTION",
		Protocol:     envConfig.Protocol,
		Endpoint:     envConfig.Endpoint,
	}
}
