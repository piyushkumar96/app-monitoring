package prometheus

import (
	l "github.com/piyushkumar96/generic-logger"
	"github.com/prometheus/client_golang/prometheus"
)

// GetPromHistogramVec creates and registers a new Prometheus HistogramVec metric.
// A histogram samples observations (usually things like request durations or response sizes)
// and counts them in configurable buckets.
//
// Parameters:
//   - namespace: The metric namespace (typically the application name)
//   - name: The metric name
//   - help: Description of what the metric measures
//   - labelNames: Slice of label names for the metric dimensions
//   - buckets: Histogram bucket boundaries (e.g., []float64{10, 50, 100, 500, 1000})
//
// Returns a HistogramVec that can be used to observe values with different label combinations.
// If registration fails (e.g., duplicate metric), an error is logged but the histogram is still returned.
func GetPromHistogramVec(namespace, name, help string, labelNames []string, buckets []float64) *prometheus.HistogramVec {
	histogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      name,
			Help:      help,
			Buckets:   buckets,
		}, labelNames,
	)
	if err := prometheus.Register(histogram); err != nil {
		l.Logger.Error("failed to register histogram vec metric", "code", "OnHistogramMetricRegisterFailure", "err", err.Error())
	}
	return histogram
}

// GetPromSummaryVec creates and registers a new Prometheus SummaryVec metric.
// A summary samples observations and provides a total count and sum of observations,
// along with configurable quantiles over a sliding time window.
//
// Parameters:
//   - namespace: The metric namespace (typically the application name)
//   - name: The metric name
//   - help: Description of what the metric measures
//   - labelNames: Slice of label names for the metric dimensions
//
// Returns a SummaryVec that can be used to observe values with different label combinations.
// If registration fails (e.g., duplicate metric), an error is logged but the summary is still returned.
func GetPromSummaryVec(namespace, name, help string, labelNames []string) *prometheus.SummaryVec {
	summary := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: namespace,
			Name:      name,
			Help:      help,
		}, labelNames,
	)
	if err := prometheus.Register(summary); err != nil {
		l.Logger.Error("failed to register summary vec metric", "code", "OnSummaryVecMetricRegisterFailure", "err", err.Error())
	}
	return summary
}

// GetPromCounterVec creates and registers a new Prometheus CounterVec metric.
// A counter is a cumulative metric that only increases (or resets to zero on restart).
// Use counters for things like number of requests, errors, or completed tasks.
//
// Parameters:
//   - namespace: The metric namespace (typically the application name)
//   - name: The metric name
//   - help: Description of what the metric measures
//   - labelNames: Slice of label names for the metric dimensions
//
// Returns a CounterVec that can be used to increment counts with different label combinations.
// If registration fails (e.g., duplicate metric), an error is logged but the counter is still returned.
func GetPromCounterVec(namespace, name, help string, labelNames []string) *prometheus.CounterVec {
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      name,
			Help:      help,
		}, labelNames,
	)
	if err := prometheus.Register(counter); err != nil {
		l.Logger.Error("failed to register counter vec metric", "code", "OnCounterVecMetricRegisterFailure", "err", err.Error())
	}
	return counter
}

// GetPromGaugeVec creates and registers a new Prometheus GaugeVec metric.
// A gauge is a metric that can go up and down, representing a current value.
// Use gauges for things like current temperature, memory usage, or active connections.
//
// Parameters:
//   - namespace: The metric namespace (typically the application name)
//   - name: The metric name
//   - help: Description of what the metric measures
//   - labelNames: Slice of label names for the metric dimensions
//
// Returns a GaugeVec that can be used to set, increment, or decrement values with different label combinations.
// If registration fails (e.g., duplicate metric), an error is logged but the gauge is still returned.
func GetPromGaugeVec(namespace, name, help string, labelNames []string) *prometheus.GaugeVec {
	gauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      name,
			Help:      help,
		}, labelNames,
	)
	if err := prometheus.Register(gauge); err != nil {
		l.Logger.Error("failed to register gaugevec metric", "code", "OnGaugeVecMetricRegisterFailure", "err", err.Error())
	}
	return gauge
}

// GetPromExponentialBuckets generates exponentially increasing bucket boundaries for histograms.
// This is useful for latency measurements where you expect a wide range of values.
//
// Parameters:
//   - start: The lower bound of the first bucket (must be > 0)
//   - factor: The growth factor between consecutive buckets (must be > 1)
//   - count: The total number of buckets to generate
//
// Example: GetPromExponentialBuckets(10, 2, 5) returns []float64{10, 20, 40, 80, 160}
//
// Returns a slice of float64 bucket boundaries suitable for use with GetPromHistogramVec.
func GetPromExponentialBuckets(start, factor float64, count int) []float64 {
	return prometheus.ExponentialBuckets(start, factor, count)
}
