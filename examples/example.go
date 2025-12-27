// Package main demonstrates how to use the app-monitoring package for collecting
// metrics in a Go application using the Prometheus backend.
//
// This example shows:
//   - Setting up router-level HTTP metrics with Gin middleware
//   - Tracking database operation metrics
//   - Monitoring downstream service calls
//   - Recording cron job execution metrics
//   - Tracking pub/sub messaging metrics
//   - Collecting application-level error metrics
package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	ae "github.com/piyushkumar96/app-error"
	"github.com/piyushkumar96/app-monitoring/interfaces"
	"github.com/piyushkumar96/app-monitoring/models"
	prom "github.com/piyushkumar96/app-monitoring/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Global metric instances - typically initialized once at application startup
// Using interfaces allows for easy mocking in tests and swapping implementations
var (
	routerMetrics     interfaces.RouterMetricsInterface
	dbMetrics         interfaces.DBMetricsInterface
	downstreamMetrics interfaces.DownstreamServiceMetricsInterface
	cronMetrics       interfaces.CronJobMetricsInterface
	pubsubMetrics     interfaces.PSMetricsInterface
	appMetrics        interfaces.AppMetricsInterface
)

func main() {
	// Initialize all metrics
	initializeMetrics()

	// Set up Gin router with metrics middleware
	router := gin.Default()

	// Add router-level metrics middleware
	// This automatically tracks all HTTP requests to your endpoints
	router.Use(routerMetrics.LogMetrics("/metrics"))

	// Expose Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Example API endpoints
	router.GET("/api/users", getUsersHandler)
	router.POST("/api/users", createUserHandler)
	router.GET("/api/orders", getOrdersHandler)

	// Start server
	_ = router.Run(":8080")
}

// initializeMetrics sets up all Prometheus metrics for the application.
// Call this once during application startup.
func initializeMetrics() {
	namespace := "myapp"

	// ============================================
	// Router-Level Metrics (HTTP Endpoint Metrics)
	// ============================================
	// These metrics are automatically collected by the middleware
	routerMetrics = prom.NewPromRouterMetrics(&models.RouterMetricsMeta{
		Namespace: namespace,
		HTTPRequests: &models.MetricMeta{
			Labels: []string{"method", "code", "path", "status"},
		},
		HTTPRequestsLatencyMillis: &models.MetricMeta{
			Labels:  []string{"method", "code", "path"},
			Buckets: prom.GetPromExponentialBuckets(10, 2, 10), // 10ms to ~10s
		},
		HTTPRequestSizeBytes: &models.MetricMeta{
			Labels:  []string{"method", "code", "path"},
			Buckets: prom.GetPromExponentialBuckets(100, 2, 10), // 100B to ~100KB
		},
		HTTPResponseSizeBytes: &models.MetricMeta{
			Labels:  []string{"method", "code", "path"},
			Buckets: prom.GetPromExponentialBuckets(100, 2, 10),
		},
	})

	// ============================================
	// Database Metrics
	// ============================================
	// Track database operations (queries, inserts, updates, deletes)
	dbMetrics = prom.NewPromDatabaseMetrics(&models.DBMetricsMeta{
		Namespace: namespace,
		OperationsTotal: &models.MetricMeta{
			Labels: []string{"op_type", "source", "entity", "is_txn", "status"},
		},
		OperationsLatencyMillis: &models.MetricMeta{
			Labels:  []string{"op_type", "source", "entity", "is_txn"},
			Buckets: prom.GetPromExponentialBuckets(1, 2, 12), // 1ms to ~4s
		},
	})

	// ============================================
	// Downstream Service Metrics
	// ============================================
	// Track HTTP calls to external/downstream services
	downstreamMetrics = prom.NewPromDownstreamServiceMetrics(&models.DownstreamServiceMetricsMeta{
		Namespace: namespace,
		HTTPRequests: &models.MetricMeta{
			Labels: []string{"service", "method", "code", "api", "status"},
		},
		HTTPRequestsLatencyMillis: &models.MetricMeta{
			Labels:  []string{"service", "method", "code", "api"},
			Buckets: prom.GetPromExponentialBuckets(10, 2, 10),
		},
		HTTPRequestSizeBytes: &models.MetricMeta{
			Labels:  []string{"service", "method", "code", "api"},
			Buckets: prom.GetPromExponentialBuckets(100, 2, 10),
		},
		HTTPResponseSizeBytes: &models.MetricMeta{
			Labels:  []string{"service", "method", "code", "api"},
			Buckets: prom.GetPromExponentialBuckets(100, 2, 10),
		},
	})

	// ============================================
	// Cron Job Metrics
	// ============================================
	// Track cron job executions and latencies
	cronMetrics = prom.NewPromCronJobMetrics(&models.CronJobMetricsMeta{
		Namespace: namespace,
		JobExecutionTotal: &models.MetricMeta{
			Labels: []string{"job_name", "status"},
		},
		JobExecutionLatencyMillis: &models.MetricMeta{
			Labels:  []string{"job_name"},
			Buckets: prom.GetPromExponentialBuckets(100, 2, 12), // 100ms to ~400s
		},
	})

	// ============================================
	// Pub/Sub Metrics
	// ============================================
	// Track message publishing and consumption
	pubsubMetrics = prom.NewPromPubSubMetrics(&models.PSMetricsMeta{
		Namespace: namespace,
		TotalMessagesConsumed: &models.MetricMeta{
			Labels: []string{"source", "entity", "op_type", "status", "error_code"},
		},
		TotalMessagesPublished: &models.MetricMeta{
			Labels: []string{"entity", "op_type", "status"},
		},
		MessagesPublishedLatencyMillis: &models.MetricMeta{
			Labels:  []string{"entity", "op_type"},
			Buckets: prom.GetPromExponentialBuckets(10, 2, 10),
		},
		MessagesPublishedSizeBytes: &models.MetricMeta{
			Labels:  []string{"entity", "op_type"},
			Buckets: prom.GetPromExponentialBuckets(100, 2, 12),
		},
	})

	// ============================================
	// Application Error Metrics
	// ============================================
	// Track application-level errors by error code
	appMetrics = prom.NewPromAppMetrics(&models.AppMetricsMeta{
		Namespace: namespace,
		ApplicationErrorsCounter: &models.MetricMeta{
			Labels: []string{"error_code"},
		},
	})
}

// getUsersHandler demonstrates database metrics usage
func getUsersHandler(c *gin.Context) {
	// Start database operation tracking
	labelValues := &models.DBMetricsLabelValues{
		OpType:   "select",
		Source:   "UserHandler",
		AdEntity: "users",
		IsTxn:    "false",
	}
	startTime := dbMetrics.LogMetricsPre(labelValues)

	// Simulate database operation
	users, appErr := fetchUsersFromDB()

	// Record success/failure and latency
	dbMetrics.LogMetricsPost(appErr, labelValues, startTime)

	if appErr != nil {
		// Track application error
		appMetrics.LogMetrics([]string{"ERR_DB_QUERY"})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// createUserHandler demonstrates downstream service metrics usage
func createUserHandler(c *gin.Context) {
	// Track downstream service call to notification service
	labelValues := &models.DownstreamServiceMetricsLabelValues{
		Name:          "notification-service",
		HTTPMethod:    "POST",
		APIIdentifier: "/api/v1/notifications",
	}
	downstreamMetrics.LogMetricsPre(labelValues)

	// Simulate downstream HTTP call
	startTime := time.Now()
	resp, err := callNotificationService()

	// Record downstream call metrics
	httpMetrics := &models.HTTPMetrics{
		Method:                "POST",
		Code:                  resp.StatusCode,
		RequestBodySizeBytes:  1024,
		ResponseBodySizeBytes: 256,
		ResponseTime:          time.Since(startTime),
	}

	success := err == nil && resp.StatusCode >= 200 && resp.StatusCode <= 299
	downstreamMetrics.LogMetricsPost(success, labelValues, httpMetrics)

	if !success {
		appMetrics.LogMetrics([]string{"ERR_NOTIFICATION_SERVICE"})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send notification"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created"})
}

// getOrdersHandler demonstrates pub/sub metrics usage for consumed messages
func getOrdersHandler(c *gin.Context) {
	// This would typically be called in a message consumer
	labelValues := &models.PSMetricsLabelValues{
		Source:       "orders-subscription",
		Entity:       "order",
		EntityOpType: "fetch",
		ErrorCode:    "", // Empty for success
	}

	// Track message consumption start
	_ = pubsubMetrics.LogMetricsPre(labelValues)

	// Process the message...
	// On success, ErrorCode remains empty
	// On failure, set ErrorCode to the appropriate error code

	// Track message consumption completion
	pubsubMetrics.LogMetricsPost(labelValues, nil)

	c.JSON(http.StatusOK, gin.H{"orders": []string{}})
}

// runDailyCleanupJob demonstrates cron job metrics usage
// nolint:unused
func runDailyCleanupJob() {
	labelValues := &models.CronJobMetricsLabelValues{
		JobName: "daily_cleanup",
	}

	// Track job start
	startTime := cronMetrics.LogMetricsPre(labelValues)

	// Execute the job
	appErr := performCleanup()

	// Track job completion with success/failure
	cronMetrics.LogMetricsPost(appErr, labelValues, startTime)
}

// Stub functions for demonstration

// fetchUsersFromDB simulates a database query operation
func fetchUsersFromDB() ([]string, *ae.AppError) {
	// In real implementation, this would query the database
	// Return nil error for success
	return []string{"user1", "user2"}, nil
}

// callNotificationService simulates an HTTP call to an external service
func callNotificationService() (*http.Response, error) {
	// In real implementation, this would make an HTTP request
	return &http.Response{StatusCode: 200}, nil
}

// performCleanup simulates a cleanup job execution
func performCleanup() *ae.AppError {
	// In real implementation, this would perform cleanup tasks
	// Return nil for success, or an AppError for failure
	return nil
}
