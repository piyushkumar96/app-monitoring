package prometheus

import (
	"strconv"

	"github.com/piyushkumar96/app-monitoring/constants"
	"github.com/piyushkumar96/app-monitoring/interfaces"
	"github.com/piyushkumar96/app-monitoring/models"

	"github.com/prometheus/client_golang/prometheus"
)

// NewPromDownstreamServiceMetrics creates and registers Prometheus metrics for downstream HTTP service calls.
// It initializes counters for request counts and histograms for latencies and payload sizes.
//
// The metrics track:
//   - HTTPRequests: Counter for total/success/failure HTTP requests to downstream services
//   - HTTPRequestsLatencyMillis: Histogram for request latency in milliseconds
//   - HTTPRequestSizeBytes: Histogram for request body size in bytes
//   - HTTPResponseSizeBytes: Histogram for response body size in bytes
//
// Parameters:
//   - meta: Configuration containing the namespace and metric settings.
//     Set individual metric configs to nil to disable them.
//
// Returns an interfaces.DownstreamServiceMetricsInterface instance for logging downstream call metrics.
func NewPromDownstreamServiceMetrics(meta *models.DownstreamServiceMetricsMeta) interfaces.DownstreamServiceMetricsInterface {
	var httpRequests *prometheus.CounterVec
	var httpRequestsLatencyMillis, httpRequestSizeBytes, httpResponseSizeBytes *prometheus.HistogramVec

	if meta.HTTPRequests != nil {
		httpRequests = GetPromCounterVec(meta.Namespace, "downstream_service_http_requests", "Tracks the number of HTTP requests at downstream service level", meta.HTTPRequests.Labels)
	}
	if meta.HTTPRequestsLatencyMillis != nil {
		httpRequestsLatencyMillis = GetPromHistogramVec(meta.Namespace, "downstream_service_http_request_latency_millis", "Tracks the latencies for HTTP requests at downstream service level", meta.HTTPRequestsLatencyMillis.Labels, meta.HTTPRequestsLatencyMillis.Buckets)
	}
	if meta.HTTPRequestSizeBytes != nil {
		httpRequestSizeBytes = GetPromHistogramVec(meta.Namespace, "downstream_service_http_request_size_bytes", "Tracks the size of HTTP requests at downstream service level.", meta.HTTPRequestSizeBytes.Labels, meta.HTTPRequestSizeBytes.Buckets)
	}
	if meta.HTTPResponseSizeBytes != nil {
		httpResponseSizeBytes = GetPromHistogramVec(meta.Namespace, "downstream_service_http_response_size_bytes", "Tracks the size of HTTP responses at downstream service level", meta.HTTPResponseSizeBytes.Labels, meta.HTTPResponseSizeBytes.Buckets)
	}

	return &PromDownstreamServiceMetrics{
		httpRequests:              httpRequests,
		httpRequestsLatencyMillis: httpRequestsLatencyMillis,
		httpRequestSizeBytes:      httpRequestSizeBytes,
		httpResponseSizeBytes:     httpResponseSizeBytes,
	}
}

// LogMetricsPre should be called before making a downstream service HTTP call.
// It increments the total request counter for the service.
func (dsm *PromDownstreamServiceMetrics) LogMetricsPre(dssMetricsLabelValues *models.DownstreamServiceMetricsLabelValues) {
	if dsm.httpRequests != nil {
		dsm.httpRequests.WithLabelValues(string(dssMetricsLabelValues.Name), dssMetricsLabelValues.HTTPMethod, "", dssMetricsLabelValues.APIIdentifier, constants.Total).Inc()
	}
}

// LogMetricsPost should be called after a downstream service HTTP call completes.
// It records the success/failure status, latency, and payload sizes.
func (dsm *PromDownstreamServiceMetrics) LogMetricsPost(success bool, dssMetricsLabelValues *models.DownstreamServiceMetricsLabelValues, httpMetrics *models.HTTPMetrics) {
	httpCodeStr := strconv.Itoa(httpMetrics.Code)
	if dsm.httpRequests != nil {
		if success {
			dsm.httpRequests.WithLabelValues(string(dssMetricsLabelValues.Name), httpMetrics.Method, httpCodeStr, dssMetricsLabelValues.APIIdentifier, constants.Success).Inc()
		} else {
			dsm.httpRequests.WithLabelValues(string(dssMetricsLabelValues.Name), httpMetrics.Method, httpCodeStr, dssMetricsLabelValues.APIIdentifier, constants.Failure).Inc()
		}
	}
	if dsm.httpRequestsLatencyMillis != nil {
		dsm.httpRequestsLatencyMillis.WithLabelValues(string(dssMetricsLabelValues.Name), httpMetrics.Method, httpCodeStr, dssMetricsLabelValues.APIIdentifier).Observe(float64(httpMetrics.ResponseTime.Milliseconds()))
	}
	if dsm.httpRequestSizeBytes != nil {
		dsm.httpRequestSizeBytes.WithLabelValues(string(dssMetricsLabelValues.Name), httpMetrics.Method, httpCodeStr, dssMetricsLabelValues.APIIdentifier).Observe(float64(httpMetrics.RequestBodySizeBytes))
	}
	if dsm.httpResponseSizeBytes != nil {
		dsm.httpResponseSizeBytes.WithLabelValues(string(dssMetricsLabelValues.Name), httpMetrics.Method, httpCodeStr, dssMetricsLabelValues.APIIdentifier).Observe(float64(httpMetrics.ResponseBodySizeBytes))
	}
}

// GetHTTPRequestsMetric returns the underlying Prometheus CounterVec
// for the HTTP requests counter. This can be used for advanced operations.
func (dsm *PromDownstreamServiceMetrics) GetHTTPRequestsMetric() *prometheus.CounterVec {
	return dsm.httpRequests
}

// GetHTTPRequestsLatencyMillisMetric returns the underlying Prometheus HistogramVec
// for the HTTP request latency. This can be used for advanced operations.
func (dsm *PromDownstreamServiceMetrics) GetHTTPRequestsLatencyMillisMetric() *prometheus.HistogramVec {
	return dsm.httpRequestsLatencyMillis
}

// GetHTTPRequestSizeBytesMetric returns the underlying Prometheus HistogramVec
// for the HTTP request size. This can be used for advanced operations.
func (dsm *PromDownstreamServiceMetrics) GetHTTPRequestSizeBytesMetric() *prometheus.HistogramVec {
	return dsm.httpRequestSizeBytes
}

// GetHTTPResponseSizeBytesMetric returns the underlying Prometheus HistogramVec
// for the HTTP response size. This can be used for advanced operations.
func (dsm *PromDownstreamServiceMetrics) GetHTTPResponseSizeBytesMetric() *prometheus.HistogramVec {
	return dsm.httpResponseSizeBytes
}
