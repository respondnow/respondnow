package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/respondnow/respond/server/pkg/prometheus"
)

func RequestMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		status := c.Writer.Status()
		statusLabel := http.StatusText(status)

		// Increment the total requests counter
		prometheus.TotalRequests.WithLabelValues(statusLabel, c.Request.URL.Path).Inc()

		// Increment the error counter for 4xx and 5xx responses
		if status >= 400 {
			prometheus.ErrorRequests.WithLabelValues(statusLabel, c.Request.URL.Path).Inc()
		}

		if status >= 400 && status < 500 {
			prometheus.Error4xxRequests.WithLabelValues(c.Request.URL.Path).Inc()
		}

		if status >= 500 {
			prometheus.Error5xxRequests.WithLabelValues(c.Request.URL.Path).Inc()
		}
	}
}

func SLIAPIResponseTimeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		duration := time.Since(start)
		prometheus.ResponseTime.WithLabelValues(c.Request.URL.Path).Observe(float64(duration.Milliseconds()))
		prometheus.ResponseTimeInSeconds.WithLabelValues(c.Request.URL.Path).Observe(duration.Seconds())
	}
}
