// Package models provides shared data models and configuration types for application monitoring.
// These models are used across all metric implementations.
package models

import "time"

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

// MetricMeta contains common metadata for configuring metrics.
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
	// Namespace is the metric namespace prefix for all router metrics.
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

// AppMetricsMeta contains configuration for application-level error metrics.
// Use this to track application errors by error code.
type AppMetricsMeta struct {
	// Namespace is the metric namespace prefix for all app metrics.
	Namespace string

	// ApplicationErrorsCounter configures the application errors gauge metric.
	// Set to nil to disable this metric.
	ApplicationErrorsCounter *MetricMeta
}

// DownstreamServiceMetricsMeta contains configuration for downstream service HTTP metrics.
// Use this to track HTTP calls made to external/downstream services.
type DownstreamServiceMetricsMeta struct {
	// Namespace is the metric namespace prefix for all downstream service metrics.
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
	// Namespace is the metric namespace prefix for all database metrics.
	Namespace string

	// OperationsTotal configures the database operations counter metric.
	// Set to nil to disable this metric.
	OperationsTotal *MetricMeta

	// OperationsLatencyMillis configures the database operation latency histogram.
	// Set to nil to disable this metric.
	OperationsLatencyMillis *MetricMeta
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
	// Namespace is the metric namespace prefix for all pub/sub metrics.
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
	// Namespace is the metric namespace prefix for all cron job metrics.
	Namespace string

	// JobExecutionTotal configures the job execution counter metric.
	// Set to nil to disable this metric.
	JobExecutionTotal *MetricMeta

	// JobExecutionLatencyMillis configures the job execution latency histogram.
	// Set to nil to disable this metric.
	JobExecutionLatencyMillis *MetricMeta
}

// CronJobMetricsLabelValues holds the label values for cron job metrics.
// These values are used when logging metrics for cron job executions.
type CronJobMetricsLabelValues struct {
	// JobName is the unique name/identifier of the cron job.
	JobName string
}
