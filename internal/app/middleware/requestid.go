package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const contextKeyRequestID = "ctx-request-id"

func (m *middleware) RequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := ctx.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx.Header("X-Request-ID", requestID)
		ctx.Set(contextKeyRequestID, requestID)
	}
}

func RequestIDFromContext(ctx context.Context) string {
	v := ctx.Value(contextKeyRequestID)
	if v == nil {
		return ""
	}
	t, _ := v.(string)
	return t
}
