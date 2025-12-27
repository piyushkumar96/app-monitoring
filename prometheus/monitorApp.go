package prometheus

import (
	"github.com/piyushkumar96/app-monitoring/interfaces"
	"github.com/piyushkumar96/app-monitoring/models"

	"github.com/prometheus/client_golang/prometheus"
)

// NewPromAppMetrics creates and registers Prometheus application-level metrics.
// It initializes an ApplicationErrorsCounter gauge for tracking application errors by error code.
//
// The ApplicationErrorsCounter metric tracks the count of errors at the application level,
// allowing you to monitor error rates and identify problematic error codes.
//
// Parameters:
//   - meta: Configuration containing the namespace and metric settings.
//     Set ApplicationErrorsCounter to nil to disable error tracking.
//
// Returns an interfaces.AppMetricsInterface instance that can be used to log and query error metrics.
func NewPromAppMetrics(meta *models.AppMetricsMeta) interfaces.AppMetricsInterface {
	var appErrorsCounter *prometheus.GaugeVec
	if meta.ApplicationErrorsCounter != nil {
		appErrorsCounter = GetPromGaugeVec(meta.Namespace, "application_errors_total", "Tracks the counts of app errors at application level", meta.ApplicationErrorsCounter.Labels)
	}
	return &PromAppMetrics{
		applicationErrorsCounter: appErrorsCounter,
	}
}

// LogMetrics increments the application error counter for each provided error code.
// Call this method when application errors occur to track them in Prometheus.
func (cm *PromAppMetrics) LogMetrics(errCodes []string) {
	if cm.applicationErrorsCounter != nil {
		for _, errCode := range errCodes {
			cm.applicationErrorsCounter.WithLabelValues(errCode).Inc()
		}
	}
}

// GetApplicationErrorsCounterMetric returns the underlying Prometheus GaugeVec
// for the application errors counter. This can be used for advanced operations
// like resetting metrics or custom queries.
func (cm *PromAppMetrics) GetApplicationErrorsCounterMetric() *prometheus.GaugeVec {
	return cm.applicationErrorsCounter
}

// DecrementAppErrorCount decrements the application error counter for a specific error code.
// Use this when an error condition has been resolved or corrected.
func (cm *PromAppMetrics) DecrementAppErrorCount(errCode string) {
	cm.applicationErrorsCounter.WithLabelValues(errCode).Dec()
}
