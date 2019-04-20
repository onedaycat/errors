package errors

import "fmt"

type ErrorDefinition struct {
	Code    string
	Type    string
	Status  int
	Message string
}

func Def(code, errType string, status int, msg ...string) *ErrorDefinition {
	e := &ErrorDefinition{
		Code:   code,
		Type:   errType,
		Status: status,
	}

	if len(msg) > 0 {
		e.Message = msg[0]
	}

	return e
}

func (e *ErrorDefinition) Msg(msg string) *ErrorDefinition {
	e.Message = msg

	return e
}

func (e *ErrorDefinition) Msgf(format string, v ...interface{}) *ErrorDefinition {
	e.Message = fmt.Sprintf(format, v...)

	return e
}

func (e *ErrorDefinition) New(msg ...string) *AppError {
	if len(msg) > 0 {
		return NewError(e.Status, e.Type, e.Code, msg[0])
	}

	return NewError(e.Status, e.Type, e.Code, e.Message)
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
	return err != nil && e.Code == err.GetCode()
}

func DefBadRequest(code string, msg ...string) *ErrorDefinition {
	e := &ErrorDefinition{
		Code:   code,
		Type:   BadRequestType,
		Status: 400,
	}
	if len(msg) > 0 {
		e.Message = msg[0]
	}

	return e
}

func DefUnauthorized(code string, msg ...string) *ErrorDefinition {
	e := &ErrorDefinition{
		Code:   code,
		Type:   UnauthorizedType,
		Status: 401,
	}
	if len(msg) > 0 {
		e.Message = msg[0]
	}

	return e
}

func DefForbidden(code string, msg ...string) *ErrorDefinition {
	e := &ErrorDefinition{
		Code:   code,
		Type:   ForbiddenType,
		Status: 403,
	}
	if len(msg) > 0 {
		e.Message = msg[0]
	}

	return e
}

func DefNotFound(code string, msg ...string) *ErrorDefinition {
	e := &ErrorDefinition{
		Code:   code,
		Type:   NotFoundType,
		Status: 404,
	}
	if len(msg) > 0 {
		e.Message = msg[0]
	}

	return e
}

func DefTimeout(code string, msg ...string) *ErrorDefinition {
	e := &ErrorDefinition{
		Code:   code,
		Type:   TimeoutType,
		Status: 441,
	}
	if len(msg) > 0 {
		e.Message = msg[0]
	}

	return e
}

func DefInternalError(code string, msg ...string) *ErrorDefinition {
	e := &ErrorDefinition{
		Code:   code,
		Type:   InternalErrorType,
		Status: 500,
	}
	if len(msg) > 0 {
		e.Message = msg[0]
	}

	return e
}

func DefNotImplement(code string, msg ...string) *ErrorDefinition {
	e := &ErrorDefinition{
		Code:   code,
		Type:   NotImplementType,
		Status: 501,
	}
	if len(msg) > 0 {
		e.Message = msg[0]
	}

	return e
}

func DefUnavailable(code string, msg ...string) *ErrorDefinition {
	e := &ErrorDefinition{
		Code:   code,
		Type:   UnavailableType,
		Status: 503,
	}
	if len(msg) > 0 {
		e.Message = msg[0]
	}

	return e
}

func DefUnknownError(code string, msg ...string) *ErrorDefinition {
	e := &ErrorDefinition{
		Code:   code,
		Type:   UnknownErrorType,
		Status: 520,
	}
	if len(msg) > 0 {
		e.Message = msg[0]
	}

	return e
}
