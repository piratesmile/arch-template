package middleware

import (
	"github.com/gin-gonic/gin"
)

const contextKeyLogger = "ctx-logger"

func (m *middleware) PopulateLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var attrs = []interface{}{
			"url", ctx.Request.URL,
		}
		requestID := RequestIDFromContext(ctx)
		if requestID != "" {
			attrs = append(attrs, "request-id", requestID)
		}
		logger := m.logger.With(attrs...)
		ctx.Set(contextKeyLogger, logger)
	}
}
