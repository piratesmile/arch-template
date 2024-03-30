package middleware

import (
	"arch-template/configs"
	"arch-template/pkg/response"
	"arch-template/pkg/tlog"
	"arch-template/utils"
	"bytes"
	"fmt"
	"net/http/httputil"
	"strings"

	"github.com/gin-gonic/gin"
)

type responseWriter struct {
	gin.ResponseWriter
	buf *bytes.Buffer
}

func (w *responseWriter) Write(data []byte) (n int, err error) {
	w.buf.Write(data)
	return w.ResponseWriter.Write(data)
}

func (w *responseWriter) WriteString(s string) (n int, err error) {
	w.buf.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func (m *middleware) DebugLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		if m.config.APP.Env != configs.Dev || strings.EqualFold(ctx.FullPath(), "/favicon.ico") {
			return
		}

		request, err := httputil.DumpRequest(ctx.Request, true)
		if err != nil {
			response.Error(ctx, fmt.Errorf("dump request err:%w", err))
			return
		}
		tlog.Debug(ctx, "debug request", tlog.Fields{"body": utils.BytesToString(request)})

		rspWriter := &responseWriter{
			ResponseWriter: ctx.Writer,
			buf:            bytes.NewBuffer(nil),
		}
		ctx.Writer = rspWriter

		ctx.Next()
		tlog.Debug(ctx, "debug response", tlog.Fields{"body": rspWriter.buf.String()})
	}
}
