package errs

import "errors"

type Error interface {
	Is(err error) bool
	Wrap() error
	WrapMsg(message string) error
	error
}

func New(message string) Error {
	return &errorString{
		msg:   message,
		stack: callers(stackSkip - 1),
	}
}

type errorString struct {
	msg string
	*stack
}

func (e *errorString) Is(err error) bool {
	if err == nil {
		return false
	}
	var t *errorString
	ok := errors.As(err, &t)
	return ok && e.msg == t.msg
}

func (e *errorString) Wrap() error {
	return newErrorWrapper(e)
}

func (e *errorString) WrapMsg(message string) error {
	return newErrorWrapperMsg(e, message)
}

func (e *errorString) Error() string {
	if e.stack == nil {
		return "Error: " + e.msg
	}
	return "Error: " + e.msg + e.stack.String()
}
