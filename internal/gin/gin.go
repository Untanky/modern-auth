package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"time"
)

const (
	RequestIdHeader = "request-id"
)

var (
	engine                 *gin.Engine
	requestsInstrument     metric.Int64Counter
	errorsInstrument       metric.Int64Counter
	latencyInstrument      metric.Int64Histogram
	responseSizeInstrument metric.Int64Histogram
)

func GetRouter(relativePath string) gin.IRouter {
	return engine.Group(relativePath)
}

func Start() error {
	return engine.Run()
}

type writeFunc func([]byte) (int, error)

func (fn writeFunc) Write(data []byte) (int, error) {
	return fn(data)
}

func ConfigureGin() error {
	gin.DefaultWriter = writeFunc(func(data []byte) (int, error) {
		slog.Debug(string(data))
		return 0, nil
	})
	engine = gin.Default()
	engine.Use(gin.Recovery())

	return nil
}

func ConfigureTelemetry() error {
	// TODO: Fix path
	meter := otel.GetMeterProvider().Meter("github.com/Untanky/modern-auth/http")

	var err error
	requestsInstrument, err = meter.Int64Counter("requestsInstrument")
	if err != nil {
		return err
	}
	errorsInstrument, err = meter.Int64Counter("errorsInstrument")
	if err != nil {
		return err
	}
	latencyInstrument, err = meter.Int64Histogram("latencyInstrument", metric.WithUnit("µs"))
	if err != nil {
		return err
	}
	responseSizeInstrument, err = meter.Int64Histogram("response_size", metric.WithUnit("byte"))
	if err != nil {
		return err
	}

	engine.Use(handleRequestId, handleRequestTelemetry, handleError)

	return nil
}

func handleRequestId(c *gin.Context) {
	requestId := c.GetHeader(RequestIdHeader)
	if requestId == "" {
		requestId = uuid.New().String()
	}

	c.Set("requestId", requestId)
	c.Header(RequestIdHeader, requestId)

	c.Next()
}

func handleRequestTelemetry(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path

	c.Next()

	status := c.Writer.Status()
	size := c.Writer.Size()
	isError := status >= 500
	latency := time.Since(start)

	// METRICS
	// increase counters for observable metrics
	requestsInstrument.Add(c.Request.Context(), 1)
	if isError {
		errorsInstrument.Add(c.Request.Context(), 1)
	}
	latencyInstrument.Record(c.Request.Context(), latency.Microseconds())
	responseSizeInstrument.Record(c.Request.Context(), int64(size))

	// LOGGING
	// create meta data fields for logging
	msg := ""
	fields := []any{
		slog.String("method", c.Request.Method),
		slog.String("path", path),
		slog.String("ip", c.ClientIP()),
		slog.Int("status", status),
		slog.String("user-agent", c.Request.UserAgent()),
		slog.Duration("latencyInstrument", latency),
		slog.Int("body-size", size),
		slog.String("request-id", c.GetString("requestId")),
	}

	if isError {
		slog.Error(msg, fields...)
	} else {
		slog.Debug(msg, fields...)
	}
}

func handleError(ctx *gin.Context) {
	ctx.Next()

	for _, err := range ctx.Errors {
		trace.SpanFromContext(ctx).RecordError(err)
		slog.ErrorContext(ctx, "Encountered error", "err", err)
	}
}
