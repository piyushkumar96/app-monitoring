// Package constants provides shared constants for application monitoring.
// These constants are used across all metric implementations.
package constants

// Constants for metric status labels used across all metric implementations.
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
