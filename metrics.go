package main

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func initMetrics() {
	log.Info().Msg("Initialize metrics")

	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(ServiceName),
		),
	)
	HandleErr(err, "failed to create resource")

	cfg := prometheus.Config{
		DefaultHistogramBoundaries: []float64{50, 100, 200, 300, 400},
	}

	c := controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries(cfg.DefaultHistogramBoundaries),
			),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
		controller.WithCollectPeriod(2*time.Second),
		controller.WithResource(res),
	)

	exporter, err := prometheus.New(cfg, c)
	if err != nil {
		log.Panic().Err(err).Msg("failed to initialize prometheus exporter")
	}
	global.SetMeterProvider(exporter.MeterProvider())
	http.HandleFunc("/metrics", exporter.ServeHTTP)
}
