package errs

import (
	"fmt"
)

func ErrPanic(r any) error {
	return ErrPanicMsg(r, UnknownError, "panic error", 9)
}

func ErrPanicMsg(r any, code int, msg string, skip int) error {
	if r == nil {
		return nil
	}
	err := &errorCode{
		code:   code,
		msg:    fmt.Sprint(r),
		detail: fmt.Sprint(r),
	}
	return &errorWrapper{
		err:   err,
		stack: callers(skip),
	}
}
