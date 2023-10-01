package main

//import (
//	"log"
//	"log/slog"
//	"os"
//	"time"
//
//	"github.com/gin-gonic/gin"
//	"go.opentelemetry.io/otel"
//	"go.opentelemetry.io/otel/exporters/jaeger"
//	"go.opentelemetry.io/otel/exporters/prometheus"
//	api "go.opentelemetry.io/otel/metric"
//	"go.opentelemetry.io/otel/sdk/metric"
//	"go.opentelemetry.io/otel/sdk/resource"
//	"go.opentelemetry.io/otel/sdk/trace"
//	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
//)
//
//var meter api.Meter
//
//func init() {
//	traceExporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
//	if err != nil {
//		log.Fatal(err)
//	}
//	mergedResource, err := resource.Merge(
//		resource.Default(),
//		resource.NewWithAttributes(
//			semconv.SchemaURL,
//			semconv.ServiceName("ModernAuth"),
//		),
//	)
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	tracerProvider := trace.NewTracerProvider(
//		trace.WithBatcher(traceExporter),
//		trace.WithResource(mergedResource),
//	)
//
//	otel.SetTracerProvider(tracerProvider)
//
//	slogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
//		Level: slog.LevelDebug,
//	}))
//	slog.SetDefault(slogger)
//
//	meterExporter, err := prometheus.New()
//	if err != nil {
//		log.Fatal(err)
//	}
//	meterProvider := metric.NewMeterProvider(metric.WithReader(meterExporter))
//
//	meter = meterProvider.Meter("github.com/Untanky/modern-auth")
//
//	otel.SetMeterProvider(meterProvider)
//}
//
//type requestTelemetry struct {
//	requests     api.Int64Counter
//	errors       api.Int64Counter
//	latency      api.Int64Histogram
//	responseSize api.Int64Histogram
//}
//
//func newRequestTelemetry(meter api.Meter) (*requestTelemetry, error) {
//	requestsInstrument, err := meter.Int64Counter("requests")
//	if err != nil {
//		return nil, err
//	}
//
//	errorsInstrument, err := meter.Int64Counter("errors")
//	if err != nil {
//		return nil, err
//	}
//
//	latencyInstrument, err := meter.Int64Histogram("latency", api.WithUnit("µs"))
//	if err != nil {
//		return nil, err
//	}
//
//	responseSizeInstrument, err := meter.Int64Histogram("response_size", api.WithUnit("µs"))
//	if err != nil {
//		return nil, err
//	}
//
//	return &requestTelemetry{
//		requests:     requestsInstrument,
//		errors:       errorsInstrument,
//		latency:      latencyInstrument,
//		responseSize: responseSizeInstrument,
//	}, nil
//}
//
//func (r *requestTelemetry) handleTelemetry() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		start := time.Now()
//		path := c.Request.URL.Path
//
//		c.Next()
//
//		status := c.Writer.Status()
//		size := c.Writer.Size()
//		isError := status >= 500
//		latency := time.Since(start)
//
//		// METRICS
//		// increase counters for observable metrics
//		r.requests.Add(c.Request.Context(), 1)
//		if isError {
//			r.errors.Add(c.Request.Context(), 1)
//		}
//		r.latency.Record(c.Request.Context(), latency.Microseconds())
//		r.responseSize.Record(c.Request.Context(), int64(size))
//
//		// LOGGING
//		// create meta data fields for logging
//		msg := ""
//		fields := []any{
//			slog.String("method", c.Request.Method),
//			slog.String("path", path),
//			slog.String("ip", c.ClientIP()),
//			slog.Int("status", status),
//			slog.String("user-agent", c.Request.UserAgent()),
//			slog.Duration("latency", latency),
//			slog.Int("body-size", size),
//			slog.String("request-id", c.GetString("requestId")),
//		}
//
//		if isError {
//			slog.Error(msg, fields...)
//		} else {
//			slog.Debug(msg, fields...)
//		}
//	}
//}
