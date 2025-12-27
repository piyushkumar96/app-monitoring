package app_monitoring

import "github.com/prometheus/client_golang/prometheus"

// NewAppMetrics creates and registers application-level Prometheus metrics.
// It initializes an ApplicationErrorsCounter gauge for tracking application errors by error code.
//
// The ApplicationErrorsCounter metric tracks the count of errors at the application level,
// allowing you to monitor error rates and identify problematic error codes.
//
// Parameters:
//   - meta: Configuration containing the namespace and metric settings.
//     Set ApplicationErrorsCounter to nil to disable error tracking.
//
// Returns an AppMetrics instance that can be used to log and query error metrics.
//
// Example:
//
//	appMetrics := monitoring.NewAppMetrics(&monitoring.AppMetricsMeta{
//	    Namespace: "myapp",
//	    ApplicationErrorsCounter: &monitoring.MetricMeta{
//	        Labels: []string{"error_code"},
//	    },
//	})
func NewAppMetrics(meta *AppMetricsMeta) *AppMetrics {
	var appErrorsCounter *prometheus.GaugeVec
	if meta.ApplicationErrorsCounter != nil {
		appErrorsCounter = GetGaugeVec(meta.Namespace, "application_errors_total", "Tracks the counts of app errors at application level", meta.ApplicationErrorsCounter.Labels)
	}
	return &AppMetrics{
		applicationErrorsCounter: appErrorsCounter,
	}
}

// LogMetrics increments the application error counter for each provided error code.
// Call this method when application errors occur to track them in Prometheus.
//
// Parameters:
//   - errCodes: Slice of error codes to increment counters for.
//
// Example:
//
//	appMetrics.LogMetrics([]string{"ERR_DB_CONNECTION", "ERR_VALIDATION"})
func (cm *AppMetrics) LogMetrics(errCodes []string) {
	if cm.applicationErrorsCounter != nil {
		for _, errCode := range errCodes {
			cm.applicationErrorsCounter.WithLabelValues(errCode).Inc()
		}
	}
}

// GetApplicationErrorsCounterMetric returns the underlying Prometheus GaugeVec
// for the application errors counter. This can be used for advanced operations
// like resetting metrics or custom queries.
//
// Returns nil if the metric was not configured during initialization.
func (cm *AppMetrics) GetApplicationErrorsCounterMetric() *prometheus.GaugeVec {
	return cm.applicationErrorsCounter
}

// DecrementAppErrorCount decrements the application error counter for a specific error code.
// Use this when an error condition has been resolved or corrected.
//
// Parameters:
//   - errCode: The error code to decrement the counter for.
//
// Example:
//
//	appMetrics.DecrementAppErrorCount("ERR_DB_CONNECTION")
func (cm *AppMetrics) DecrementAppErrorCount(errCode string) {
	cm.applicationErrorsCounter.WithLabelValues(errCode).Dec()
}
