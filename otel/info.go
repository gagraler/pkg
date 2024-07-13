package otel

import (
	"encoding/json"
	"github.com/gagraler/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"net/http"
)

/**
 * @author: gagral.x@gmail.com
 * @time: 2024/7/13 23:00
 * @file: info.go
 * @description:
 */

type InfoResp struct {
	Version     string `json:"version"`
	ServiceName string `json:"serviceName"`
}

var (
	tracer      = otel.Tracer("")
	meter       = otel.Meter("")
	viewCounter metric.Int64Counter
	log         = logger.SugaredLogger()
)

func init() {
	var err error
	viewCounter, err = meter.Int64Counter("user.views",
		metric.WithDescription("The number of views"),
		metric.WithUnit("{views}"))
	if err != nil {
		log.DPanic(err)
	}
}

// info
func info(c *gin.Context, tracerName string) {
	ctx, span := tracer.Start(c.Request.Context(), "info")
	defer span.End()

	viewCounter.Add(ctx, 1)

	c.Header("Content-Type", "application/json")
	response := InfoResp{Version: "0.1.0", ServiceName: tracerName}
	err := json.NewEncoder(c.Writer).Encode(response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode response"})
		return
	}
}

// NewPropagator
func NewPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
	)
}
