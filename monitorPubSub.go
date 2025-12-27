package app_monitoring

import (
	"time"

	pubsub "github.com/piyushkumar96/generic-pubsub"
	"github.com/prometheus/client_golang/prometheus"
)

// NewPubSubMetrics creates and registers Prometheus metrics for pub/sub messaging operations.
// It initializes counters for message counts and histograms for latencies and message sizes.
//
// The metrics track:
//   - TotalMessagesConsumed: Counter for consumed messages (total/success/failure)
//   - TotalMessagesPublished: Counter for published messages (total/success/failure)
//   - MessagesPublishedLatencyMillis: Histogram for publish latency in milliseconds
//   - MessagesPublishedSizeBytes: Histogram for published message size in bytes
//
// Parameters:
//   - meta: Configuration containing the namespace and metric settings.
//     Set individual metric configs to nil to disable them.
//
// Returns a PSMetrics instance for logging pub/sub messaging metrics.
//
// Example:
//
//	psMetrics := monitoring.NewPubSubMetrics(&monitoring.PSMetricsMeta{
//	    Namespace: "myapp",
//	    TotalMessagesConsumed: &monitoring.MetricMeta{
//	        Labels: []string{"source", "entity", "op_type", "status", "error_code"},
//	    },
//	    TotalMessagesPublished: &monitoring.MetricMeta{
//	        Labels: []string{"entity", "op_type", "status"},
//	    },
//	    MessagesPublishedLatencyMillis: &monitoring.MetricMeta{
//	        Labels:  []string{"entity", "op_type"},
//	        Buckets: monitoring.GetExponentialBuckets(10, 2, 10),
//	    },
//	})
func NewPubSubMetrics(meta *PSMetricsMeta) *PSMetrics {
	var totalMessagesConsumed, totalMessagesPublished *prometheus.CounterVec
	var messagesPublishedLatencyMillis, messagesPublishedSizeBytes *prometheus.HistogramVec
	if meta.TotalMessagesConsumed != nil {
		totalMessagesConsumed = GetCounterVec(meta.Namespace, "pubsub_messages_consumed", "Number of messages consumed for total/success/failure scenario", meta.TotalMessagesConsumed.Labels)
	}
	if meta.TotalMessagesPublished != nil {
		totalMessagesPublished = GetCounterVec(meta.Namespace, "pubsub_messages_published", "Tracks the number of published messages at pubSub service level", meta.TotalMessagesPublished.Labels)
	}
	if meta.MessagesPublishedLatencyMillis != nil {
		messagesPublishedLatencyMillis = GetHistogramVec(meta.Namespace, "pubsub_messages_published_latency_millis", "Tracks the latencies to publish message at pubSub service level", meta.MessagesPublishedLatencyMillis.Labels, meta.MessagesPublishedLatencyMillis.Buckets)
	}
	if meta.MessagesPublishedSizeBytes != nil {
		messagesPublishedSizeBytes = GetHistogramVec(meta.Namespace, "pubsub_messages_published_size_bytes", "Tracks the message size pubSub service level", meta.MessagesPublishedSizeBytes.Labels, meta.MessagesPublishedSizeBytes.Buckets)
	}

	return &PSMetrics{
		totalMessagesConsumed:          totalMessagesConsumed,
		totalMessagesPublished:         totalMessagesPublished,
		messagesPublishedLatencyMillis: messagesPublishedLatencyMillis,
		messagesPublishedSizeBytes:     messagesPublishedSizeBytes,
	}
}

// LogMetricsPre should be called before publishing a message or when starting to process a consumed message.
// It increments the total message counters and returns the start time for latency calculation.
//
// Parameters:
//   - psMetricsLabelValues: Label values containing source, entity, operation type, and error code.
//
// Returns the start time to be passed to LogMetricsPost for latency calculation.
//
// Example:
//
//	startTime := psMetrics.LogMetricsPre(&monitoring.PSMetricsLabelValues{
//	    Source:       "orders-subscription",
//	    Entity:       "order",
//	    EntityOpType: "create",
//	})
//	// ... process message ...
//	psMetrics.LogMetricsPost(labelValues, eventTxnData)
func (psm *PSMetrics) LogMetricsPre(psMetricsLabelValues *PSMetricsLabelValues) time.Time {
	if psm.totalMessagesPublished != nil {
		psm.totalMessagesPublished.WithLabelValues(psMetricsLabelValues.Entity, psMetricsLabelValues.EntityOpType, Total).Inc()
	}
	if psm.totalMessagesConsumed != nil {
		psm.totalMessagesConsumed.WithLabelValues(string(psMetricsLabelValues.Source), psMetricsLabelValues.Entity, psMetricsLabelValues.EntityOpType, Total, "").Inc()
	}
	return time.Now()
}

// LogMetricsPost should be called after a pub/sub operation completes.
// It records the success/failure status, latency, and message size for publishing operations,
// and success/failure status for consumption operations.
//
// Parameters:
//   - psMetricsLabelValues: Label values containing source, entity, operation type, and error code.
//     Set ErrorCode to a non-empty string to indicate failure for consumed messages.
//   - eventTxnData: Transaction data from the publish operation (can be nil for consumption-only metrics).
//     Contains IsPublished flag, TimeTakenToPublish, and MessageSizeInBytes.
//
// Example (Publishing):
//
//	psMetrics.LogMetricsPost(labelValues, &pubsub.EventTxnData{
//	    IsPublished:         true,
//	    TimeTakenToPublish:  100 * time.Millisecond,
//	    MessageSizeInBytes:  2048,
//	})
//
// Example (Consumption - Success):
//
//	psMetrics.LogMetricsPost(&monitoring.PSMetricsLabelValues{
//	    Source:       "orders-subscription",
//	    Entity:       "order",
//	    EntityOpType: "create",
//	    ErrorCode:    "", // empty = success
//	}, nil)
//
// Example (Consumption - Failure):
//
//	psMetrics.LogMetricsPost(&monitoring.PSMetricsLabelValues{
//	    Source:       "orders-subscription",
//	    Entity:       "order",
//	    EntityOpType: "create",
//	    ErrorCode:    "ERR_VALIDATION",
//	}, nil)
func (psm *PSMetrics) LogMetricsPost(psMetricsLabelValues *PSMetricsLabelValues, eventTxnData *pubsub.EventTxnData) {
	if psm.totalMessagesPublished != nil && eventTxnData != nil {
		if eventTxnData.IsPublished {
			psm.totalMessagesPublished.WithLabelValues(psMetricsLabelValues.Entity, psMetricsLabelValues.EntityOpType, Success).Inc()
		} else {
			psm.totalMessagesPublished.WithLabelValues(psMetricsLabelValues.Entity, psMetricsLabelValues.EntityOpType, Failure).Inc()
		}
	}
	if psm.messagesPublishedLatencyMillis != nil && eventTxnData != nil {
		psm.messagesPublishedLatencyMillis.WithLabelValues(psMetricsLabelValues.Entity, psMetricsLabelValues.EntityOpType).Observe(float64(eventTxnData.TimeTakenToPublish.Milliseconds()))
	}
	if psm.messagesPublishedSizeBytes != nil && eventTxnData != nil {
		psm.messagesPublishedSizeBytes.WithLabelValues(psMetricsLabelValues.Entity, psMetricsLabelValues.EntityOpType).Observe(float64(eventTxnData.MessageSizeInBytes))
	}
	if psm.totalMessagesConsumed != nil {
		if psMetricsLabelValues.ErrorCode != "" {
			psm.totalMessagesConsumed.WithLabelValues(string(psMetricsLabelValues.Source), psMetricsLabelValues.Entity, psMetricsLabelValues.EntityOpType, Failure, psMetricsLabelValues.ErrorCode).Inc()
		} else {
			psm.totalMessagesConsumed.WithLabelValues(string(psMetricsLabelValues.Source), psMetricsLabelValues.Entity, psMetricsLabelValues.EntityOpType, Success, psMetricsLabelValues.ErrorCode).Inc()
		}
	}
}

// GetTotalMessagesConsumedMetric returns the underlying Prometheus CounterVec
// for the messages consumed counter. This can be used for advanced operations.
//
// Returns nil if the metric was not configured during initialization.
func (psm *PSMetrics) GetTotalMessagesConsumedMetric() *prometheus.CounterVec {
	return psm.totalMessagesConsumed
}

// GetTotalMessagesPublishedMetric returns the underlying Prometheus CounterVec
// for the messages published counter. This can be used for advanced operations.
//
// Returns nil if the metric was not configured during initialization.
func (psm *PSMetrics) GetTotalMessagesPublishedMetric() *prometheus.CounterVec {
	return psm.totalMessagesPublished
}

// GetMessagesPublishedLatencyMillisMetric returns the underlying Prometheus HistogramVec
// for the message publish latency. This can be used for advanced operations.
//
// Returns nil if the metric was not configured during initialization.
func (psm *PSMetrics) GetMessagesPublishedLatencyMillisMetric() *prometheus.HistogramVec {
	return psm.messagesPublishedLatencyMillis
}

// GetMessagesPublishedSizeBytesMetric returns the underlying Prometheus HistogramVec
// for the published message size. This can be used for advanced operations.
//
// Returns nil if the metric was not configured during initialization.
func (psm *PSMetrics) GetMessagesPublishedSizeBytesMetric() *prometheus.HistogramVec {
	return psm.messagesPublishedSizeBytes
}
