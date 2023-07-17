package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

type Person struct {
	Name string
}

func (p Person) String() string {
	return p.Name
}

func main() {
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatal(err)
	}
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter("github.com/Untanky/modern-auth")

	requestMetrics, err := newRequestTelemetry(meter)
	if err != nil {
		log.Fatal(err)
	}

	app := App{}
	app.Start()
	defer app.Stop()

	app.engine.Use(requestMetrics.handleTelemetry())
	app.engine.GET("/metrics", gin.WrapH(promhttp.Handler()))
	app.engine.Run(":3000")
}
