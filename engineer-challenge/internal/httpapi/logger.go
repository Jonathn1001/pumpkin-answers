package httpapi

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// requestLogger emits one structured log line per request via slog: method,
// path, status, latency, client IP, and the tenant slug when the route has one.
// Level reflects the outcome: 5xx -> Error, 4xx -> Warn, otherwise Info.
func requestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		// Health probes hit on a fixed interval; logging them buries real traffic.
		if c.Request.URL.Path == "/healthz" {
			return
		}

		status := c.Writer.Status()
		attrs := []any{
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", status),
			slog.Duration("latency", time.Since(start)),
			slog.String("client_ip", c.ClientIP()),
		}
		if slug := c.Param("slug"); slug != "" {
			attrs = append(attrs, slog.String("tenant", slug))
		}

		switch {
		case status >= 500:
			logger.Error("request", attrs...)
		case status >= 400:
			logger.Warn("request", attrs...)
		default:
			logger.Info("request", attrs...)
		}
	}
}
