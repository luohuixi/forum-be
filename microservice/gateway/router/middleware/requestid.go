package middleware

import (
	"forum/pkg/tracer"

	"github.com/gin-gonic/gin"
)

func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for incoming header, use it if exists
		//requestId := c.Request.Header.Get("X-Request-Id")

		// Create request id with UUID4
		//if requestId == "" {
		//	u4 := uuid.NewV4()
		//	requestId = u4.String()
		//}

		// 用 traceId 代替 requestId
		traceId := tracer.GetTraceId(c.Request.Context())

		// Expose it for use in the application
		c.Set("X-Request-Id", traceId)

		// Set X-Request-Id header
		c.Writer.Header().Set("X-Request-Id", traceId)
		c.Next()
	}
}
