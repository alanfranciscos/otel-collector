package meter

import (
	"context"

	"github.com/alanfranciscos/otel-collector/internal/pkg/telemetry/config"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

type Shutdowner interface {
	Shutdown(context.Context) error
}

type MeterProviderInterface interface {
	InitMeterProvider(ctx context.Context) (*sdkmetric.MeterProvider, error)
}

type MeterProvider struct {
	res *resource.Resource
	cfg *config.Config
}

func NewMeterProvider(res *resource.Resource, cfg *config.Config) *MeterProvider {
	return &MeterProvider{
		res: res,
		cfg: cfg,
	}
}

func (m *MeterProvider) InitMeterProvider(ctx context.Context) (*sdkmetric.MeterProvider, error) {
	var exporter sdkmetric.Exporter
	var err error

	if m.cfg.IsProduction {
		if m.cfg.Protocol == config.ProtocolGRPC {
			exporter, err = otlpmetricgrpc.New(ctx)
		} else {
			exporter, err = otlpmetrichttp.New(ctx)
		}

		if err != nil {
			return nil, err
		}

		mp := sdkmetric.NewMeterProvider(
			sdkmetric.WithResource(m.res),
			sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		)
		return mp, nil
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(m.res),
	)
	return mp, nil
}
