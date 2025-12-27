package app_monitoring

import (
	"time"

	ae "github.com/piyushkumar96/app-error"
	"github.com/prometheus/client_golang/prometheus"
)

// NewCronJobMetrics creates and registers cron job execution Prometheus metrics.
// It initializes counters for job execution counts and histograms for job latencies.
//
// The metrics track:
//   - JobExecutionTotal: Counter for total/success/failure job executions
//   - JobExecutionLatencyMillis: Histogram for job execution duration in milliseconds
//
// Parameters:
//   - meta: Configuration containing the namespace and metric settings.
//     Set individual metric configs to nil to disable them.
//
// Returns a CronJobMetrics instance that can be used to log job execution metrics.
//
// Example:
//
//	cronMetrics := monitoring.NewCronJobMetrics(&monitoring.CronJobMetricsMeta{
//	    Namespace: "myapp",
//	    JobExecutionTotal: &monitoring.MetricMeta{
//	        Labels: []string{"job_name", "status"},
//	    },
//	    JobExecutionLatencyMillis: &monitoring.MetricMeta{
//	        Labels:  []string{"job_name"},
//	        Buckets: monitoring.GetExponentialBuckets(100, 2, 10),
//	    },
//	})
func NewCronJobMetrics(meta *CronJobMetricsMeta) *CronJobMetrics {
	var jobExecutionTotal *prometheus.CounterVec
	var jobExecutionLatencyMillis *prometheus.HistogramVec

	if meta.JobExecutionTotal != nil {
		jobExecutionTotal = GetCounterVec(meta.Namespace, "cron_job_execution_count", "Number of times cron jobs executed for total/success/failure", meta.JobExecutionTotal.Labels)
	}
	if meta.JobExecutionLatencyMillis != nil {
		jobExecutionLatencyMillis = GetHistogramVec(meta.Namespace, "cron_job_execution_latency_millis", "Tracks the latencies for cron jobs run", meta.JobExecutionLatencyMillis.Labels, meta.JobExecutionLatencyMillis.Buckets)
	}

	return &CronJobMetrics{
		jobExecutionTotal:         jobExecutionTotal,
		jobExecutionLatencyMillis: jobExecutionLatencyMillis,
	}
}

// LogMetricsPre should be called at the start of a cron job execution.
// It increments the total execution counter and returns the start time for latency calculation.
//
// Parameters:
//   - cjMetricsLabelValues: Label values containing the job name.
//
// Returns the start time to be passed to LogMetricsPost for latency calculation.
//
// Example:
//
//	startTime := cronMetrics.LogMetricsPre(&monitoring.CronJobMetricsLabelValues{
//	    JobName: "daily_cleanup",
//	})
//	// ... execute job ...
//	cronMetrics.LogMetricsPost(err, labelValues, startTime)
func (cjm *CronJobMetrics) LogMetricsPre(cjMetricsLabelValues *CronJobMetricsLabelValues) time.Time {
	if cjm.jobExecutionTotal != nil {
		cjm.jobExecutionTotal.WithLabelValues(cjMetricsLabelValues.JobName, Total).Inc()
	}
	return time.Now()
}

// LogMetricsPost should be called after a cron job execution completes.
// It records the success/failure status and the execution latency.
//
// Parameters:
//   - appErr: The error returned by the job (nil for success, non-nil for failure).
//   - cjMetricsLabelValues: Label values containing the job name.
//   - opsExecTime: The start time returned by LogMetricsPre.
//
// Example:
//
//	cronMetrics.LogMetricsPost(nil, labelValues, startTime) // Success
//	cronMetrics.LogMetricsPost(err, labelValues, startTime) // Failure
func (cjm *CronJobMetrics) LogMetricsPost(appErr *ae.AppError, cjMetricsLabelValues *CronJobMetricsLabelValues, opsExecTime time.Time) {
	if cjm.jobExecutionTotal != nil {
		if appErr != nil {
			cjm.jobExecutionTotal.WithLabelValues(cjMetricsLabelValues.JobName, Failure).Inc()
		} else {
			cjm.jobExecutionTotal.WithLabelValues(cjMetricsLabelValues.JobName, Success).Inc()
		}
	}
	if cjm.jobExecutionLatencyMillis != nil {
		cjm.jobExecutionLatencyMillis.WithLabelValues(cjMetricsLabelValues.JobName).Observe(float64(time.Since(opsExecTime).Milliseconds()))
	}
}

// GetJobExecutionTotalMetric returns the underlying Prometheus CounterVec
// for the job execution counter. This can be used for advanced operations.
//
// Returns nil if the metric was not configured during initialization.
func (cjm *CronJobMetrics) GetJobExecutionTotalMetric() *prometheus.CounterVec {
	return cjm.jobExecutionTotal
}

// GetJobExecutionLatencyMillisMetric returns the underlying Prometheus HistogramVec
// for the job execution latency. This can be used for advanced operations.
//
// Returns nil if the metric was not configured during initialization.
func (cjm *CronJobMetrics) GetJobExecutionLatencyMillisMetric() *prometheus.HistogramVec {
	return cjm.jobExecutionLatencyMillis
}
