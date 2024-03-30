package errdefs

const (
	codeIncorrectPassword Code = "user.incorrect-password"
	codeUserAlreadyExists Code = "user.already-exists"
)

var (
	ErrUserAlreadyExists = NewBizError(codeUserAlreadyExists, "")
	ErrIncorrectPassword = NewBizError(codeIncorrectPassword, "")
)
