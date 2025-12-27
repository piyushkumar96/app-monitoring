package app_monitoring

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// HTTPMetrics holds HTTP request/response metrics data captured during an HTTP call.
// It is used to record metrics for downstream service calls and router-level monitoring.
type HTTPMetrics struct {
	// Method is the HTTP method used (GET, POST, PUT, DELETE, etc.).
	Method string

	// URL is the request URL path.
	URL string

	// Code is the HTTP response status code.
	Code int

	// RequestBodySizeBytes is the size of the HTTP request body in bytes.
	RequestBodySizeBytes int64

	// ResponseBodySizeBytes is the size of the HTTP response body in bytes.
	ResponseBodySizeBytes int64

	// ResponseTime is the duration taken to complete the HTTP request.
	ResponseTime time.Duration
}

// MetricMeta contains common metadata for configuring Prometheus metrics.
// It defines the labels and histogram buckets for a metric.
type MetricMeta struct {
	// Labels are the label names used for the metric.
	Labels []string

	// Buckets are the histogram bucket boundaries (only used for histogram metrics).
	Buckets []float64
}

// RouterMetricsMeta contains configuration for router-level HTTP metrics.
// Use this to configure which metrics to collect at the HTTP router/endpoint level.
type RouterMetricsMeta struct {
	// Namespace is the Prometheus namespace prefix for all router metrics.
	Namespace string

	// HTTPRequests configures the HTTP request counter metric.
	// Set to nil to disable this metric.
	HTTPRequests *MetricMeta

	// HTTPRequestsLatencyMillis configures the HTTP request latency histogram.
	// Set to nil to disable this metric.
	HTTPRequestsLatencyMillis *MetricMeta

	// HTTPRequestSizeBytes configures the HTTP request size histogram.
	// Set to nil to disable this metric.
	HTTPRequestSizeBytes *MetricMeta

	// HTTPResponseSizeBytes configures the HTTP response size histogram.
	// Set to nil to disable this metric.
	HTTPResponseSizeBytes *MetricMeta
}

// RouterMetrics holds the registered Prometheus metrics for router-level monitoring.
// It is created by NewRouterLevelMetrics and used to log HTTP endpoint metrics.
type RouterMetrics struct {
	httpRequests              *prometheus.CounterVec
	httpRequestsLatencyMillis *prometheus.HistogramVec
	httpRequestSizeBytes      *prometheus.HistogramVec
	httpResponseSizeBytes     *prometheus.HistogramVec
}

// AppMetricsMeta contains configuration for application-level error metrics.
// Use this to track application errors by error code.
type AppMetricsMeta struct {
	// Namespace is the Prometheus namespace prefix for all app metrics.
	Namespace string

	// ApplicationErrorsCounter configures the application errors gauge metric.
	// Set to nil to disable this metric.
	ApplicationErrorsCounter *MetricMeta
}

// AppMetrics holds the registered Prometheus metrics for application-level monitoring.
// It is created by NewAppMetrics and used to track application errors.
type AppMetrics struct {
	applicationErrorsCounter *prometheus.GaugeVec
}

// DownstreamServiceMetricsMeta contains configuration for downstream service HTTP metrics.
// Use this to track HTTP calls made to external/downstream services.
type DownstreamServiceMetricsMeta struct {
	// Namespace is the Prometheus namespace prefix for all downstream service metrics.
	Namespace string

	// HTTPRequests configures the HTTP request counter metric for downstream calls.
	// Set to nil to disable this metric.
	HTTPRequests *MetricMeta

	// HTTPRequestsLatencyMillis configures the HTTP request latency histogram for downstream calls.
	// Set to nil to disable this metric.
	HTTPRequestsLatencyMillis *MetricMeta

	// HTTPRequestSizeBytes configures the HTTP request size histogram for downstream calls.
	// Set to nil to disable this metric.
	HTTPRequestSizeBytes *MetricMeta

	// HTTPResponseSizeBytes configures the HTTP response size histogram for downstream calls.
	// Set to nil to disable this metric.
	HTTPResponseSizeBytes *MetricMeta
}

// DownstreamServiceMetrics holds the registered Prometheus metrics for downstream service monitoring.
// It is created by NewDownstreamServiceMetrics and used to log downstream HTTP call metrics.
type DownstreamServiceMetrics struct {
	httpRequests              *prometheus.CounterVec
	httpRequestsLatencyMillis *prometheus.HistogramVec
	httpRequestSizeBytes      *prometheus.HistogramVec
	httpResponseSizeBytes     *prometheus.HistogramVec
}

// DownstreamServiceMetricsLabelValues holds the label values for downstream service metrics.
// These values are used when logging metrics for downstream HTTP calls.
type DownstreamServiceMetricsLabelValues struct {
	// Name is the name/identifier of the downstream service being called.
	Name string

	// HTTPMethod is the HTTP method used for the downstream call.
	HTTPMethod string

	// APIIdentifier is a unique identifier for the API endpoint being called.
	APIIdentifier string
}

// DBMetricsMeta contains configuration for database operation metrics.
// Use this to track database operations (queries, inserts, updates, deletes).
type DBMetricsMeta struct {
	// Namespace is the Prometheus namespace prefix for all database metrics.
	Namespace string

	// OperationsTotal configures the database operations counter metric.
	// Set to nil to disable this metric.
	OperationsTotal *MetricMeta

	// OperationsLatencyMillis configures the database operation latency histogram.
	// Set to nil to disable this metric.
	OperationsLatencyMillis *MetricMeta
}

// DBMetrics holds the registered Prometheus metrics for database monitoring.
// It is created by NewDatabaseMetrics and used to log database operation metrics.
type DBMetrics struct {
	operationsTotal         *prometheus.CounterVec
	operationsLatencyMillis *prometheus.HistogramVec
}

// DBMetricsLabelValues holds the label values for database metrics.
// These values are used when logging metrics for database operations.
type DBMetricsLabelValues struct {
	// OpType is the type of database operation (e.g., "select", "insert", "update", "delete").
	OpType string

	// Source is the source/caller of the database operation.
	Source string

	// AdEntity is the entity/table being operated on.
	AdEntity string

	// IsTxn indicates whether the operation is part of a transaction ("true" or "false").
	IsTxn string
}

// PSMetricsMeta contains configuration for pub/sub messaging metrics.
// Use this to track message publishing and consumption operations.
type PSMetricsMeta struct {
	// Namespace is the Prometheus namespace prefix for all pub/sub metrics.
	Namespace string

	// TotalMessagesConsumed configures the message consumption counter metric.
	// Set to nil to disable this metric.
	TotalMessagesConsumed *MetricMeta

	// TotalMessagesPublished configures the message publishing counter metric.
	// Set to nil to disable this metric.
	TotalMessagesPublished *MetricMeta

	// MessagesPublishedLatencyMillis configures the message publishing latency histogram.
	// Set to nil to disable this metric.
	MessagesPublishedLatencyMillis *MetricMeta

	// MessagesPublishedSizeBytes configures the published message size histogram.
	// Set to nil to disable this metric.
	MessagesPublishedSizeBytes *MetricMeta
}

// PSMetrics holds the registered Prometheus metrics for pub/sub monitoring.
// It is created by NewPubSubMetrics and used to log messaging metrics.
type PSMetrics struct {
	totalMessagesConsumed          *prometheus.CounterVec
	totalMessagesPublished         *prometheus.CounterVec
	messagesPublishedLatencyMillis *prometheus.HistogramVec
	messagesPublishedSizeBytes     *prometheus.HistogramVec
}

// PSMetricsLabelValues holds the label values for pub/sub metrics.
// These values are used when logging metrics for messaging operations.
type PSMetricsLabelValues struct {
	// Source is the source of the message (e.g., subscription name).
	Source string

	// Entity is the entity type the message relates to.
	Entity string

	// EntityOpType is the operation type for the entity (e.g., "create", "update", "delete").
	EntityOpType string

	// ErrorCode is the error code if the operation failed (empty string for success).
	ErrorCode string
}

// CronJobMetricsMeta contains configuration for cron job execution metrics.
// Use this to track cron job executions and their latencies.
type CronJobMetricsMeta struct {
	// Namespace is the Prometheus namespace prefix for all cron job metrics.
	Namespace string

	// JobExecutionTotal configures the job execution counter metric.
	// Set to nil to disable this metric.
	JobExecutionTotal *MetricMeta

	// JobExecutionLatencyMillis configures the job execution latency histogram.
	// Set to nil to disable this metric.
	JobExecutionLatencyMillis *MetricMeta
}

// CronJobMetrics holds the registered Prometheus metrics for cron job monitoring.
// It is created by NewCronJobMetrics and used to log cron job execution metrics.
type CronJobMetrics struct {
	jobExecutionTotal         *prometheus.CounterVec
	jobExecutionLatencyMillis *prometheus.HistogramVec
}

// CronJobMetricsLabelValues holds the label values for cron job metrics.
// These values are used when logging metrics for cron job executions.
type CronJobMetricsLabelValues struct {
	// JobName is the unique name/identifier of the cron job.
	JobName string
}

// AdsAlertingMetricsMeta contains configuration for ads alerting metrics.
// Use this to track alerts generated by the ads alerting system.
type AdsAlertingMetricsMeta struct {
	// Namespace is the Prometheus namespace prefix for all alerting metrics.
	Namespace string

	// Alerts configures the alerts counter metric.
	// Set to nil to disable this metric.
	Alerts *MetricMeta
}

// AdsAlertingMetrics holds the registered Prometheus metrics for ads alerting monitoring.
type AdsAlertingMetrics struct {
	alerts *prometheus.CounterVec
}

// AdsAlertingMetricsLabelValues holds the label values for ads alerting metrics.
// These values are used when logging metrics for generated alerts.
type AdsAlertingMetricsLabelValues struct {
	// AccountID is the account identifier associated with the alert.
	AccountID string

	// MetricName is the name of the metric that triggered the alert.
	MetricName string

	// AlertLevel is the level/tier of the alert.
	AlertLevel string

	// AlertType is the type/category of the alert.
	AlertType string

	// Frequency is how often the alert is evaluated.
	Frequency string

	// Severity is the severity level of the alert (e.g., "warning", "critical").
	Severity string

	// EntityType is the type of entity the alert relates to.
	EntityType string
}

// AdsAlertingMetricsLogInfo contains additional information logged with alerts.
type AdsAlertingMetricsLogInfo struct {
	// NumberOfAlertsGenerated is the count of alerts generated.
	NumberOfAlertsGenerated float64
}
