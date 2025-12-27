package prometheus

import (
	"time"

	"github.com/piyushkumar96/app-monitoring/constants"
	"github.com/piyushkumar96/app-monitoring/interfaces"
	"github.com/piyushkumar96/app-monitoring/models"

	pubsub "github.com/piyushkumar96/generic-pubsub"
	"github.com/prometheus/client_golang/prometheus"
)

// NewPromPubSubMetrics creates and registers Prometheus metrics for pub/sub messaging operations.
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
// Returns an interfaces.PSMetricsInterface instance for logging pub/sub messaging metrics.
func NewPromPubSubMetrics(meta *models.PSMetricsMeta) interfaces.PSMetricsInterface {
	var totalMessagesConsumed, totalMessagesPublished *prometheus.CounterVec
	var messagesPublishedLatencyMillis, messagesPublishedSizeBytes *prometheus.HistogramVec
	if meta.TotalMessagesConsumed != nil {
		totalMessagesConsumed = GetPromCounterVec(meta.Namespace, "pubsub_messages_consumed", "Number of messages consumed for total/success/failure scenario", meta.TotalMessagesConsumed.Labels)
	}
	if meta.TotalMessagesPublished != nil {
		totalMessagesPublished = GetPromCounterVec(meta.Namespace, "pubsub_messages_published", "Tracks the number of published messages at pubSub service level", meta.TotalMessagesPublished.Labels)
	}
	if meta.MessagesPublishedLatencyMillis != nil {
		messagesPublishedLatencyMillis = GetPromHistogramVec(meta.Namespace, "pubsub_messages_published_latency_millis", "Tracks the latencies to publish message at pubSub service level", meta.MessagesPublishedLatencyMillis.Labels, meta.MessagesPublishedLatencyMillis.Buckets)
	}
	if meta.MessagesPublishedSizeBytes != nil {
		messagesPublishedSizeBytes = GetPromHistogramVec(meta.Namespace, "pubsub_messages_published_size_bytes", "Tracks the message size pubSub service level", meta.MessagesPublishedSizeBytes.Labels, meta.MessagesPublishedSizeBytes.Buckets)
	}

	return &PromPSMetrics{
		totalMessagesConsumed:          totalMessagesConsumed,
		totalMessagesPublished:         totalMessagesPublished,
		messagesPublishedLatencyMillis: messagesPublishedLatencyMillis,
		messagesPublishedSizeBytes:     messagesPublishedSizeBytes,
	}
}

// LogMetricsPre should be called before publishing a message or when starting to process a consumed message.
// It increments the total message counters and returns the start time for latency calculation.
func (psm *PromPSMetrics) LogMetricsPre(psMetricsLabelValues *models.PSMetricsLabelValues) time.Time {
	if psm.totalMessagesPublished != nil {
		psm.totalMessagesPublished.WithLabelValues(psMetricsLabelValues.Entity, psMetricsLabelValues.EntityOpType, constants.Total).Inc()
	}
	if psm.totalMessagesConsumed != nil {
		psm.totalMessagesConsumed.WithLabelValues(string(psMetricsLabelValues.Source), psMetricsLabelValues.Entity, psMetricsLabelValues.EntityOpType, constants.Total, "").Inc()
	}
	return time.Now()
}

// LogMetricsPost should be called after a pub/sub operation completes.
// It records the success/failure status, latency, and message size for publishing operations,
// and success/failure status for consumption operations.
func (psm *PromPSMetrics) LogMetricsPost(psMetricsLabelValues *models.PSMetricsLabelValues, eventTxnData *pubsub.EventTxnData) {
	if psm.totalMessagesPublished != nil && eventTxnData != nil {
		if eventTxnData.IsPublished {
			psm.totalMessagesPublished.WithLabelValues(psMetricsLabelValues.Entity, psMetricsLabelValues.EntityOpType, constants.Success).Inc()
		} else {
			psm.totalMessagesPublished.WithLabelValues(psMetricsLabelValues.Entity, psMetricsLabelValues.EntityOpType, constants.Failure).Inc()
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
			psm.totalMessagesConsumed.WithLabelValues(string(psMetricsLabelValues.Source), psMetricsLabelValues.Entity, psMetricsLabelValues.EntityOpType, constants.Failure, psMetricsLabelValues.ErrorCode).Inc()
		} else {
			psm.totalMessagesConsumed.WithLabelValues(string(psMetricsLabelValues.Source), psMetricsLabelValues.Entity, psMetricsLabelValues.EntityOpType, constants.Success, psMetricsLabelValues.ErrorCode).Inc()
		}
	}
}

// GetTotalMessagesConsumedMetric returns the underlying Prometheus CounterVec
// for the messages consumed counter. This can be used for advanced operations.
func (psm *PromPSMetrics) GetTotalMessagesConsumedMetric() *prometheus.CounterVec {
	return psm.totalMessagesConsumed
}

// GetTotalMessagesPublishedMetric returns the underlying Prometheus CounterVec
// for the messages published counter. This can be used for advanced operations.
func (psm *PromPSMetrics) GetTotalMessagesPublishedMetric() *prometheus.CounterVec {
	return psm.totalMessagesPublished
}

// GetMessagesPublishedLatencyMillisMetric returns the underlying Prometheus HistogramVec
// for the message publish latency. This can be used for advanced operations.
func (psm *PromPSMetrics) GetMessagesPublishedLatencyMillisMetric() *prometheus.HistogramVec {
	return psm.messagesPublishedLatencyMillis
}

// GetMessagesPublishedSizeBytesMetric returns the underlying Prometheus HistogramVec
// for the published message size. This can be used for advanced operations.
func (psm *PromPSMetrics) GetMessagesPublishedSizeBytesMetric() *prometheus.HistogramVec {
	return psm.messagesPublishedSizeBytes
}
