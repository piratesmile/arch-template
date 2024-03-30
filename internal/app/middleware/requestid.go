package middleware

import (
	"arch-template/internal/app/public"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (m *middleware) RequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := ctx.GetHeader(public.XRequestIDKey)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx.Header(public.XRequestIDKey, requestID)
		ctx.Set(public.XRequestIDKey, requestID)
		ctx.Next()
	}
}
