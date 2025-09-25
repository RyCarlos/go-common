package errs

import (
	"errors"
)

type ErrCode interface {
	Code() int
	Msg() string
	Error
}

func NewErrorCode(code int, msg string) ErrCode {
	return &errorCode{
		code: code,
		msg:  msg,
	}
}

type errorCode struct {
	error
	code   int
	msg    string
	detail string
}

func (e *errorCode) Code() int {
	return e.code
}

func (e *errorCode) Msg() string {
	return e.msg
}

func (e *errorCode) Is(err error) bool {
	if err == nil {
		return false
	}
	var codeErr *errorCode
	ok := errors.As(err, &codeErr)
	return ok && e.Code() == codeErr.Code() && e.Msg() == codeErr.Msg()
}

func (e *errorCode) Wrap() error {
	return newErrorWrapper(e)
}

func (e *errorCode) WrapMsg(message string) error {
	return newErrorWrapperMsg(e, message)
}

func (e *errorCode) Error() string {
	return "Error: " + e.msg
}
