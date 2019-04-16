package errors

import "fmt"

type ErrorDefinition struct {
	Code   string
	Type   string
	Status int
	Msg    string
}

func Def(code string) *ErrorDefinition {
	e := &ErrorDefinition{
		Code:   code,
		Status: UnknownErrorStatus,
		Type:   UnknownErrorType,
	}

	return e
}

func DefM(code, msg string) *ErrorDefinition {
	e := Def(code)
	e.Msg = msg

	return e
}

func (e *ErrorDefinition) Message(msg string) *ErrorDefinition {
	e.Msg = msg

	return e
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

func (e *ErrorDefinition) BadRequest() *ErrorDefinition {
	e.Type = BadRequestType
	e.Status = 400
	return e
}

func (e *ErrorDefinition) Unauthorized() *ErrorDefinition {
	e.Type = UnauthorizedType
	e.Status = 401
	return e
}

func (e *ErrorDefinition) Forbidden() *ErrorDefinition {
	e.Type = ForbiddenType
	e.Status = 403
	return e
}

func (e *ErrorDefinition) NotFound() *ErrorDefinition {
	e.Type = NotFoundType
	e.Status = 404
	return e
}

func (e *ErrorDefinition) Timeout() *ErrorDefinition {
	e.Type = TimeoutType
	e.Status = 441
	return e
}

func (e *ErrorDefinition) InternalError() *ErrorDefinition {
	e.Type = InternalErrorType
	e.Status = 500
	return e
}

func (e *ErrorDefinition) NotImplement() *ErrorDefinition {
	e.Type = NotImplementType
	e.Status = 501
	return e
}

func (e *ErrorDefinition) Unavailable() *ErrorDefinition {
	e.Type = UnavailableType
	e.Status = 503
	return e
}

func (e *ErrorDefinition) UnknownError() *ErrorDefinition {
	e.Type = UnknownErrorType
	e.Status = 520
	return e
}
