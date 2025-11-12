package middleware

import (
	"strconv"
	"time"

	"chatserver/metrics"

	"github.com/gin-gonic/gin"
)

// PrometheusMiddleware tracks HTTP metrics for all requests
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Increment active connections
		metrics.IncrementActiveConnections()
		defer metrics.DecrementActiveConnections()

		// Process request
		c.Next()

		// Record metrics
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		endpoint := c.FullPath()
		method := c.Request.Method

		// If endpoint is empty, use the raw path
		if endpoint == "" {
			endpoint = c.Request.URL.Path
		}

		metrics.HttpRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
		metrics.HttpRequestDuration.WithLabelValues(method, endpoint, status).Observe(duration)
	}
}
