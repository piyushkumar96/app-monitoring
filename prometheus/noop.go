package prometheus

import (
	"time"

	"github.com/piyushkumar96/app-monitoring/interfaces"
	"github.com/piyushkumar96/app-monitoring/models"

	"github.com/gin-gonic/gin"
	ae "github.com/piyushkumar96/app-error"
	pubsub "github.com/piyushkumar96/generic-pubsub"
)

// NoOpPromRouterMetrics is a no-operation implementation of RouterMetricsInterface.
// Use this for testing or when you want to disable Prometheus router metrics collection.
type NoOpPromRouterMetrics struct{}

// NewNoOpPromRouterMetrics creates a new no-op Prometheus router metrics instance.
func NewNoOpPromRouterMetrics() interfaces.RouterMetricsInterface {
	return &NoOpPromRouterMetrics{}
}

// LogMetrics returns a pass-through middleware that does nothing.
func (n *NoOpPromRouterMetrics) LogMetrics(_ string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

// NoOpPromDBMetrics is a no-operation implementation of DBMetricsInterface.
// Use this for testing or when you want to disable Prometheus database metrics collection.
type NoOpPromDBMetrics struct{}

// NewNoOpPromDBMetrics creates a new no-op Prometheus database metrics instance.
func NewNoOpPromDBMetrics() interfaces.DBMetricsInterface {
	return &NoOpPromDBMetrics{}
}

// LogMetricsPre does nothing and returns the current time.
func (n *NoOpPromDBMetrics) LogMetricsPre(_ *models.DBMetricsLabelValues) time.Time {
	return time.Now()
}

// LogMetricsPost does nothing.
func (n *NoOpPromDBMetrics) LogMetricsPost(_ *ae.AppError, _ *models.DBMetricsLabelValues, _ time.Time) {
}

// NoOpPromDownstreamServiceMetrics is a no-operation implementation of DownstreamServiceMetricsInterface.
// Use this for testing or when you want to disable Prometheus downstream service metrics collection.
type NoOpPromDownstreamServiceMetrics struct{}

// NewNoOpPromDownstreamServiceMetrics creates a new no-op Prometheus downstream service metrics instance.
func NewNoOpPromDownstreamServiceMetrics() interfaces.DownstreamServiceMetricsInterface {
	return &NoOpPromDownstreamServiceMetrics{}
}

// LogMetricsPre does nothing.
func (n *NoOpPromDownstreamServiceMetrics) LogMetricsPre(_ *models.DownstreamServiceMetricsLabelValues) {
}

// LogMetricsPost does nothing.
func (n *NoOpPromDownstreamServiceMetrics) LogMetricsPost(_ bool, _ *models.DownstreamServiceMetricsLabelValues, _ *models.HTTPMetrics) {
}

// NoOpPromCronJobMetrics is a no-operation implementation of CronJobMetricsInterface.
// Use this for testing or when you want to disable Prometheus cron job metrics collection.
type NoOpPromCronJobMetrics struct{}

// NewNoOpPromCronJobMetrics creates a new no-op Prometheus cron job metrics instance.
func NewNoOpPromCronJobMetrics() interfaces.CronJobMetricsInterface {
	return &NoOpPromCronJobMetrics{}
}

// LogMetricsPre does nothing and returns the current time.
func (n *NoOpPromCronJobMetrics) LogMetricsPre(_ *models.CronJobMetricsLabelValues) time.Time {
	return time.Now()
}

// LogMetricsPost does nothing.
func (n *NoOpPromCronJobMetrics) LogMetricsPost(_ *ae.AppError, _ *models.CronJobMetricsLabelValues, _ time.Time) {
}

// NoOpPromPSMetrics is a no-operation implementation of PSMetricsInterface.
// Use this for testing or when you want to disable Prometheus pub/sub metrics collection.
type NoOpPromPSMetrics struct{}

// NewNoOpPromPSMetrics creates a new no-op Prometheus pub/sub metrics instance.
func NewNoOpPromPSMetrics() interfaces.PSMetricsInterface {
	return &NoOpPromPSMetrics{}
}

// LogMetricsPre does nothing and returns the current time.
func (n *NoOpPromPSMetrics) LogMetricsPre(_ *models.PSMetricsLabelValues) time.Time {
	return time.Now()
}

// LogMetricsPost does nothing.
func (n *NoOpPromPSMetrics) LogMetricsPost(_ *models.PSMetricsLabelValues, _ *pubsub.EventTxnData) {
}

// NoOpPromAppMetrics is a no-operation implementation of AppMetricsInterface.
// Use this for testing or when you want to disable Prometheus application error metrics collection.
type NoOpPromAppMetrics struct{}

// NewNoOpPromAppMetrics creates a new no-op Prometheus application metrics instance.
func NewNoOpPromAppMetrics() interfaces.AppMetricsInterface {
	return &NoOpPromAppMetrics{}
}

// LogMetrics does nothing.
func (n *NoOpPromAppMetrics) LogMetrics(_ []string) {
}

// DecrementAppErrorCount does nothing.
func (n *NoOpPromAppMetrics) DecrementAppErrorCount(_ string) {
}

// Compile-time interface implementation checks for NoOp types
var (
	_ interfaces.RouterMetricsInterface            = (*NoOpPromRouterMetrics)(nil)
	_ interfaces.DBMetricsInterface                = (*NoOpPromDBMetrics)(nil)
	_ interfaces.DownstreamServiceMetricsInterface = (*NoOpPromDownstreamServiceMetrics)(nil)
	_ interfaces.CronJobMetricsInterface           = (*NoOpPromCronJobMetrics)(nil)
	_ interfaces.PSMetricsInterface                = (*NoOpPromPSMetrics)(nil)
	_ interfaces.AppMetricsInterface               = (*NoOpPromAppMetrics)(nil)
)
