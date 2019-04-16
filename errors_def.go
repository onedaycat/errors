package errors

import "fmt"

type ErrorDefinition struct {
	Code   string
	Type   string
	Status int
	Msg    string
}

func Def(code, errType string) *ErrorDefinition {
	e := &ErrorDefinition{
		Code: code,
		Type: errType,
	}

	switch errType {
	case BadRequestType:
		e.Status = 400
	case UnauthorizedType:
		e.Status = 401
	case ForbiddenType:
		e.Status = 403
	case NotFoundType:
		e.Status = 404
	case TimeoutType:
		e.Status = 441
	case InternalErrorType:
		e.Status = 500
	case NotImplementType:
		e.Status = 501
	case UnavailableType:
		e.Status = 503
	case UnknownErrorType:
		e.Status = 520
	}

	return e
}

func (e *ErrorDefinition) Message(msg string) {
	e.Msg = msg
}

func (e *ErrorDefinition) New(msg ...string) *AppError {
	if len(msg) > 0 {
		return NewError(e.Status, e.Type, e.Code, msg[0])
	}

	return NewError(e.Status, e.Type, e.Code, e.Msg)
}

func (e *ErrorDefinition) Newf(format string, v ...interface{}) *AppError {
	msg := fmt.Sprintf(format, v...)
	return NewError(e.Status, e.Type, e.Code, msg)
}

func (e *ErrorDefinition) WithPanic() *AppError {
	return e.New().WithPanic()
}

func (e *ErrorDefinition) WithCause(err error) *AppError {
	return e.New().WithCause(err)
}

func (e *ErrorDefinition) WithCaller() *AppError {
	return e.New().WithCaller()
}

func (e *ErrorDefinition) WithCallerSkip(skip int) *AppError {
	return e.New().WithCallerSkip(skip)
}

func (e *ErrorDefinition) WithInput(input interface{}) *AppError {
	return e.New().WithInput(input)
}

func (e *ErrorDefinition) Is(err Error) bool {
	return e.Code == err.GetCode()
}
