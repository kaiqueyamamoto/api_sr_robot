package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP Metrics
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "status"},
	)

	// Authentication Metrics
	AuthAttemptsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_attempts_total",
			Help: "Total number of authentication attempts",
		},
		[]string{"type", "status"}, // type: login/register, status: success/failure
	)

	ActiveUsers = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_users_total",
			Help: "Number of currently active users",
		},
	)

	TokensIssued = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "jwt_tokens_issued_total",
			Help: "Total number of JWT tokens issued",
		},
	)

	TokenValidationFailures = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "jwt_token_validation_failures_total",
			Help: "Total number of JWT token validation failures",
		},
		[]string{"reason"}, // reason: expired, invalid, malformed
	)

	// Database Metrics
	DatabaseOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_operations_total",
			Help: "Total number of database operations",
		},
		[]string{"operation", "collection", "status"}, // operation: insert/find/update/delete
	)

	DatabaseOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_operation_duration_seconds",
			Help:    "Duration of database operations in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "collection"},
	)

	// Chat Metrics
	ChatMessagesTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "chat_messages_total",
			Help: "Total number of chat messages processed",
		},
	)

	// System Metrics
	ActiveConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_connections",
			Help: "Number of active connections",
		},
	)
)

// Helper functions to record metrics
func RecordAuthAttempt(authType, status string) {
	AuthAttemptsTotal.WithLabelValues(authType, status).Inc()
}

func RecordTokenIssued() {
	TokensIssued.Inc()
}

func RecordTokenValidationFailure(reason string) {
	TokenValidationFailures.WithLabelValues(reason).Inc()
}

func RecordDatabaseOperation(operation, collection, status string, duration float64) {
	DatabaseOperationsTotal.WithLabelValues(operation, collection, status).Inc()
	DatabaseOperationDuration.WithLabelValues(operation, collection).Observe(duration)
}

func RecordChatMessage() {
	ChatMessagesTotal.Inc()
}

func IncrementActiveConnections() {
	ActiveConnections.Inc()
}

func DecrementActiveConnections() {
	ActiveConnections.Dec()
}

func SetActiveUsers(count float64) {
	ActiveUsers.Set(count)
}
