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

type Config struct {
	ServiceName  string
	Environment  EnvironmentEnum
	IsProduction bool
	Protocol     ExporterProtocol
	Endpoint     string
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
		return EnvLocal
	}
}

func LoadConfig() *Config {
	env := getEnvironment()
	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	protocolStr := strings.ToLower(os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL"))

	protocol := ProtocolHTTP
	if protocolStr == "grpc" {
		protocol = ProtocolGRPC
	}

	return &Config{
		ServiceName:  strings.ToUpper(serviceName),
		Environment:  env,
		IsProduction: env == EnvProduction,
		Protocol:     protocol,
		Endpoint:     strings.ToUpper(endpoint),
	}
}
