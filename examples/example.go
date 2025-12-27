// Package main demonstrates how to use the appMnt package for collecting
// Prometheus metrics in a Go application.
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
	"github.com/gin-gonic/gin"
	ae "github.com/piyushkumar96/app-error"
	appMnt "github.com/piyushkumar96/app-monitoring"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

// Global metric instances - typically initialized once at application startup
var (
	routerMetrics     *appMnt.RouterMetrics
	dbMetrics         *appMnt.DBMetrics
	downstreamMetrics *appMnt.DownstreamServiceMetrics
	cronMetrics       *appMnt.CronJobMetrics
	pubsubMetrics     *appMnt.PSMetrics
	appMetrics        *appMnt.AppMetrics
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
	routerMetrics = appMnt.NewRouterLevelMetrics(&appMnt.RouterMetricsMeta{
		Namespace: namespace,
		HTTPRequests: &appMnt.MetricMeta{
			Labels: []string{"method", "code", "path", "status"},
		},
		HTTPRequestsLatencyMillis: &appMnt.MetricMeta{
			Labels:  []string{"method", "code", "path"},
			Buckets: appMnt.GetExponentialBuckets(10, 2, 10), // 10ms to ~10s
		},
		HTTPRequestSizeBytes: &appMnt.MetricMeta{
			Labels:  []string{"method", "code", "path"},
			Buckets: appMnt.GetExponentialBuckets(100, 2, 10), // 100B to ~100KB
		},
		HTTPResponseSizeBytes: &appMnt.MetricMeta{
			Labels:  []string{"method", "code", "path"},
			Buckets: appMnt.GetExponentialBuckets(100, 2, 10),
		},
	})

	// ============================================
	// Database Metrics
	// ============================================
	// Track database operations (queries, inserts, updates, deletes)
	dbMetrics = appMnt.NewDatabaseMetrics(&appMnt.DBMetricsMeta{
		Namespace: namespace,
		OperationsTotal: &appMnt.MetricMeta{
			Labels: []string{"op_type", "source", "entity", "is_txn", "status"},
		},
		OperationsLatencyMillis: &appMnt.MetricMeta{
			Labels:  []string{"op_type", "source", "entity", "is_txn"},
			Buckets: appMnt.GetExponentialBuckets(1, 2, 12), // 1ms to ~4s
		},
	})

	// ============================================
	// Downstream Service Metrics
	// ============================================
	// Track HTTP calls to external/downstream services
	downstreamMetrics = appMnt.NewDownstreamServiceMetrics(&appMnt.DownstreamServiceMetricsMeta{
		Namespace: namespace,
		HTTPRequests: &appMnt.MetricMeta{
			Labels: []string{"service", "method", "code", "api", "status"},
		},
		HTTPRequestsLatencyMillis: &appMnt.MetricMeta{
			Labels:  []string{"service", "method", "code", "api"},
			Buckets: appMnt.GetExponentialBuckets(10, 2, 10),
		},
		HTTPRequestSizeBytes: &appMnt.MetricMeta{
			Labels:  []string{"service", "method", "code", "api"},
			Buckets: appMnt.GetExponentialBuckets(100, 2, 10),
		},
		HTTPResponseSizeBytes: &appMnt.MetricMeta{
			Labels:  []string{"service", "method", "code", "api"},
			Buckets: appMnt.GetExponentialBuckets(100, 2, 10),
		},
	})

	// ============================================
	// Cron Job Metrics
	// ============================================
	// Track cron job executions and latencies
	cronMetrics = appMnt.NewCronJobMetrics(&appMnt.CronJobMetricsMeta{
		Namespace: namespace,
		JobExecutionTotal: &appMnt.MetricMeta{
			Labels: []string{"job_name", "status"},
		},
		JobExecutionLatencyMillis: &appMnt.MetricMeta{
			Labels:  []string{"job_name"},
			Buckets: appMnt.GetExponentialBuckets(100, 2, 12), // 100ms to ~400s
		},
	})

	// ============================================
	// Pub/Sub Metrics
	// ============================================
	// Track message publishing and consumption
	pubsubMetrics = appMnt.NewPubSubMetrics(&appMnt.PSMetricsMeta{
		Namespace: namespace,
		TotalMessagesConsumed: &appMnt.MetricMeta{
			Labels: []string{"source", "entity", "op_type", "status", "error_code"},
		},
		TotalMessagesPublished: &appMnt.MetricMeta{
			Labels: []string{"entity", "op_type", "status"},
		},
		MessagesPublishedLatencyMillis: &appMnt.MetricMeta{
			Labels:  []string{"entity", "op_type"},
			Buckets: appMnt.GetExponentialBuckets(10, 2, 10),
		},
		MessagesPublishedSizeBytes: &appMnt.MetricMeta{
			Labels:  []string{"entity", "op_type"},
			Buckets: appMnt.GetExponentialBuckets(100, 2, 12),
		},
	})

	// ============================================
	// Application Error Metrics
	// ============================================
	// Track application-level errors by error code
	appMetrics = appMnt.NewAppMetrics(&appMnt.AppMetricsMeta{
		Namespace: namespace,
		ApplicationErrorsCounter: &appMnt.MetricMeta{
			Labels: []string{"error_code"},
		},
	})
}

// getUsersHandler demonstrates database metrics usage
func getUsersHandler(c *gin.Context) {
	// Start database operation tracking
	labelValues := &appMnt.DBMetricsLabelValues{
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
	labelValues := &appMnt.DownstreamServiceMetricsLabelValues{
		Name:          "notification-service",
		HTTPMethod:    "POST",
		APIIdentifier: "/api/v1/notifications",
	}
	downstreamMetrics.LogMetricsPre(labelValues)

	// Simulate downstream HTTP call
	startTime := time.Now()
	resp, err := callNotificationService()

	// Record downstream call metrics
	httpMetrics := &appMnt.HTTPMetrics{
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
	labelValues := &appMnt.PSMetricsLabelValues{
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
	labelValues := &appMnt.CronJobMetricsLabelValues{
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
