package main

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/metric"
)

type requestTelemetry struct {
	requests     metric.Int64Counter
	errors       metric.Int64Counter
	latency      metric.Int64Histogram
	responseSize metric.Int64Histogram
}

func newRequestTelemetry(meter metric.Meter) (*requestTelemetry, error) {
	requestsInstrument, err := meter.Int64Counter("requests")
	if err != nil {
		return nil, err
	}

	errorsInstrument, err := meter.Int64Counter("errors")
	if err != nil {
		return nil, err
	}

	latencyInstrument, err := meter.Int64Histogram("latency", metric.WithUnit("µs"))
	if err != nil {
		return nil, err
	}

	responseSizeInstrument, err := meter.Int64Histogram("response_size", metric.WithUnit("µs"))
	if err != nil {
		return nil, err
	}

	return &requestTelemetry{
		requests:     requestsInstrument,
		errors:       errorsInstrument,
		latency:      latencyInstrument,
		responseSize: responseSizeInstrument,
	}, nil
}

func (r *requestTelemetry) handleTelemetry() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		status := c.Writer.Status()
		size := c.Writer.Size()
		isError := status >= 500
		latency := time.Since(start)

		// METRICS
		// increase counters for observable metrics
		r.requests.Add(c.Request.Context(), 1)
		if isError {
			r.errors.Add(c.Request.Context(), 1)
		}
		r.latency.Record(c.Request.Context(), latency.Microseconds())
		r.responseSize.Record(c.Request.Context(), int64(size))

		// LOGGING
		// create meta data fields for logging
		msg := ""
		fields := []any{
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.String("ip", c.ClientIP()),
			slog.Int("status", status),
			slog.String("user-agent", c.Request.UserAgent()),
			slog.Duration("latency", latency),
			slog.Int("body-size", size),
			slog.String("request-id", c.GetString("requestId")),
		}

		// select logger func
		var logFunc func(msg string, fields ...any)
		if isError {
			logFunc = slog.Error
		} else {
			logFunc = slog.Debug
		}

		// actually log
		logFunc(msg, fields...)
	}
}
