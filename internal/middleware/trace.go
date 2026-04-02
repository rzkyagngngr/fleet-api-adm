package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const TraceIDKey = "trace_id"
const TraceHeader = "X-Trace-Id"

// TraceMiddleware handles the creation and propagation of a unique Trace ID
// for every request-response cycle.
func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get from header (frontend generated)
		traceID := c.GetHeader(TraceHeader)

		// 2. If empty, generate new UUID as fallback (short version hex)
		if traceID == "" {
			u := uuid.New().String()
			traceID = strings.ReplaceAll(u, "-", "")[:16]
		}

		// 3. Set in Gin Context for logging
		c.Set(TraceIDKey, traceID)

		// 4. Set in Response Header for frontend to see
		c.Writer.Header().Set(TraceHeader, traceID)

		c.Next()
	}
}
