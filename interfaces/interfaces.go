// Package interfaces provides generic interfaces for application monitoring.
// These interfaces can be implemented by different backends (Prometheus, OpenTelemetry, etc.).
package interfaces

import (
	"time"

	"github.com/gin-gonic/gin"
	ae "github.com/piyushkumar96/app-error"
	"github.com/piyushkumar96/app-monitoring/models"
	pubsub "github.com/piyushkumar96/generic-pubsub"
)

// RouterMetricsInterface defines the contract for router-level HTTP metrics.
// Implement this interface to provide custom router metrics implementations
// for different backends (Prometheus, OpenTelemetry, StatsD, etc.).
type RouterMetricsInterface interface {
	// LogMetrics returns a Gin middleware that logs HTTP request metrics.
	LogMetrics(metricsPath string) gin.HandlerFunc
}

// DBMetricsInterface defines the contract for database operation metrics.
// Implement this interface to provide custom database metrics implementations
// for different backends (Prometheus, OpenTelemetry, StatsD, etc.).
type DBMetricsInterface interface {
	// LogMetricsPre should be called before a database operation.
	// Returns the start time for latency calculation.
	LogMetricsPre(dbMetricsLabelValues *models.DBMetricsLabelValues) time.Time

	// LogMetricsPost should be called after a database operation completes.
	LogMetricsPost(appErr *ae.AppError, dbMetricsLabelValues *models.DBMetricsLabelValues, opsExecTime time.Time)
}

// DownstreamServiceMetricsInterface defines the contract for downstream HTTP service metrics.
// Implement this interface to provide custom downstream metrics implementations
// for different backends (Prometheus, OpenTelemetry, StatsD, etc.).
type DownstreamServiceMetricsInterface interface {
	// LogMetricsPre should be called before making a downstream HTTP call.
	LogMetricsPre(dssMetricsLabelValues *models.DownstreamServiceMetricsLabelValues)

	// LogMetricsPost should be called after a downstream HTTP call completes.
	LogMetricsPost(success bool, dssMetricsLabelValues *models.DownstreamServiceMetricsLabelValues, httpMetrics *models.HTTPMetrics)
}

// CronJobMetricsInterface defines the contract for cron job execution metrics.
// Implement this interface to provide custom cron job metrics implementations
// for different backends (Prometheus, OpenTelemetry, StatsD, etc.).
type CronJobMetricsInterface interface {
	// LogMetricsPre should be called at the start of a cron job execution.
	// Returns the start time for latency calculation.
	LogMetricsPre(cjMetricsLabelValues *models.CronJobMetricsLabelValues) time.Time

	// LogMetricsPost should be called after a cron job execution completes.
	LogMetricsPost(appErr *ae.AppError, cjMetricsLabelValues *models.CronJobMetricsLabelValues, opsExecTime time.Time)
}

// PSMetricsInterface defines the contract for pub/sub messaging metrics.
// Implement this interface to provide custom pub/sub metrics implementations
// for different backends (Prometheus, OpenTelemetry, StatsD, etc.).
type PSMetricsInterface interface {
	// LogMetricsPre should be called before publishing a message or when starting to process a consumed message.
	// Returns the start time for latency calculation.
	LogMetricsPre(psMetricsLabelValues *models.PSMetricsLabelValues) time.Time

	// LogMetricsPost should be called after a pub/sub operation completes.
	LogMetricsPost(psMetricsLabelValues *models.PSMetricsLabelValues, eventTxnData *pubsub.EventTxnData)
}

// AppMetricsInterface defines the contract for application-level error metrics.
// Implement this interface to provide custom app metrics implementations
// for different backends (Prometheus, OpenTelemetry, StatsD, etc.).
type AppMetricsInterface interface {
	// LogMetrics increments the application error counter for each provided error code.
	LogMetrics(errCodes []string)

	// DecrementAppErrorCount decrements the application error counter for a specific error code.
	DecrementAppErrorCount(errCode string)
}
