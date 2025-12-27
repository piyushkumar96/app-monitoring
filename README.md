# App Monitoring

A Go package for capturing application metrics with pluggable backends. Currently supports Prometheus, with an extensible architecture for adding other backends (OpenTelemetry, StatsD, etc.) in the future.

## Features

This package provides easy-to-use abstractions for common metric types:

| Metric Type | Description | Use Case |
|-------------|-------------|----------|
| **Router** | HTTP endpoint metrics | Track requests, latencies, and payload sizes at the application level |
| **Database** | DB operation metrics | Monitor query counts and latencies by operation type |
| **Downstream Service** | External HTTP call metrics | Track outbound HTTP requests to other services |
| **Pub/Sub** | Messaging metrics | Monitor message publishing and consumption |
| **Cron Job** | Scheduled job metrics | Track job executions and durations |
| **Application** | Error tracking | Count application-level errors by error code |

## Installation

```bash
go get github.com/piyushkumar96/app-monitoring@latest
```

## Project Structure

```
app-monitoring/
├── constants/            # Shared constants package
│   └── constants.go      # Total, Success, Failure, HTTP status constants
├── interfaces/           # Generic interfaces package
│   ├── interfaces.go     # Interface definitions for all metric types
│   └── mock.go           # Mock implementations for testing
├── models/               # Shared data models package
│   └── model.go          # Configuration types and label value types
├── prometheus/           # Prometheus-specific implementation
│   ├── metric.go
│   ├── model.go
│   ├── monitorApp.go
│   ├── monitorCronJob.go
│   ├── monitorDatabase.go
│   ├── monitorDownstreamService.go
│   ├── monitorPubSub.go
│   ├── monitorRouter.go
│   └── noop.go           # NoOp implementations for testing
├── examples/
│   └── example.go
├── go.mod
└── README.md
```

## Quick Start

### 1. Initialize Prometheus Metrics

```go
package main

import (
    "github.com/piyushkumar96/app-monitoring/interfaces"
    "github.com/piyushkumar96/app-monitoring/models"
    prom "github.com/piyushkumar96/app-monitoring/prometheus"
    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
    // Initialize router metrics using Prometheus backend
    routerMetrics := prom.NewPromRouterMetrics(&models.RouterMetricsMeta{
        Namespace: "myapp",
        HTTPRequests: &models.MetricMeta{
            Labels: []string{"method", "code", "path", "status"},
        },
        HTTPRequestsLatencyMillis: &models.MetricMeta{
            Labels:  []string{"method", "code", "path"},
            Buckets: prom.GetPromExponentialBuckets(10, 2, 10),
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
dbMetrics := prom.NewPromDatabaseMetrics(&models.DBMetricsMeta{
    Namespace: "myapp",
    OperationsTotal: &models.MetricMeta{
        Labels: []string{"op_type", "source", "entity", "is_txn", "status"},
    },
    OperationsLatencyMillis: &models.MetricMeta{
        Labels:  []string{"op_type", "source", "entity", "is_txn"},
        Buckets: prom.GetPromExponentialBuckets(1, 2, 12),
    },
})

// In your repository/handler
func GetUser(id string) (*User, error) {
    labelValues := &models.DBMetricsLabelValues{
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
dsMetrics := prom.NewPromDownstreamServiceMetrics(&models.DownstreamServiceMetricsMeta{
    Namespace: "myapp",
    HTTPRequests: &models.MetricMeta{
        Labels: []string{"service", "method", "code", "api", "status"},
    },
    HTTPRequestsLatencyMillis: &models.MetricMeta{
        Labels:  []string{"service", "method", "code", "api"},
        Buckets: prom.GetPromExponentialBuckets(10, 2, 10),
    },
})

// Before making HTTP call
labelValues := &models.DownstreamServiceMetricsLabelValues{
    Name:          "payment-service",
    HTTPMethod:    "POST",
    APIIdentifier: "/api/v1/payments",
}
dsMetrics.LogMetricsPre(labelValues)

// Make HTTP call
startTime := time.Now()
resp, err := http.Post(url, "application/json", body)

// After call completes
httpMetrics := &models.HTTPMetrics{
    Method:       "POST",
    Code:         resp.StatusCode,
    ResponseTime: time.Since(startTime),
}
dsMetrics.LogMetricsPost(resp.StatusCode >= 200 && resp.StatusCode <= 299, labelValues, httpMetrics)
```

### 4. Track Cron Job Executions

```go
cronMetrics := prom.NewPromCronJobMetrics(&models.CronJobMetricsMeta{
    Namespace: "myapp",
    JobExecutionTotal: &models.MetricMeta{
        Labels: []string{"job_name", "status"},
    },
    JobExecutionLatencyMillis: &models.MetricMeta{
        Labels:  []string{"job_name"},
        Buckets: prom.GetPromExponentialBuckets(100, 2, 12),
    },
})

// In your cron job
func RunCleanupJob() {
    labelValues := &models.CronJobMetricsLabelValues{
        JobName: "daily_cleanup",
    }
    
    startTime := cronMetrics.LogMetricsPre(labelValues)
    err := performCleanup()
    cronMetrics.LogMetricsPost(err, labelValues, startTime)
}
```

### 5. Track Pub/Sub Operations

```go
psMetrics := prom.NewPromPubSubMetrics(&models.PSMetricsMeta{
    Namespace: "myapp",
    TotalMessagesConsumed: &models.MetricMeta{
        Labels: []string{"source", "entity", "op_type", "status", "error_code"},
    },
    TotalMessagesPublished: &models.MetricMeta{
        Labels: []string{"entity", "op_type", "status"},
    },
})

// For message consumption
labelValues := &models.PSMetricsLabelValues{
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
appMetrics := prom.NewPromAppMetrics(&models.AppMetricsMeta{
    Namespace: "myapp",
    ApplicationErrorsCounter: &models.MetricMeta{
        Labels: []string{"error_code"},
    },
})

// When errors occur
appMetrics.LogMetrics([]string{"ERR_DB_CONNECTION", "ERR_VALIDATION"})

// When error is resolved
appMetrics.DecrementAppErrorCount("ERR_DB_CONNECTION")
```

## Interface-Based Architecture

All metric types are defined as generic interfaces in the `interfaces` package, enabling:
- **Easy mocking** for unit tests
- **Dependency injection** in your application
- **Future extensibility** for other metric backends (OpenTelemetry, StatsD, etc.)

### Using Interfaces

```go
import (
    "github.com/piyushkumar96/app-monitoring/interfaces"
    "github.com/piyushkumar96/app-monitoring/models"
    prom "github.com/piyushkumar96/app-monitoring/prometheus"
)

// Declare metrics using generic interfaces
var (
    routerMetrics     interfaces.RouterMetricsInterface
    dbMetrics         interfaces.DBMetricsInterface
    downstreamMetrics interfaces.DownstreamServiceMetricsInterface
    cronMetrics       interfaces.CronJobMetricsInterface
    pubsubMetrics     interfaces.PSMetricsInterface
    appMetrics        interfaces.AppMetricsInterface
)

// Initialize with Prometheus implementations
routerMetrics = prom.NewPromRouterMetrics(&models.RouterMetricsMeta{...})

// Or use NoOp implementations for testing
routerMetrics = prom.NewNoOpPromRouterMetrics()

// Or use Mock implementations for unit testing with assertions
routerMetrics = interfaces.NewMockRouterMetrics()
```

### Available Interfaces & Implementations

| Interface (interfaces pkg) | Prometheus Constructor | NoOp Constructor | Mock Constructor |
|---------------------------|------------------------|------------------|------------------|
| `RouterMetricsInterface` | `prom.NewPromRouterMetrics()` | `prom.NewNoOpPromRouterMetrics()` | `interfaces.NewMockRouterMetrics()` |
| `DBMetricsInterface` | `prom.NewPromDatabaseMetrics()` | `prom.NewNoOpPromDBMetrics()` | `interfaces.NewMockDBMetrics()` |
| `DownstreamServiceMetricsInterface` | `prom.NewPromDownstreamServiceMetrics()` | `prom.NewNoOpPromDownstreamServiceMetrics()` | `interfaces.NewMockDownstreamServiceMetrics()` |
| `CronJobMetricsInterface` | `prom.NewPromCronJobMetrics()` | `prom.NewNoOpPromCronJobMetrics()` | `interfaces.NewMockCronJobMetrics()` |
| `PSMetricsInterface` | `prom.NewPromPubSubMetrics()` | `prom.NewNoOpPromPSMetrics()` | `interfaces.NewMockPSMetrics()` |
| `AppMetricsInterface` | `prom.NewPromAppMetrics()` | `prom.NewNoOpPromAppMetrics()` | `interfaces.NewMockAppMetrics()` |

### Testing with Mock Implementations

The `interfaces` package provides mock implementations that track method calls for assertions:

```go
func TestUserHandler(t *testing.T) {
    // Create mock metrics for testing
    mockDBMetrics := interfaces.NewMockDBMetrics()
    mockAppMetrics := interfaces.NewMockAppMetrics()
    
    handler := NewUserHandler(mockDBMetrics, mockAppMetrics)
    
    // Call the handler
    handler.GetUser("123")
    
    // Assert metrics were logged
    if !mockDBMetrics.LogMetricsPreCalled {
        t.Error("Expected LogMetricsPre to be called")
    }
    if mockDBMetrics.LogMetricsPreLabelValues.OpType != "select" {
        t.Error("Expected OpType to be 'select'")
    }
}
```

### Testing with NoOp Implementations

For simpler tests where you don't need to assert on metrics:

```go
func TestUserHandler(t *testing.T) {
    // Use NoOp metrics - no actual metrics recorded, no assertions
    dbMetrics := prom.NewNoOpPromDBMetrics()
    appMetrics := prom.NewNoOpPromAppMetrics()
    
    handler := NewUserHandler(dbMetrics, appMetrics)
    // ... test your handler
}
```

## Configuration Options

### Metric Labels

Each metric type supports customizable labels. The labels you specify in `MetricMeta.Labels` must match the order of label values you provide when logging metrics.

### Histogram Buckets

Use `prom.GetPromExponentialBuckets(start, factor, count)` to generate exponential bucket boundaries:

```go
// Generate buckets: 10, 20, 40, 80, 160, 320, 640, 1280, 2560, 5120
buckets := prom.GetPromExponentialBuckets(10, 2, 10)
```

### Disabling Metrics

Set any metric configuration to `nil` to disable it:

```go
routerMetrics := prom.NewPromRouterMetrics(&models.RouterMetricsMeta{
    Namespace:                 "myapp",
    HTTPRequests:              &models.MetricMeta{Labels: []string{"method", "code", "path", "status"}},
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
6. **Use interfaces** - Declare variables using interfaces for better testability
7. **Use mocks for assertions** - Use mock implementations when you need to verify metrics are logged correctly

## License

See [LICENSE](LICENSE) for details.
