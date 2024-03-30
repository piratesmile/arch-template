package errdefs

type Code string

const (
	codeSuccess             Code = "success"
	codeInternalServerError Code = "internal"
	codeInvalidParams       Code = "invalid"
	codeUnauthorized        Code = "unauthorized"
	codeTooManyRequests     Code = "frequency"
	codeNotFound            Code = "not-found"
)

var (
	ErrTooManyRequests  = NewBizError(codeTooManyRequests, "")
	ErrUnauthorized     = NewBizError(codeUnauthorized, "")
	ErrResourceNotFound = NewBizError(codeNotFound, "")
)

func InvalidParams(err error) error {
	return NewBizError(codeInvalidParams, err.Error())
}
