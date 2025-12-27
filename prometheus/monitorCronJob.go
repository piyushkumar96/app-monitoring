package prometheus

import (
	"time"

	"github.com/piyushkumar96/app-monitoring/constants"
	"github.com/piyushkumar96/app-monitoring/interfaces"
	"github.com/piyushkumar96/app-monitoring/models"

	ae "github.com/piyushkumar96/app-error"
	"github.com/prometheus/client_golang/prometheus"
)

// NewPromCronJobMetrics creates and registers Prometheus cron job execution metrics.
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
// Returns an interfaces.CronJobMetricsInterface instance that can be used to log job execution metrics.
func NewPromCronJobMetrics(meta *models.CronJobMetricsMeta) interfaces.CronJobMetricsInterface {
	var jobExecutionTotal *prometheus.CounterVec
	var jobExecutionLatencyMillis *prometheus.HistogramVec

	if meta.JobExecutionTotal != nil {
		jobExecutionTotal = GetPromCounterVec(meta.Namespace, "cron_job_execution_count", "Number of times cron jobs executed for total/success/failure", meta.JobExecutionTotal.Labels)
	}
	if meta.JobExecutionLatencyMillis != nil {
		jobExecutionLatencyMillis = GetPromHistogramVec(meta.Namespace, "cron_job_execution_latency_millis", "Tracks the latencies for cron jobs run", meta.JobExecutionLatencyMillis.Labels, meta.JobExecutionLatencyMillis.Buckets)
	}

	return &PromCronJobMetrics{
		jobExecutionTotal:         jobExecutionTotal,
		jobExecutionLatencyMillis: jobExecutionLatencyMillis,
	}
}

// LogMetricsPre should be called at the start of a cron job execution.
// It increments the total execution counter and returns the start time for latency calculation.
func (cjm *PromCronJobMetrics) LogMetricsPre(cjMetricsLabelValues *models.CronJobMetricsLabelValues) time.Time {
	if cjm.jobExecutionTotal != nil {
		cjm.jobExecutionTotal.WithLabelValues(cjMetricsLabelValues.JobName, constants.Total).Inc()
	}
	return time.Now()
}

// LogMetricsPost should be called after a cron job execution completes.
// It records the success/failure status and the execution latency.
func (cjm *PromCronJobMetrics) LogMetricsPost(appErr *ae.AppError, cjMetricsLabelValues *models.CronJobMetricsLabelValues, opsExecTime time.Time) {
	if cjm.jobExecutionTotal != nil {
		if appErr != nil {
			cjm.jobExecutionTotal.WithLabelValues(cjMetricsLabelValues.JobName, constants.Failure).Inc()
		} else {
			cjm.jobExecutionTotal.WithLabelValues(cjMetricsLabelValues.JobName, constants.Success).Inc()
		}
	}
	if cjm.jobExecutionLatencyMillis != nil {
		cjm.jobExecutionLatencyMillis.WithLabelValues(cjMetricsLabelValues.JobName).Observe(float64(time.Since(opsExecTime).Milliseconds()))
	}
}

// GetJobExecutionTotalMetric returns the underlying Prometheus CounterVec
// for the job execution counter. This can be used for advanced operations.
func (cjm *PromCronJobMetrics) GetJobExecutionTotalMetric() *prometheus.CounterVec {
	return cjm.jobExecutionTotal
}

// GetJobExecutionLatencyMillisMetric returns the underlying Prometheus HistogramVec
// for the job execution latency. This can be used for advanced operations.
func (cjm *PromCronJobMetrics) GetJobExecutionLatencyMillisMetric() *prometheus.HistogramVec {
	return cjm.jobExecutionLatencyMillis
}
