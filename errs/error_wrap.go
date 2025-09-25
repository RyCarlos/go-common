package errs

import "fmt"

type ErrWrapper interface {
	Unwrap() error
	error
}

func Wrap(err error) error {
	return newErrorWrapper(err)
}

func newErrorWrapper(err error) error {
	if err == nil {
		return nil
	}
	return &errorWrapper{
		err:   err,
		stack: callers(stackSkip),
	}
}

func newErrorWrapperMsg(err error, msg string) error {
	if err == nil {
		return nil
	}
	err = fmt.Errorf("%w -> %s", err, msg)
	return &errorWrapper{
		err,
		msg,
		callers(stackSkip),
	}
}

type errorWrapper struct {
	err   error
	msg   string
	stack *stack
}

func (e *errorWrapper) Unwrap() error {
	return e.err
}

func (e *errorWrapper) Error() string {
	return e.err.Error() + e.stack.String()
}

func (e *errorWrapper) Cause() error {
	return e.err
}

func WrapMsg(err error, message string) error {
	if err == nil {
		return nil
	}
	return newErrorWrapperMsg(err, message)
}

func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}
