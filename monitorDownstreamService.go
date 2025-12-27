package app_monitoring

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

// NewDownstreamServiceMetrics creates and registers Prometheus metrics for downstream HTTP service calls.
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
// Returns a DownstreamServiceMetrics instance for logging downstream call metrics.
//
// Example:
//
//	dsMetrics := monitoring.NewDownstreamServiceMetrics(&monitoring.DownstreamServiceMetricsMeta{
//	    Namespace: "myapp",
//	    HTTPRequests: &monitoring.MetricMeta{
//	        Labels: []string{"service", "method", "code", "api", "status"},
//	    },
//	    HTTPRequestsLatencyMillis: &monitoring.MetricMeta{
//	        Labels:  []string{"service", "method", "code", "api"},
//	        Buckets: monitoring.GetExponentialBuckets(10, 2, 10),
//	    },
//	})
func NewDownstreamServiceMetrics(meta *DownstreamServiceMetricsMeta) *DownstreamServiceMetrics {
	var httpRequests *prometheus.CounterVec
	var httpRequestsLatencyMillis, httpRequestSizeBytes, httpResponseSizeBytes *prometheus.HistogramVec

	if meta.HTTPRequests != nil {
		httpRequests = GetCounterVec(meta.Namespace, "downstream_service_http_requests", "Tracks the number of HTTP requests at downstream service level", meta.HTTPRequests.Labels)
	}
	if meta.HTTPRequestsLatencyMillis != nil {
		httpRequestsLatencyMillis = GetHistogramVec(meta.Namespace, "downstream_service_http_request_latency_millis", "Tracks the latencies for HTTP requests at downstream service level", meta.HTTPRequestsLatencyMillis.Labels, meta.HTTPRequestsLatencyMillis.Buckets)
	}
	if meta.HTTPRequestSizeBytes != nil {
		httpRequestSizeBytes = GetHistogramVec(meta.Namespace, "downstream_service_http_request_size_bytes", "Tracks the size of HTTP requests at downstream service level.", meta.HTTPRequestSizeBytes.Labels, meta.HTTPRequestSizeBytes.Buckets)
	}
	if meta.HTTPResponseSizeBytes != nil {
		httpResponseSizeBytes = GetHistogramVec(meta.Namespace, "downstream_service_http_response_size_bytes", "Tracks the size of HTTP responses at downstream service level", meta.HTTPResponseSizeBytes.Labels, meta.HTTPResponseSizeBytes.Buckets)
	}

	return &DownstreamServiceMetrics{
		httpRequests:              httpRequests,
		httpRequestsLatencyMillis: httpRequestsLatencyMillis,
		httpRequestSizeBytes:      httpRequestSizeBytes,
		httpResponseSizeBytes:     httpResponseSizeBytes,
	}
}

// LogMetricsPre should be called before making a downstream service HTTP call.
// It increments the total request counter for the service.
//
// Parameters:
//   - dssMetricsLabelValues: Label values containing service name, HTTP method, and API identifier.
//
// Example:
//
//	dsMetrics.LogMetricsPre(&monitoring.DownstreamServiceMetricsLabelValues{
//	    Name:          "payment-service",
//	    HTTPMethod:    "POST",
//	    APIIdentifier: "/api/v1/payments",
//	})
func (dsm *DownstreamServiceMetrics) LogMetricsPre(dssMetricsLabelValues *DownstreamServiceMetricsLabelValues) {
	if dsm.httpRequests != nil {
		dsm.httpRequests.WithLabelValues(string(dssMetricsLabelValues.Name), dssMetricsLabelValues.HTTPMethod, "", dssMetricsLabelValues.APIIdentifier, Total).Inc()
	}
}

// LogMetricsPost should be called after a downstream service HTTP call completes.
// It records the success/failure status, latency, and payload sizes.
//
// Parameters:
//   - success: Whether the call was successful (typically based on HTTP status code).
//   - dssMetricsLabelValues: Label values containing service name, HTTP method, and API identifier.
//   - httpMetrics: HTTP metrics containing response code, latency, and payload sizes.
//
// Example:
//
//	dsMetrics.LogMetricsPost(true, labelValues, &monitoring.HTTPMetrics{
//	    Method:                "POST",
//	    Code:                  200,
//	    RequestBodySizeBytes:  1024,
//	    ResponseBodySizeBytes: 512,
//	    ResponseTime:          150 * time.Millisecond,
//	})
func (dsm *DownstreamServiceMetrics) LogMetricsPost(success bool, dssMetricsLabelValues *DownstreamServiceMetricsLabelValues, httpMetrics *HTTPMetrics) {
	httpCodeStr := strconv.Itoa(httpMetrics.Code)
	if dsm.httpRequests != nil {
		if success {
			dsm.httpRequests.WithLabelValues(string(dssMetricsLabelValues.Name), httpMetrics.Method, httpCodeStr, dssMetricsLabelValues.APIIdentifier, Success).Inc()
		} else {
			dsm.httpRequests.WithLabelValues(string(dssMetricsLabelValues.Name), httpMetrics.Method, httpCodeStr, dssMetricsLabelValues.APIIdentifier, Failure).Inc()
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
//
// Returns nil if the metric was not configured during initialization.
func (dsm *DownstreamServiceMetrics) GetHTTPRequestsMetric() *prometheus.CounterVec {
	return dsm.httpRequests
}

// GetHTTPRequestsLatencyMillisMetric returns the underlying Prometheus HistogramVec
// for the HTTP request latency. This can be used for advanced operations.
//
// Returns nil if the metric was not configured during initialization.
func (dsm *DownstreamServiceMetrics) GetHTTPRequestsLatencyMillisMetric() *prometheus.HistogramVec {
	return dsm.httpRequestsLatencyMillis
}

// GetHTTPRequestSizeBytesMetric returns the underlying Prometheus HistogramVec
// for the HTTP request size. This can be used for advanced operations.
//
// Returns nil if the metric was not configured during initialization.
func (dsm *DownstreamServiceMetrics) GetHTTPRequestSizeBytesMetric() *prometheus.HistogramVec {
	return dsm.httpRequestSizeBytes
}

// GetHTTPResponseSizeBytesMetric returns the underlying Prometheus HistogramVec
// for the HTTP response size. This can be used for advanced operations.
//
// Returns nil if the metric was not configured during initialization.
func (dsm *DownstreamServiceMetrics) GetHTTPResponseSizeBytesMetric() *prometheus.HistogramVec {
	return dsm.httpResponseSizeBytes
}
