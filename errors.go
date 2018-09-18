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

// Stack return Sentry Stack trace
func (e *AppError) Stack() *raven.Stacktrace { return e.stack }

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
	return &AppError{InternalErrorStatus, "", msg, nil, nil, nil}
}

func Newf(format string, v ...interface{}) error {
	return &AppError{InternalErrorStatus, "", fmt.Sprintf(format, v...), nil, nil, nil}
}

func Error(status int, code, msg string) *AppError {
	return &AppError{status, code, msg, nil, nil, nil}
}

func Errorf(status int, code, format string, v ...interface{}) *AppError {
	return &AppError{status, code, fmt.Sprintf(format, v...), nil, nil, nil}
}

func BadRequest(code, msg string) *AppError {
	return &AppError{BadRequestStatus, code, msg, nil, nil, nil}
}

func BadRequestf(code, format string, v ...interface{}) *AppError {
	return &AppError{BadRequestStatus, code, fmt.Sprintf(format, v...), nil, nil, nil}
}

func Unauthorized(code, msg string) *AppError {
	return &AppError{UnauthorizedStatus, code, msg, nil, nil, nil}
}

func Unauthorizedf(code, format string, v ...interface{}) *AppError {
	return &AppError{UnauthorizedStatus, code, fmt.Sprintf(format, v...), nil, nil, nil}
}

func Forbidden(code, msg string) *AppError {
	return &AppError{ForbiddenStatus, code, msg, nil, nil, nil}
}

func Forbiddenf(code, format string, v ...interface{}) *AppError {
	return &AppError{ForbiddenStatus, code, fmt.Sprintf(format, v...), nil, nil, nil}
}

func NotFound(code, msg string) *AppError {
	return &AppError{NotFoundStatus, code, msg, nil, nil, nil}
}

func NotFoundf(code, format string, v ...interface{}) *AppError {
	return &AppError{NotFoundStatus, code, fmt.Sprintf(format, v...), nil, nil, nil}
}

func InternalError(code, msg string) *AppError {
	return &AppError{InternalErrorStatus, code, msg, nil, nil, nil}
}

func InternalErrorf(code, format string, v ...interface{}) *AppError {
	return &AppError{InternalErrorStatus, code, fmt.Sprintf(format, v...), nil, nil, nil}
}

func Timeout(code, msg string) *AppError {
	return &AppError{TimeoutStatus, code, msg, nil, nil, nil}
}

func Timeoutf(code, format string, v ...interface{}) *AppError {
	return &AppError{TimeoutStatus, code, fmt.Sprintf(format, v...), nil, nil, nil}
}

func NotImplement(code, msg string) *AppError {
	return &AppError{NotFoundStatus, code, msg, nil, nil, nil}
}

func NotImplementf(code, format string, v ...interface{}) *AppError {
	return &AppError{NotFoundStatus, code, fmt.Sprintf(format, v...), nil, nil, nil}
}

func Unavailable(code, msg string) *AppError {
	return &AppError{UnavailableStatus, code, msg, nil, nil, nil}
}

func Unavailablef(code, format string, v ...interface{}) *AppError {
	return &AppError{UnavailableStatus, code, fmt.Sprintf(format, v...), nil, nil, nil}
}

func UnknownError(code, msg string) *AppError {
	return &AppError{UnknownErrorStatus, code, msg, nil, nil, nil}
}

func UnknownErrorf(code, format string, v ...interface{}) *AppError {
	return &AppError{UnknownErrorStatus, code, fmt.Sprintf(format, v...), nil, nil, nil}
}

func WithCaller(err error) *AppError {
	if err == nil {
		return nil
	}

	herr, ok := err.(*AppError)
	if !ok {
		appErr := &AppError{InternalErrorStatus, "", err.Error(), nil, nil, nil}
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
		return 520
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
