package otel

import (
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	"time"
)

/**
 * @author: gagral.x@gmail.com
 * @time: 2024/7/13 23:04
 * @file: otel.go
 * @description:
 */

// NewTraceExporter 创建新的跟踪导出器
func NewTraceExporter() (trace.SpanExporter, error) {
	return stdouttrace.New(stdouttrace.WithPrettyPrint())
}

// NewMetricExporter 创建新的指标导出器
func NewMetricExporter() (metric.Exporter, error) {
	return stdoutmetric.New()
}

// NewMeterProvider 使用给定的指标导出器创建新的指标提供程序
func NewMeterProvider(meterExporter metric.Exporter) *metric.MeterProvider {
	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(meterExporter, metric.WithInterval(10*time.Second))),
	)

	return meterProvider
}

// NewTraceProvider 使用给定的跟踪导出器创建新的跟踪提供程序
func NewTraceProvider(traceExporter trace.SpanExporter) *trace.TracerProvider {
	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter, trace.WithBatchTimeout(time.Second)),
	)

	return traceProvider
}
