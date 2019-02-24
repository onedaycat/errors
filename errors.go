package errors

import (
	"fmt"
	"io"

	"github.com/getsentry/raven-go"
)

var (
	TraceContextLines = 3
	TraceSkipFrames   = 1
)

type Input map[string]interface{}

const (
	GolangErrorType   = "GolangError"
	BadRequestType    = "BadRequest"
	UnauthorizedType  = "Unauthorized"
	ForbiddenType     = "Forbidden"
	NotFoundType      = "NotFound"
	TimeoutType       = "Timeout"
	InternalErrorType = "InternalError"
	NotImplementType  = "NotImplement"
	UnavailableType   = "Unavailable"
	UnknownErrorType  = "UnknownError"
)

const (
	GolangErrorStatus   = 0
	BadRequestStatus    = 400
	UnauthorizedStatus  = 401
	ForbiddenStatus     = 403
	NotFoundStatus      = 404
	TimeoutStatus       = 441
	InternalErrorStatus = 500
	NotImplementStatus  = 501
	UnavailableStatus   = 503
	UnknownErrorStatus  = 520
)

// AppError error
type AppError struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
	Input   Input  `json:"input"`
	Cause   error
	stack   *raven.Stacktrace
}

// Error error
func (e *AppError) Error() string {
	if e.Code == "" {
		return e.Message
	}

	return e.Code + ": " + e.Message
}

// String interface
func (e *AppError) String() string {
	if e.Code == "" {
		return e.Message
	}

	return e.Code + ": " + e.Message
}

// Stack return Sentry Stack trace
func (e *AppError) StackTrace() *raven.Stacktrace { return e.stack }

// WithCause error
func (e *AppError) WithCause(err error) *AppError {
	e.Cause = err
	return e
}

// WithCause error
func (e *AppError) WithCaller() *AppError {
	e.stack = raven.NewStacktrace(TraceSkipFrames, TraceContextLines, nil)

	return e
}

// WithCause error
func (e *AppError) WithCallerSkip(skip int) *AppError {
	e.stack = raven.NewStacktrace(skip, TraceContextLines, nil)

	return e
}

// WithInput add input
func (e *AppError) WithInput(input Input) *AppError {
	e.Input = input
	return e
}

func (e *AppError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		fmt.Fprintln(s, e.Message)
		if e.Cause != nil {
			fmt.Fprintln(s, e.Cause)
		}

		if s.Flag('+') {
			for _, frame := range e.stack.Frames {
				fmt.Fprintf(s, "%s:%d\t%s\n", frame.Function, frame.Lineno, frame.AbsolutePath)
			}
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.Message)
	case 'q':
		fmt.Fprintf(s, "%q", e.Message)
	}
}

func Wrap(err error, AppError *AppError) error {
	if err != nil {
		return nil
	}

	return AppError.WithCause(err)
}

func New(msg string) error {
	return &AppError{Status: InternalErrorStatus, Type: InternalErrorType, Code: "", Message: msg, Input: nil, Cause: nil, stack: nil}
}

func Newf(format string, v ...interface{}) error {
	return &AppError{Status: InternalErrorStatus, Type: InternalErrorType, Code: "", Message: fmt.Sprintf(format, v...), Input: nil, Cause: nil, stack: nil}
}

func Error(status int, errorType, code, msg string) *AppError {
	return &AppError{Status: status, Type: errorType, Code: code, Message: msg, Input: nil, Cause: nil, stack: nil}
}

func Errorf(status int, errorType, code, format string, v ...interface{}) *AppError {
	return &AppError{Status: status, Type: errorType, Code: code, Message: fmt.Sprintf(format, v...), Input: nil, Cause: nil, stack: nil}
}

func BadRequest(code, msg string) *AppError {
	return &AppError{Status: BadRequestStatus, Type: BadRequestType, Code: code, Message: msg, Input: nil, Cause: nil, stack: nil}
}

func BadRequestf(code, format string, v ...interface{}) *AppError {
	return &AppError{Status: BadRequestStatus, Type: BadRequestType, Code: code, Message: fmt.Sprintf(format, v...), Input: nil, Cause: nil, stack: nil}
}

func Unauthorized(code, msg string) *AppError {
	return &AppError{Status: UnauthorizedStatus, Type: UnauthorizedType, Code: code, Message: msg, Input: nil, Cause: nil, stack: nil}
}

func Unauthorizedf(code, format string, v ...interface{}) *AppError {
	return &AppError{Status: UnauthorizedStatus, Type: UnauthorizedType, Code: code, Message: fmt.Sprintf(format, v...), Input: nil, Cause: nil, stack: nil}
}

func Forbidden(code, msg string) *AppError {
	return &AppError{Status: ForbiddenStatus, Type: ForbiddenType, Code: code, Message: msg, Input: nil, Cause: nil, stack: nil}
}

func Forbiddenf(code, format string, v ...interface{}) *AppError {
	return &AppError{Status: ForbiddenStatus, Type: ForbiddenType, Code: code, Message: fmt.Sprintf(format, v...), Input: nil, Cause: nil, stack: nil}
}

func NotFound(code, msg string) *AppError {
	return &AppError{Status: NotFoundStatus, Type: NotFoundType, Code: code, Message: msg, Input: nil, Cause: nil, stack: nil}
}

func NotFoundf(code, format string, v ...interface{}) *AppError {
	return &AppError{Status: NotFoundStatus, Type: NotFoundType, Code: code, Message: fmt.Sprintf(format, v...), Input: nil, Cause: nil, stack: nil}
}

func InternalError(code, msg string) *AppError {
	return &AppError{Status: InternalErrorStatus, Type: InternalErrorType, Code: code, Message: msg, Input: nil, Cause: nil, stack: nil}
}

func InternalErrorf(code, format string, v ...interface{}) *AppError {
	return &AppError{Status: InternalErrorStatus, Type: InternalErrorType, Code: code, Message: fmt.Sprintf(format, v...), Input: nil, Cause: nil, stack: nil}
}

func Timeout(code, msg string) *AppError {
	return &AppError{Status: TimeoutStatus, Type: TimeoutType, Code: code, Message: msg, Input: nil, Cause: nil, stack: nil}
}

func Timeoutf(code, format string, v ...interface{}) *AppError {
	return &AppError{Status: TimeoutStatus, Type: TimeoutType, Code: code, Message: fmt.Sprintf(format, v...), Input: nil, Cause: nil, stack: nil}
}

func NotImplement(code, msg string) *AppError {
	return &AppError{Status: NotFoundStatus, Type: NotImplementType, Code: code, Message: msg, Input: nil, Cause: nil, stack: nil}
}

func NotImplementf(code, format string, v ...interface{}) *AppError {
	return &AppError{Status: NotFoundStatus, Type: NotImplementType, Code: code, Message: fmt.Sprintf(format, v...), Input: nil, Cause: nil, stack: nil}
}

func Unavailable(code, msg string) *AppError {
	return &AppError{Status: UnavailableStatus, Type: UnavailableType, Code: code, Message: msg, Input: nil, Cause: nil, stack: nil}
}

func Unavailablef(code, format string, v ...interface{}) *AppError {
	return &AppError{Status: UnavailableStatus, Type: UnavailableType, Code: code, Message: fmt.Sprintf(format, v...), Input: nil, Cause: nil, stack: nil}
}

func UnknownError(code, msg string) *AppError {
	return &AppError{Status: UnknownErrorStatus, Type: UnknownErrorType, Code: code, Message: msg, Input: nil, Cause: nil, stack: nil}
}

func UnknownErrorf(code, format string, v ...interface{}) *AppError {
	return &AppError{Status: UnknownErrorStatus, Type: UnknownErrorType, Code: code, Message: fmt.Sprintf(format, v...), Input: nil, Cause: nil, stack: nil}
}

func WithCaller(err error) *AppError {
	if err == nil {
		return nil
	}

	herr, ok := err.(*AppError)
	if !ok {
		appErr := &AppError{Status: InternalErrorStatus, Type: InternalErrorType, Code: "", Message: err.Error(), Input: nil, Cause: nil, stack: nil}
		appErr.WithCallerSkip(2)
		return appErr
	}

	return herr.WithCallerSkip(2)
}

func Cause(err error) error {
	herr, ok := err.(*AppError)
	if !ok {
		return err
	}

	return herr.Cause
}

func FromError(err error) (*AppError, bool) {
	herr, ok := err.(*AppError)
	return herr, ok
}

func ErrStatus(err error) int {
	herr, ok := err.(*AppError)
	if !ok {
		return GolangErrorStatus
	}

	return herr.Status
}

func IsNotFound(err error) bool {
	return err != nil && ErrStatus(err) == 404
}

func IsInternalError(err error) bool {
	return err != nil && ErrStatus(err) == 500
}

func IsBadRequest(err error) bool {
	return err != nil && ErrStatus(err) == 400
}

func IsUnauthorized(err error) bool {
	return err != nil && ErrStatus(err) == 401
}

func IsForbidden(err error) bool {
	return err != nil && ErrStatus(err) == 403
}
