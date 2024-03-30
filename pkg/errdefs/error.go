package errdefs

import (
	"errors"
	"fmt"
	"net/http"
)

type BizError struct {
	Code    Code
	Message string
}

func NewBizError(code Code, message string) BizError {
	return BizError{code, message}
}

func (b BizError) Error() string {
	return fmt.Sprintf("[BizError] code=%v, message=%s", b.Code, b.Message)
}

func Decode(err error) (httpCode int, code Code, message string) {
	if err == nil {
		return http.StatusOK, codeSuccess, ""
	}

	var berr BizError
	switch {
	case errors.As(err, &berr):
		return http.StatusBadRequest, berr.Code, berr.Message
	}

	return http.StatusInternalServerError, codeInternalServerError, http.StatusText(http.StatusInternalServerError)
}
