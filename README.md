# Monitoring

A Go package for capturing Prometheus metrics to monitor application state. These metrics can be visualized in Grafana dashboards for comprehensive application observability.

## Features

This package provides easy-to-use abstractions for common Prometheus metric types:

| Metric Type | Description | Use Case |
|-------------|-------------|----------|
| **Router** | HTTP endpoint metrics | Track requests, latencies, and payload sizes at the application level |
| **Database** | DB operation metrics | Monitor query counts and latencies by operation type |
| **Downstream Service** | External HTTP call metrics | Track outbound HTTP requests to other services |
| **Pub/Sub** | Messaging metrics | Monitor message publishing and consumption |
| **Cron Job** | Scheduled job metrics | Track job executions and durations |
| **Application** | Error tracking | Count application-level errors by error code |

## Supported Metrics

### Router Metrics
- `http_requests` - Counter for HTTP requests (total/success/failure)
- `http_request_latency_millis` - Histogram for request latencies
- `http_request_size_bytes` - Histogram for request body sizes
- `http_response_size_bytes` - Histogram for response body sizes

### Database Metrics
- `db_operations` - Counter for DB operations (total/success/failure)
- `db_operations_latency_millis` - Histogram for operation latencies

### Downstream Service Metrics
- `downstream_service_http_requests` - Counter for downstream HTTP requests
- `downstream_service_http_request_latency_millis` - Histogram for request latencies
- `downstream_service_http_request_size_bytes` - Histogram for request sizes
- `downstream_service_http_response_size_bytes` - Histogram for response sizes

### Pub/Sub Metrics
- `pubsub_messages_consumed` - Counter for consumed messages
- `pubsub_messages_published` - Counter for published messages
- `pubsub_messages_published_latency_millis` - Histogram for publish latencies
- `pubsub_messages_published_size_bytes` - Histogram for message sizes

### Cron Job Metrics
- `cron_job_execution_count` - Counter for job executions (total/success/failure)
- `cron_job_execution_latency_millis` - Histogram for job durations

### Application Metrics
- `application_errors_total` - Gauge for application errors by error code

## Installation

```bash
go get github.com/piyushkumar96/app-monitoring@latest
```

## Quick Start

### 1. Initialize Metrics

```go
package main

import (
    monitoring "github.com/piyushkumar96/app-monitoring"
    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
    // Initialize router metrics
    routerMetrics := monitoring.NewRouterLevelMetrics(&monitoring.RouterMetricsMeta{
        Namespace: "myapp",
        HTTPRequests: &monitoring.MetricMeta{
            Labels: []string{"method", "code", "path", "status"},
        },
        HTTPRequestsLatencyMillis: &monitoring.MetricMeta{
            Labels:  []string{"method", "code", "path"},
            Buckets: monitoring.GetExponentialBuckets(10, 2, 10),
        },
    })

    // Set up Gin router with metrics middleware
    router := gin.Default()
    router.Use(routerMetrics.LogMetrics("/metrics"))
    router.GET("/metrics", gin.WrapH(promhttp.Handler()))
    
    router.Run(":8080")
}
```

### 2. Track Database Operations

```go
dbMetrics := monitoring.NewDatabaseMetrics(&monitoring.DBMetricsMeta{
    Namespace: "myapp",
    OperationsTotal: &monitoring.MetricMeta{
        Labels: []string{"op_type", "source", "entity", "is_txn", "status"},
    },
    OperationsLatencyMillis: &monitoring.MetricMeta{
        Labels:  []string{"op_type", "source", "entity", "is_txn"},
        Buckets: monitoring.GetExponentialBuckets(1, 2, 12),
    },
})

// In your repository/handler
func GetUser(id string) (*User, error) {
    labelValues := &monitoring.DBMetricsLabelValues{
        OpType:   "select",
        Source:   "UserRepository",
        AdEntity: "users",
        IsTxn:    "false",
    }
    
    startTime := dbMetrics.LogMetricsPre(labelValues)
    user, err := db.Query("SELECT * FROM users WHERE id = ?", id)
    dbMetrics.LogMetricsPost(err, labelValues, startTime)
    
    return user, err
}
```

### 3. Track Downstream Service Calls

```go
dsMetrics := monitoring.NewDownstreamServiceMetrics(&monitoring.DownstreamServiceMetricsMeta{
    Namespace: "myapp",
    HTTPRequests: &monitoring.MetricMeta{
        Labels: []string{"service", "method", "code", "api", "status"},
    },
    HTTPRequestsLatencyMillis: &monitoring.MetricMeta{
        Labels:  []string{"service", "method", "code", "api"},
        Buckets: monitoring.GetExponentialBuckets(10, 2, 10),
    },
})

// Before making HTTP call
labelValues := &monitoring.DownstreamServiceMetricsLabelValues{
    Name:          "payment-service",
    HTTPMethod:    "POST",
    APIIdentifier: "/api/v1/payments",
}
dsMetrics.LogMetricsPre(labelValues)

// Make HTTP call
startTime := time.Now()
resp, err := http.Post(url, "application/json", body)

// After call completes
httpMetrics := &monitoring.HTTPMetrics{
    Method:       "POST",
    Code:         resp.StatusCode,
    ResponseTime: time.Since(startTime),
}
dsMetrics.LogMetricsPost(resp.StatusCode >= 200 && resp.StatusCode <= 299, labelValues, httpMetrics)
```

### 4. Track Cron Job Executions

```go
cronMetrics := monitoring.NewCronJobMetrics(&monitoring.CronJobMetricsMeta{
    Namespace: "myapp",
    JobExecutionTotal: &monitoring.MetricMeta{
        Labels: []string{"job_name", "status"},
    },
    JobExecutionLatencyMillis: &monitoring.MetricMeta{
        Labels:  []string{"job_name"},
        Buckets: monitoring.GetExponentialBuckets(100, 2, 12),
    },
})

// In your cron job
func RunCleanupJob() {
    labelValues := &monitoring.CronJobMetricsLabelValues{
        JobName: "daily_cleanup",
    }
    
    startTime := cronMetrics.LogMetricsPre(labelValues)
    err := performCleanup()
    cronMetrics.LogMetricsPost(err, labelValues, startTime)
}
```

### 5. Track Pub/Sub Operations

```go
psMetrics := monitoring.NewPubSubMetrics(&monitoring.PSMetricsMeta{
    Namespace: "myapp",
    TotalMessagesConsumed: &monitoring.MetricMeta{
        Labels: []string{"source", "entity", "op_type", "status", "error_code"},
    },
    TotalMessagesPublished: &monitoring.MetricMeta{
        Labels: []string{"entity", "op_type", "status"},
    },
})

// For message consumption
labelValues := &monitoring.PSMetricsLabelValues{
    Source:       "orders-subscription",
    Entity:       "order",
    EntityOpType: "create",
    ErrorCode:    "", // Set error code on failure
}

startTime := psMetrics.LogMetricsPre(labelValues)
err := processMessage(msg)
if err != nil {
    labelValues.ErrorCode = "ERR_PROCESSING"
}
psMetrics.LogMetricsPost(labelValues, nil)
```

### 6. Track Application Errors

```go
appMetrics := monitoring.NewAppMetrics(&monitoring.AppMetricsMeta{
    Namespace: "myapp",
    ApplicationErrorsCounter: &monitoring.MetricMeta{
        Labels: []string{"error_code"},
    },
})

// When errors occur
appMetrics.LogMetrics([]string{"ERR_DB_CONNECTION", "ERR_VALIDATION"})

// When error is resolved
appMetrics.DecrementAppErrorCount("ERR_DB_CONNECTION")
```

## Configuration Options

### Metric Labels

Each metric type supports customizable labels. The labels you specify in `MetricMeta.Labels` must match the order of label values you provide when logging metrics.

### Histogram Buckets

Use `monitoring.GetExponentialBuckets(start, factor, count)` to generate exponential bucket boundaries:

```go
// Generate buckets: 10, 20, 40, 80, 160, 320, 640, 1280, 2560, 5120
buckets := monitoring.GetExponentialBuckets(10, 2, 10)
```

### Disabling Metrics

Set any metric configuration to `nil` to disable it:

```go
routerMetrics := monitoring.NewRouterLevelMetrics(&monitoring.RouterMetricsMeta{
    Namespace:                 "myapp",
    HTTPRequests:              &monitoring.MetricMeta{Labels: []string{"method", "code", "path", "status"}},
    HTTPRequestsLatencyMillis: nil, // Disabled
    HTTPRequestSizeBytes:      nil, // Disabled
    HTTPResponseSizeBytes:     nil, // Disabled
})
```

## Complete Example

See [examples/example.go](examples/example.go) for a complete working example demonstrating all metric types.

## Best Practices

1. **Initialize metrics once** - Create metric instances during application startup
2. **Use consistent namespaces** - Use your application name as the namespace
3. **Choose appropriate buckets** - Match bucket ranges to your expected latency distributions
4. **Limit cardinality** - Avoid high-cardinality labels (e.g., user IDs, request IDs)
5. **Handle nil metrics** - The package handles nil metric configs gracefully

## License

See [LICENSE](LICENSE) for details.
