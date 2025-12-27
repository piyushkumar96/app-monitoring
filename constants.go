// Package app_monitoring provides Prometheus metrics collection utilities for application monitoring.
// It supports various metric types including counters, gauges, histograms, and summaries
// for tracking HTTP requests, database operations, cron jobs, pub/sub messaging, and more.
package app_monitoring

// Constants for metric status labels used across all metric types.
// These are used to categorize metrics into total, success, and failure buckets.
const (
	// Total represents the total count label value for metrics.
	Total = "total"

	// Success represents the successful operation label value for metrics.
	Success = "success"

	// Failure represents the failed operation label value for metrics.
	Failure = "failure"

	// HTTPStatus2XXMaxValue is the maximum HTTP status code considered successful (inclusive).
	HTTPStatus2XXMaxValue = 299

	// HTTPStatus2XXMinValue is the minimum HTTP status code considered successful (inclusive).
	HTTPStatus2XXMinValue = 200
)
