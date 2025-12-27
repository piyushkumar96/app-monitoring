package interfaces

import (
	"time"

	"github.com/gin-gonic/gin"
	ae "github.com/piyushkumar96/app-error"
	"github.com/piyushkumar96/app-monitoring/models"
	pubsub "github.com/piyushkumar96/generic-pubsub"
)

// MockRouterMetrics is a mock implementation of RouterMetricsInterface for testing.
type MockRouterMetrics struct {
	// LogMetricsCalled tracks if LogMetrics was called.
	LogMetricsCalled bool
	// LogMetricsPath stores the metricsPath argument.
	LogMetricsPath string
}

// NewMockRouterMetrics creates a new mock router metrics instance.
func NewMockRouterMetrics() *MockRouterMetrics {
	return &MockRouterMetrics{}
}

// LogMetrics returns a pass-through middleware and records the call.
func (m *MockRouterMetrics) LogMetrics(metricsPath string) gin.HandlerFunc {
	m.LogMetricsCalled = true
	m.LogMetricsPath = metricsPath
	return func(c *gin.Context) {
		c.Next()
	}
}

// MockDBMetrics is a mock implementation of DBMetricsInterface for testing.
type MockDBMetrics struct {
	// LogMetricsPreCalled tracks if LogMetricsPre was called.
	LogMetricsPreCalled bool
	// LogMetricsPreLabelValues stores the label values from LogMetricsPre.
	LogMetricsPreLabelValues *models.DBMetricsLabelValues

	// LogMetricsPostCalled tracks if LogMetricsPost was called.
	LogMetricsPostCalled bool
	// LogMetricsPostAppErr stores the appErr from LogMetricsPost.
	LogMetricsPostAppErr *ae.AppError
	// LogMetricsPostLabelValues stores the label values from LogMetricsPost.
	LogMetricsPostLabelValues *models.DBMetricsLabelValues
}

// NewMockDBMetrics creates a new mock database metrics instance.
func NewMockDBMetrics() *MockDBMetrics {
	return &MockDBMetrics{}
}

// LogMetricsPre records the call and returns the current time.
func (m *MockDBMetrics) LogMetricsPre(dbMetricsLabelValues *models.DBMetricsLabelValues) time.Time {
	m.LogMetricsPreCalled = true
	m.LogMetricsPreLabelValues = dbMetricsLabelValues
	return time.Now()
}

// LogMetricsPost records the call.
func (m *MockDBMetrics) LogMetricsPost(appErr *ae.AppError, dbMetricsLabelValues *models.DBMetricsLabelValues, _ time.Time) {
	m.LogMetricsPostCalled = true
	m.LogMetricsPostAppErr = appErr
	m.LogMetricsPostLabelValues = dbMetricsLabelValues
}

// MockDownstreamServiceMetrics is a mock implementation of DownstreamServiceMetricsInterface for testing.
type MockDownstreamServiceMetrics struct {
	// LogMetricsPreCalled tracks if LogMetricsPre was called.
	LogMetricsPreCalled bool
	// LogMetricsPreLabelValues stores the label values from LogMetricsPre.
	LogMetricsPreLabelValues *models.DownstreamServiceMetricsLabelValues

	// LogMetricsPostCalled tracks if LogMetricsPost was called.
	LogMetricsPostCalled bool
	// LogMetricsPostSuccess stores the success flag from LogMetricsPost.
	LogMetricsPostSuccess bool
	// LogMetricsPostLabelValues stores the label values from LogMetricsPost.
	LogMetricsPostLabelValues *models.DownstreamServiceMetricsLabelValues
	// LogMetricsPostHTTPMetrics stores the HTTP metrics from LogMetricsPost.
	LogMetricsPostHTTPMetrics *models.HTTPMetrics
}

// NewMockDownstreamServiceMetrics creates a new mock downstream service metrics instance.
func NewMockDownstreamServiceMetrics() *MockDownstreamServiceMetrics {
	return &MockDownstreamServiceMetrics{}
}

// LogMetricsPre records the call.
func (m *MockDownstreamServiceMetrics) LogMetricsPre(dssMetricsLabelValues *models.DownstreamServiceMetricsLabelValues) {
	m.LogMetricsPreCalled = true
	m.LogMetricsPreLabelValues = dssMetricsLabelValues
}

// LogMetricsPost records the call.
func (m *MockDownstreamServiceMetrics) LogMetricsPost(success bool, dssMetricsLabelValues *models.DownstreamServiceMetricsLabelValues, httpMetrics *models.HTTPMetrics) {
	m.LogMetricsPostCalled = true
	m.LogMetricsPostSuccess = success
	m.LogMetricsPostLabelValues = dssMetricsLabelValues
	m.LogMetricsPostHTTPMetrics = httpMetrics
}

// MockCronJobMetrics is a mock implementation of CronJobMetricsInterface for testing.
type MockCronJobMetrics struct {
	// LogMetricsPreCalled tracks if LogMetricsPre was called.
	LogMetricsPreCalled bool
	// LogMetricsPreLabelValues stores the label values from LogMetricsPre.
	LogMetricsPreLabelValues *models.CronJobMetricsLabelValues

	// LogMetricsPostCalled tracks if LogMetricsPost was called.
	LogMetricsPostCalled bool
	// LogMetricsPostAppErr stores the appErr from LogMetricsPost.
	LogMetricsPostAppErr *ae.AppError
	// LogMetricsPostLabelValues stores the label values from LogMetricsPost.
	LogMetricsPostLabelValues *models.CronJobMetricsLabelValues
}

// NewMockCronJobMetrics creates a new mock cron job metrics instance.
func NewMockCronJobMetrics() *MockCronJobMetrics {
	return &MockCronJobMetrics{}
}

// LogMetricsPre records the call and returns the current time.
func (m *MockCronJobMetrics) LogMetricsPre(cjMetricsLabelValues *models.CronJobMetricsLabelValues) time.Time {
	m.LogMetricsPreCalled = true
	m.LogMetricsPreLabelValues = cjMetricsLabelValues
	return time.Now()
}

// LogMetricsPost records the call.
func (m *MockCronJobMetrics) LogMetricsPost(appErr *ae.AppError, cjMetricsLabelValues *models.CronJobMetricsLabelValues, _ time.Time) {
	m.LogMetricsPostCalled = true
	m.LogMetricsPostAppErr = appErr
	m.LogMetricsPostLabelValues = cjMetricsLabelValues
}

// MockPSMetrics is a mock implementation of PSMetricsInterface for testing.
type MockPSMetrics struct {
	// LogMetricsPreCalled tracks if LogMetricsPre was called.
	LogMetricsPreCalled bool
	// LogMetricsPreLabelValues stores the label values from LogMetricsPre.
	LogMetricsPreLabelValues *models.PSMetricsLabelValues

	// LogMetricsPostCalled tracks if LogMetricsPost was called.
	LogMetricsPostCalled bool
	// LogMetricsPostLabelValues stores the label values from LogMetricsPost.
	LogMetricsPostLabelValues *models.PSMetricsLabelValues
	// LogMetricsPostEventTxnData stores the event txn data from LogMetricsPost.
	LogMetricsPostEventTxnData *pubsub.EventTxnData
}

// NewMockPSMetrics creates a new mock pub/sub metrics instance.
func NewMockPSMetrics() *MockPSMetrics {
	return &MockPSMetrics{}
}

// LogMetricsPre records the call and returns the current time.
func (m *MockPSMetrics) LogMetricsPre(psMetricsLabelValues *models.PSMetricsLabelValues) time.Time {
	m.LogMetricsPreCalled = true
	m.LogMetricsPreLabelValues = psMetricsLabelValues
	return time.Now()
}

// LogMetricsPost records the call.
func (m *MockPSMetrics) LogMetricsPost(psMetricsLabelValues *models.PSMetricsLabelValues, eventTxnData *pubsub.EventTxnData) {
	m.LogMetricsPostCalled = true
	m.LogMetricsPostLabelValues = psMetricsLabelValues
	m.LogMetricsPostEventTxnData = eventTxnData
}

// MockAppMetrics is a mock implementation of AppMetricsInterface for testing.
type MockAppMetrics struct {
	// LogMetricsCalled tracks if LogMetrics was called.
	LogMetricsCalled bool
	// LogMetricsErrCodes stores the error codes from LogMetrics.
	LogMetricsErrCodes []string

	// DecrementAppErrorCountCalled tracks if DecrementAppErrorCount was called.
	DecrementAppErrorCountCalled bool
	// DecrementAppErrorCountErrCode stores the error code from DecrementAppErrorCount.
	DecrementAppErrorCountErrCode string
}

// NewMockAppMetrics creates a new mock application metrics instance.
func NewMockAppMetrics() *MockAppMetrics {
	return &MockAppMetrics{}
}

// LogMetrics records the call.
func (m *MockAppMetrics) LogMetrics(errCodes []string) {
	m.LogMetricsCalled = true
	m.LogMetricsErrCodes = errCodes
}

// DecrementAppErrorCount records the call.
func (m *MockAppMetrics) DecrementAppErrorCount(errCode string) {
	m.DecrementAppErrorCountCalled = true
	m.DecrementAppErrorCountErrCode = errCode
}

// Compile-time interface implementation checks for Mock types
var (
	_ RouterMetricsInterface            = (*MockRouterMetrics)(nil)
	_ DBMetricsInterface                = (*MockDBMetrics)(nil)
	_ DownstreamServiceMetricsInterface = (*MockDownstreamServiceMetrics)(nil)
	_ CronJobMetricsInterface           = (*MockCronJobMetrics)(nil)
	_ PSMetricsInterface                = (*MockPSMetrics)(nil)
	_ AppMetricsInterface               = (*MockAppMetrics)(nil)
)
