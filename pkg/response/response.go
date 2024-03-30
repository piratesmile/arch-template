package response

import (
	"arch-template/pkg/errdefs"
	"arch-template/pkg/tlog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code    errdefs.Code `json:"code"`
	Message string       `json:"message"`
}

func Response(ctx *gin.Context, httpCode int, data interface{}) {
	if data == nil {
		ctx.Status(httpCode)
	} else {
		ctx.AbortWithStatusJSON(httpCode, data)
	}
}

func Success(ctx *gin.Context, data interface{}) {
	Response(ctx, http.StatusOK, data)
}

func InternalServerError(ctx *gin.Context) {
	Response(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func Error(ctx *gin.Context, err error) {
	if err == nil {
		Response(ctx, http.StatusOK, nil)
		return
	}
	httpCode, bizCode, message := errdefs.Decode(err)

	if httpCode == http.StatusInternalServerError {
		// log error message if occur internal server error
		tlog.Error(ctx, "internal server error", tlog.Fields{"err": err})
	}
	Response(ctx, httpCode, ErrorResponse{bizCode, message})
}
