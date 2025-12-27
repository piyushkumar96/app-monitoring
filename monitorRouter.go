package app_monitoring

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// NewRouterLevelMetrics creates and registers Prometheus metrics for HTTP router/endpoint level monitoring.
// It initializes counters for request counts and histograms for latencies and payload sizes.
//
// The metrics track:
//   - HTTPRequests: Counter for total/success/failure HTTP requests
//   - HTTPRequestsLatencyMillis: Histogram for request latency in milliseconds
//   - HTTPRequestSizeBytes: Histogram for request body size in bytes
//   - HTTPResponseSizeBytes: Histogram for response body size in bytes
//
// Parameters:
//   - meta: Configuration containing the namespace and metric settings.
//     Set individual metric configs to nil to disable them.
//
// Returns a RouterMetrics instance for logging HTTP endpoint metrics.
//
// Example:
//
//	routerMetrics := monitoring.NewRouterLevelMetrics(&monitoring.RouterMetricsMeta{
//	    Namespace: "myapp",
//	    HTTPRequests: &monitoring.MetricMeta{
//	        Labels: []string{"method", "code", "path", "status"},
//	    },
//	    HTTPRequestsLatencyMillis: &monitoring.MetricMeta{
//	        Labels:  []string{"method", "code", "path"},
//	        Buckets: monitoring.GetExponentialBuckets(10, 2, 10),
//	    },
//	})
func NewRouterLevelMetrics(meta *RouterMetricsMeta) *RouterMetrics {
	var httpRequests *prometheus.CounterVec
	var httpRequestsLatencyMillis, httpRequestSizeBytes, httpResponseSizeBytes *prometheus.HistogramVec

	if meta.HTTPRequests != nil {
		httpRequests = GetCounterVec(meta.Namespace, "http_requests", "Tracks the number of HTTP requests at application level", meta.HTTPRequests.Labels)
	}
	if meta.HTTPRequestsLatencyMillis != nil {
		httpRequestsLatencyMillis = GetHistogramVec(meta.Namespace, "http_request_latency_millis", "Tracks the latencies for HTTP requests at application level", meta.HTTPRequestsLatencyMillis.Labels, meta.HTTPRequestsLatencyMillis.Buckets)
	}
	if meta.HTTPRequestSizeBytes != nil {
		httpRequestSizeBytes = GetHistogramVec(meta.Namespace, "http_request_size_bytes", "Tracks the size of HTTP requests at application level.", meta.HTTPRequestSizeBytes.Labels, meta.HTTPRequestSizeBytes.Buckets)
	}
	if meta.HTTPResponseSizeBytes != nil {
		httpResponseSizeBytes = GetHistogramVec(meta.Namespace, "http_response_size_bytes", "Tracks the size of HTTP responses at application level", meta.HTTPResponseSizeBytes.Labels, meta.HTTPResponseSizeBytes.Buckets)
	}

	return &RouterMetrics{
		httpRequests:              httpRequests,
		httpRequestsLatencyMillis: httpRequestsLatencyMillis,
		httpRequestSizeBytes:      httpRequestSizeBytes,
		httpResponseSizeBytes:     httpResponseSizeBytes,
	}
}

// LogMetrics returns a Gin middleware that automatically logs Prometheus metrics for all HTTP requests.
// It captures request counts, latencies, and payload sizes for each endpoint.
//
// The middleware:
//   - Skips metrics collection for the metrics endpoint itself (to avoid self-referential metrics)
//   - Increments total request count before processing
//   - Records success/failure based on HTTP status code (2XX = success)
//   - Measures request latency, request size, and response size
//
// Parameters:
//   - metricsPath: The path where Prometheus metrics are exposed (e.g., "/metrics").
//     Requests to this path will not be recorded to avoid metric pollution.
//
// Returns a Gin HandlerFunc that can be used as middleware.
//
// Example:
//
//	router := gin.Default()
//	routerMetrics := monitoring.NewRouterLevelMetrics(meta)
//	router.Use(routerMetrics.LogMetrics("/metrics"))
func (rlm *RouterMetrics) LogMetrics(metricsPath string) gin.HandlerFunc {
	return func(gc *gin.Context) {
		// Skip metrics collection for the metrics endpoint itself
		if gc.Request.URL.Path == metricsPath {
			gc.Next()
			return
		}

		start := time.Now()
		reqSize := float64(computeApproximateRequestSize(gc.Request))
		urlPath := gc.FullPath()

		if rlm.httpRequests != nil {
			// Increment total request counter before processing
			rlm.httpRequests.WithLabelValues(gc.Request.Method, "", urlPath, Total).Inc()
		}

		// Pass request to the next handler in chain
		gc.Next()

		// Collect response metrics after handler completes
		httpCode := strconv.Itoa(gc.Writer.Status())
		elapsed := float64(time.Since(start)) / float64(time.Millisecond)
		respSize := float64(gc.Writer.Size())

		// Parse HTTP code for success/failure determination
		httpCodeInt, err := strconv.ParseInt(httpCode, 10, 32)
		if err != nil {
			httpCodeInt = 0
		}

		// Record success/failure based on HTTP status code
		if rlm.httpRequests != nil {
			if httpCodeInt >= HTTPStatus2XXMinValue && httpCodeInt <= HTTPStatus2XXMaxValue {
				rlm.httpRequests.WithLabelValues(gc.Request.Method, httpCode, urlPath, Success).Inc()
			} else {
				rlm.httpRequests.WithLabelValues(gc.Request.Method, httpCode, urlPath, Failure).Inc()
			}
		}

		// Record latency histogram
		if rlm.httpRequestsLatencyMillis != nil {
			rlm.httpRequestsLatencyMillis.WithLabelValues(gc.Request.Method, httpCode, urlPath).Observe(elapsed)
		}

		// Record request size histogram
		if rlm.httpRequestSizeBytes != nil {
			rlm.httpRequestSizeBytes.WithLabelValues(gc.Request.Method, httpCode, urlPath).Observe(reqSize)
		}

		// Record response size histogram
		if rlm.httpResponseSizeBytes != nil {
			rlm.httpResponseSizeBytes.WithLabelValues(gc.Request.Method, httpCode, urlPath).Observe(respSize)
		}
	}
}

// computeApproximateRequestSize calculates an approximate size of the HTTP request in bytes.
// It includes the URL path, method, protocol, headers, host, and content length.
func computeApproximateRequestSize(r *http.Request) int {
	totalSize := 0
	if r.URL != nil {
		totalSize = len(r.URL.Path)
	}

	totalSize += len(r.Method) + len(r.Proto)
	for name, values := range r.Header {
		totalSize += len(name)
		for _, value := range values {
			totalSize += len(value)
		}
	}
	totalSize += len(r.Host)
	if r.ContentLength != -1 {
		totalSize += int(r.ContentLength)
	}
	return totalSize
}

// GetHTTPRequestsMetric returns the underlying Prometheus CounterVec
// for the HTTP requests counter. This can be used for advanced operations.
//
// Returns nil if the metric was not configured during initialization.
func (rlm *RouterMetrics) GetHTTPRequestsMetric() *prometheus.CounterVec {
	return rlm.httpRequests
}

// GetHTTPRequestsLatencyMillisMetric returns the underlying Prometheus HistogramVec
// for the HTTP request latency. This can be used for advanced operations.
//
// Returns nil if the metric was not configured during initialization.
func (rlm *RouterMetrics) GetHTTPRequestsLatencyMillisMetric() *prometheus.HistogramVec {
	return rlm.httpRequestsLatencyMillis
}

// GetHTTPRequestSizeBytesMetric returns the underlying Prometheus HistogramVec
// for the HTTP request size. This can be used for advanced operations.
//
// Returns nil if the metric was not configured during initialization.
func (rlm *RouterMetrics) GetHTTPRequestSizeBytesMetric() *prometheus.HistogramVec {
	return rlm.httpRequestSizeBytes
}

// GetHTTPResponseSizeBytesMetric returns the underlying Prometheus HistogramVec
// for the HTTP response size. This can be used for advanced operations.
//
// Returns nil if the metric was not configured during initialization.
func (rlm *RouterMetrics) GetHTTPResponseSizeBytesMetric() *prometheus.HistogramVec {
	return rlm.httpResponseSizeBytes
}
