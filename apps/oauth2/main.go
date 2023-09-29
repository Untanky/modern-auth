package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"log/slog"
	"time"
)

const (
	CONTEXT_PATH         = "/api/v1"
	CACHE_CONTROL_HEADER = "cache-control"
	REQUEST_ID_HEADER    = "request-id"
)

var (
	requestsInstrument     metric.Int64Counter
	errorsInstrument       metric.Int64Counter
	latencyInstrument      metric.Int64Histogram
	responseSizeInstrument metric.Int64Histogram
)

type WriteFunc func([]byte) (int, error)

func (fn WriteFunc) Write(data []byte) (int, error) {
	return fn(data)
}

func main() {
	slog.Info("Application initialization starting")

	slog.Debug("Gin configuration starting")
	engine := configureGin()
	slog.Info("Gin configuration successful")
	route := engine.Group(CONTEXT_PATH)

	slog.Debug("Telemetry configuration starting")
	err := configureTelemetry(engine, route)
	if err != nil {
		slog.Error("Telemetry configuration failed", "error", err)
		panic(err)
	}
	slog.Info("Telemetry configuration successful")

	slog.Debug("Routing configuration starting")
	configureRoutes(route)
	slog.Info("Routing configuration successful")

	slog.Info("Application initialization successful")
	slog.Info("Application now starting to listen")

	err = engine.Run(":8080")
	if err != nil {
		slog.Error("Application did not run", "err", err)
	}
}

func configureGin() *gin.Engine {
	gin.DefaultWriter = WriteFunc(func(data []byte) (int, error) {
		slog.Info(string(data))
		return 0, nil
	})
	engine := gin.Default()
	engine.Use(gin.Recovery())

	return engine
}

func configureTelemetry(engine *gin.Engine, route gin.IRouter) error {
	meterProvider := otel.GetMeterProvider()
	meter := meterProvider.Meter("github.com/Untanky/modern-auth/oauth2/http")
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
	responseSizeInstrument, err = meter.Int64Histogram("response_size", metric.WithUnit("µs"))
	if err != nil {
		return err
	}

	engine.Use(handleRequestAndCorrelationId, handleTelemetry)

	return nil
}

func handleRequestAndCorrelationId(c *gin.Context) {
	requestId := c.GetHeader(REQUEST_ID_HEADER)
	if requestId == "" {
		requestId = uuid.New().String()
	}

	c.Set("requestId", requestId)
	c.Header(REQUEST_ID_HEADER, requestId)

	c.Next()
}

func handleTelemetry(c *gin.Context) {
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

func configureRoutes(route gin.IRouter) {
	route.GET("/authorization", disableCaching)
	route.GET("/authorization/succeed", disableCaching)
	route.POST("/token", disableCaching)
	route.POST("/token/validate", disableCaching)
	route.GET("/client", disableCaching)
	route.GET("/client/:id", disableCaching)
	route.POST("/client", disableCaching)
	route.DELETE("/client/:id", disableCaching)
}

func disableCaching(c *gin.Context) {
	c.Header(CACHE_CONTROL_HEADER, "no-store")
	c.Next()
}
