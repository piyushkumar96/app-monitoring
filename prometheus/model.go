package prometheus

import "github.com/prometheus/client_golang/prometheus"

// PromRouterMetrics holds the registered Prometheus metrics for router-level monitoring.
// It implements interfaces.RouterMetricsInterface.
type PromRouterMetrics struct {
	httpRequests              *prometheus.CounterVec
	httpRequestsLatencyMillis *prometheus.HistogramVec
	httpRequestSizeBytes      *prometheus.HistogramVec
	httpResponseSizeBytes     *prometheus.HistogramVec
}

// PromAppMetrics holds the registered Prometheus metrics for application-level monitoring.
// It implements interfaces.AppMetricsInterface.
type PromAppMetrics struct {
	applicationErrorsCounter *prometheus.GaugeVec
}

// PromDownstreamServiceMetrics holds the registered Prometheus metrics for downstream service monitoring.
// It implements interfaces.DownstreamServiceMetricsInterface.
type PromDownstreamServiceMetrics struct {
	httpRequests              *prometheus.CounterVec
	httpRequestsLatencyMillis *prometheus.HistogramVec
	httpRequestSizeBytes      *prometheus.HistogramVec
	httpResponseSizeBytes     *prometheus.HistogramVec
}

// PromDBMetrics holds the registered Prometheus metrics for database monitoring.
// It implements interfaces.DBMetricsInterface.
type PromDBMetrics struct {
	operationsTotal         *prometheus.CounterVec
	operationsLatencyMillis *prometheus.HistogramVec
}

// PromPSMetrics holds the registered Prometheus metrics for pub/sub monitoring.
// It implements interfaces.PSMetricsInterface.
type PromPSMetrics struct {
	totalMessagesConsumed          *prometheus.CounterVec
	totalMessagesPublished         *prometheus.CounterVec
	messagesPublishedLatencyMillis *prometheus.HistogramVec
	messagesPublishedSizeBytes     *prometheus.HistogramVec
}

// PromCronJobMetrics holds the registered Prometheus metrics for cron job monitoring.
// It implements interfaces.CronJobMetricsInterface.
type PromCronJobMetrics struct {
	jobExecutionTotal         *prometheus.CounterVec
	jobExecutionLatencyMillis *prometheus.HistogramVec
}
