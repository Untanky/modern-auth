package main

import (
	"context"
	"fmt"
	"github.com/Untanky/modern-auth/registry"
	"log/slog"
	"net/http"
	"time"
)

var client http.Client

type healthCheckTarget struct {
	endpoint         string
	interval         time.Duration
	healthyStatus    int
	healthyThreshold int

	preparedRequest *http.Request
	ctx             context.Context

	healthyMessage   string
	unhealthyMessage string
	failedMessage    string

	lastChecked  time.Time
	lastStatus   int
	healthyCount int
}

func newHealthCheckTarget(ctx context.Context, healthCheck *registry.HealthCheck) (*healthCheckTarget, error) {
	request, err := http.NewRequest("GET", healthCheck.Endpoint, nil)
	if err != nil {
		return nil, err
	}

	return &healthCheckTarget{
		endpoint:         healthCheck.Endpoint,
		interval:         time.Duration(healthCheck.Interval),
		healthyStatus:    int(healthCheck.HealthyStatus),
		healthyThreshold: int(healthCheck.HealthyThreshold),

		preparedRequest: request,
		ctx:             ctx,

		healthyMessage:   fmt.Sprintf("Checked health of endpoint '%s'. Endpoint is healthy", healthCheck.Endpoint),
		unhealthyMessage: fmt.Sprintf("Checked health of endpoint '%s'. Endpoint is unhealthy", healthCheck.Endpoint),
		failedMessage:    fmt.Sprintf("Checked health of endpoint '%s'. Endpoint has failed", healthCheck.Endpoint),

		lastChecked:  time.Unix(0, 0),
		lastStatus:   0,
		healthyCount: 0,
	}, nil
}

func (target *healthCheckTarget) CheckHealth() {
	resp, err := client.Do(target.preparedRequest)
	if err != nil || resp.StatusCode != target.healthyStatus {
		if target.healthyCount > 0 {
			target.healthyCount = 0
		}
		target.healthyCount--
	} else {
		if target.healthyCount < 0 {
			target.healthyCount = 0
		}
		target.healthyCount++
	}
	target.lastStatus = resp.StatusCode
	target.lastChecked = time.Now()
}

func (target *healthCheckTarget) LogStatus() {
	status := target.GetStatus()

	attr := []slog.Attr{
		slog.String("endpoint", target.endpoint),
		slog.String("health", status),
		slog.Int("healthyCount", target.healthyCount),
		slog.Int("httpStatus", target.lastStatus),
	}

	switch status {
	case "healthy":
		slog.InfoContext(target.ctx, target.healthyMessage, attr)
	case "unhealthy":
		slog.WarnContext(target.ctx, target.unhealthyMessage, attr)
	case "failed":
		slog.ErrorContext(target.ctx, target.failedMessage, attr)
	}
}

func (target *healthCheckTarget) GetStatus() string {
	switch {
	case target.healthyCount > 0:
		return "healthy"
	case target.healthyCount > -target.healthyThreshold:
		return "unhealthy"
	default:
		return "failed"
	}
}
