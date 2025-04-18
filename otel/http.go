package otel

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

/**
 * @author: gagral.x@gmail.com
 * @time: 2024/7/13 23:21
 * @file: http.go
 * @description: traceInfo http
 */

// TraceHttp 配置 Trace HTTP 服务器
func TraceHttp(ctx context.Context, r *gin.Engine, tracerName string) error {

	// 初始化追踪导出器
	consoleTraceExporter, err := NewTraceExporter()
	if err != nil {
		return fmt.Errorf("failed to get console exporter (traceInfo): %v", err)
	}

	// 初始化指标导出器
	consoleMetricExporter, err := NewMetricExporter()
	if err != nil {
		return fmt.Errorf("failed to get console exporter (metric): %v", err)
	}

	// 设置追踪提供者
	tracerProvider := NewTraceProvider(consoleTraceExporter)
	defer func(tracerProvider *trace.TracerProvider, ctx context.Context) {
		err := tracerProvider.Shutdown(ctx)
		if err != nil {
		}
	}(tracerProvider, ctx)
	otel.SetTracerProvider(tracerProvider)

	// 设置指标提供者
	meterProvider := NewMeterProvider(consoleMetricExporter)
	defer func(meterProvider *metric.MeterProvider, ctx context.Context) {
		err := meterProvider.Shutdown(ctx)
		if err != nil {
		}
	}(meterProvider, ctx)
	otel.SetMeterProvider(meterProvider)

	// 设置传播器
	prop := NewPropagator()
	otel.SetTextMapPropagator(prop)

	// 配置 HTTP 服务器
	r.GET("/traceInfo", func(c *gin.Context) {
		traceInfo(c, tracerName)
	})

	return nil
}
