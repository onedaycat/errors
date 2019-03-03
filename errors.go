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

func NewInput() Input {
	return make(map[string]interface{})
}

func (i Input) Set(key string, value interface{}) {
	i[key] = value
}

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

type Error interface {
	Error() string
	String() string
	StackTrace() raven.Interface
	WithPanic() *AppError
	WithCause(err error) *AppError
	WithCaller() *AppError
	WithCallerSkip(skip int) *AppError
	WithInput(input interface{}) *AppError
	Format(s fmt.State, verb rune)
	StackStrings() []string
	GetStatus() int
	GetCode() string
	GetType() string
	GetMessage() string
	GetInput() interface{}
	GetPanic() bool
	GetCause() error
}

// AppError error
type AppError struct {
	Status  int         `json:"status"`
	Code    string      `json:"code"`
	Type    string      `json:"type"`
	Message string      `json:"message"`
	Input   interface{} `json:"input"`
	Panic   bool
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

func (e *AppError) GetStatus() int {
	return e.Status
}

func (e *AppError) GetCode() string {
	return e.Code
}

func (e *AppError) GetType() string {
	return e.Type
}

func (e *AppError) GetMessage() string {
	return e.Message
}

func (e *AppError) GetInput() interface{} {
	return e.Input
}

func (e *AppError) GetPanic() bool {
	return e.Panic
}

func (e *AppError) GetCause() error {
	return e.Cause
}

// String interface
func (e *AppError) String() string {
	if e.Code == "" {
		return e.Message
	}

	return e.Code + ": " + e.Message
}

// Stack return Sentry Stack trace
func (e *AppError) StackTrace() raven.Interface {
	if e.Cause != nil {
		return &raven.Exception{
			Stacktrace: e.stack,
			Value:      e.Cause.Error(),
			Type:       e.Code,
		}
	}

	return &raven.Exception{
		Stacktrace: e.stack,
		Value:      e.Error(),
		Type:       e.Code,
	}
}

func (e *AppError) WithPanic() *AppError {
	e.Panic = true
	return e
}

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
func (e *AppError) WithInput(input interface{}) *AppError {
	e.Input = input
	return e
}

func (e *AppError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		io.WriteString(s, e.Error())

		if s.Flag('+') {
			if e.Cause != nil {
				fmt.Fprintf(s, "\n%s", e.Cause.Error())
			}
			if e.stack != nil {
				// fmt.Fprintln(s)
				for _, frame := range e.stack.Frames {
					fmt.Fprintf(s, "\n%s:%d\t%s", frame.Function, frame.Lineno, frame.AbsolutePath)
				}
			}
		}

	case 's':
		io.WriteString(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

func (e *AppError) StackStrings() []string {
	stacks := make([]string, 0, len(e.stack.Frames))
	if e.stack != nil {
		for _, frame := range e.stack.Frames {
			stacks = append(stacks, fmt.Sprintf("%s:%d %s", frame.Function, frame.Lineno, frame.AbsolutePath))
		}
	}

	return stacks
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

func NewError(status int, errorType, code, msg string) *AppError {
	return &AppError{Status: status, Type: errorType, Code: code, Message: msg, Input: nil, Cause: nil, stack: nil}
}

func NewErrorf(status int, errorType, code, format string, v ...interface{}) *AppError {
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
