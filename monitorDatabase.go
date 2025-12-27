package app_monitoring

import (
	"time"

	ae "github.com/piyushkumar96/app-error"
	"github.com/prometheus/client_golang/prometheus"
)

// NewDatabaseMetrics creates and registers database operation Prometheus metrics.
// It initializes counters for operation counts and histograms for operation latencies.
//
// The metrics track:
//   - OperationsTotal: Counter for total/success/failure database operations
//   - OperationsLatencyMillis: Histogram for operation duration in milliseconds
//
// Parameters:
//   - meta: Configuration containing the namespace and metric settings.
//     Set individual metric configs to nil to disable them.
//
// Returns a DBMetrics instance that can be used to log database operation metrics.
//
// Example:
//
//	dbMetrics := monitoring.NewDatabaseMetrics(&monitoring.DBMetricsMeta{
//	    Namespace: "myapp",
//	    OperationsTotal: &monitoring.MetricMeta{
//	        Labels: []string{"op_type", "source", "entity", "is_txn", "status"},
//	    },
//	    OperationsLatencyMillis: &monitoring.MetricMeta{
//	        Labels:  []string{"op_type", "source", "entity", "is_txn"},
//	        Buckets: monitoring.GetExponentialBuckets(1, 2, 12),
//	    },
//	})
func NewDatabaseMetrics(meta *DBMetricsMeta) *DBMetrics {
	var operationsTotal *prometheus.CounterVec
	var operationsLatencyMillis *prometheus.HistogramVec

	if meta.OperationsTotal != nil {
		operationsTotal = GetCounterVec(meta.Namespace, "db_operations", "Number of times DB operations executed for total/success/failure", meta.OperationsTotal.Labels)
	}
	if meta.OperationsLatencyMillis != nil {
		operationsLatencyMillis = GetHistogramVec(meta.Namespace, "db_operations_latency_millis", "Tracks the latencies for HTTP requests at spring serve level", meta.OperationsLatencyMillis.Labels, meta.OperationsLatencyMillis.Buckets)
	}

	return &DBMetrics{
		operationsTotal:         operationsTotal,
		operationsLatencyMillis: operationsLatencyMillis,
	}
}

// LogMetricsPre should be called before executing a database operation.
// It increments the total operations counter and returns the start time for latency calculation.
//
// Parameters:
//   - dbMetricsLabelValues: Label values containing operation details (type, source, entity, transaction flag).
//
// Returns the start time to be passed to LogMetricsPost for latency calculation.
//
// Example:
//
//	startTime := dbMetrics.LogMetricsPre(&monitoring.DBMetricsLabelValues{
//	    OpType:   "select",
//	    Source:   "UserRepository",
//	    AdEntity: "users",
//	    IsTxn:    "false",
//	})
//	// ... execute database operation ...
//	dbMetrics.LogMetricsPost(err, labelValues, startTime)
func (dm *DBMetrics) LogMetricsPre(dbMetricsLabelValues *DBMetricsLabelValues) time.Time {
	if dm.operationsTotal != nil {
		dm.operationsTotal.WithLabelValues(string(dbMetricsLabelValues.OpType), string(dbMetricsLabelValues.Source), string(dbMetricsLabelValues.AdEntity), dbMetricsLabelValues.IsTxn, Total).Inc()
	}
	return time.Now()
}

// LogMetricsPost should be called after a database operation completes.
// It records the success/failure status and the operation latency.
//
// Parameters:
//   - appErr: The error returned by the operation (nil for success, non-nil for failure).
//   - dbMetricsLabelValues: Label values containing operation details.
//   - opsExecTime: The start time returned by LogMetricsPre.
//
// Example:
//
//	dbMetrics.LogMetricsPost(nil, labelValues, startTime) // Success
//	dbMetrics.LogMetricsPost(err, labelValues, startTime) // Failure
func (dm *DBMetrics) LogMetricsPost(appErr *ae.AppError, dbMetricsLabelValues *DBMetricsLabelValues, opsExecTime time.Time) {
	if dm.operationsTotal != nil {
		if appErr != nil {
			dm.operationsTotal.WithLabelValues(string(dbMetricsLabelValues.OpType), string(dbMetricsLabelValues.Source), dbMetricsLabelValues.AdEntity, dbMetricsLabelValues.IsTxn, Failure).Inc()
		} else {
			dm.operationsTotal.WithLabelValues(string(dbMetricsLabelValues.OpType), string(dbMetricsLabelValues.Source), dbMetricsLabelValues.AdEntity, dbMetricsLabelValues.IsTxn, Success).Inc()
		}
	}
	if dm.operationsLatencyMillis != nil {
		dm.operationsLatencyMillis.WithLabelValues(string(dbMetricsLabelValues.OpType), string(dbMetricsLabelValues.Source), dbMetricsLabelValues.AdEntity, dbMetricsLabelValues.IsTxn).Observe(float64(time.Since(opsExecTime).Milliseconds()))
	}
}

// GetOperationsTotalMetric returns the underlying Prometheus CounterVec
// for the database operations counter. This can be used for advanced operations.
//
// Returns nil if the metric was not configured during initialization.
func (dm *DBMetrics) GetOperationsTotalMetric() *prometheus.CounterVec {
	return dm.operationsTotal
}

// GetOperationsLatencyMillisMetric returns the underlying Prometheus HistogramVec
// for the database operation latency. This can be used for advanced operations.
//
// Returns nil if the metric was not configured during initialization.
func (dm *DBMetrics) GetOperationsLatencyMillisMetric() *prometheus.HistogramVec {
	return dm.operationsLatencyMillis
}
