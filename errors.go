package errors

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/getsentry/raven-go"
)

var (
	TraceContextLines = 3
	TraceSkipFrames   = 1
	delim             = ": "
)

var DumbError = InternalError("DUMB", "Dumb Error")

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
	StackTrace() *raven.Stacktrace
	WithPanic() *AppError
	WithCause(err error) *AppError
	WithCaller() *AppError
	WithCallerSkip(skip int) *AppError
	WithInput(input interface{}) *AppError
	WithStatus(status int) *AppError
	WithCode(code string) *AppError
	WithType(errType string) *AppError
	WithMessage(msg string) *AppError
	WithNotFound() *AppError
	WithBadRequest() *AppError
	WithUnauthorized() *AppError
	WithForbidden() *AppError
	WithInternalError() *AppError
	WithTimeoutError() *AppError
	WithUnknownError() *AppError
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
	Status  int         `json:"status,omitempty"`
	Code    string      `json:"code,omitempty"`
	Type    string      `json:"type,omitempty"`
	Message string      `json:"message,omitempty"`
	Input   interface{} `json:"-"`
	Panic   bool        `json:"-"`
	Cause   error       `json:"-"`
	stack   *raven.Stacktrace
}

// Error error
func (e *AppError) Error() string {
	if e.Code == "" {
		return e.Message
	}

	return e.Code + delim + e.Message
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

	return e.Code + delim + e.Message
}

// Stack return Sentry Stack trace
func (e *AppError) StackTrace() *raven.Stacktrace {
	return e.stack
}

func (e *AppError) WithStatus(status int) *AppError {
	e.Status = status
	return e
}

func (e *AppError) WithCode(code string) *AppError {
	e.Code = code
	return e
}

func (e *AppError) WithType(errType string) *AppError {
	e.Type = errType
	return e
}

func (e *AppError) WithMessage(msg string) *AppError {
	e.Message = msg
	return e
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

// WithCauseMessage error
func (e *AppError) WithCauseMessage(msg string) *AppError {
	e.Cause = errors.New(msg)
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

func (e *AppError) WithNotFound() *AppError {
	e.Status = NotFoundStatus
	e.Type = NotFoundType
	return e
}

func (e *AppError) WithBadRequest() *AppError {
	e.Status = NotFoundStatus
	e.Type = NotFoundType
	return e
}

func (e *AppError) WithUnauthorized() *AppError {
	e.Status = UnauthorizedStatus
	e.Type = UnauthorizedType
	return e
}

func (e *AppError) WithForbidden() *AppError {
	e.Status = ForbiddenStatus
	e.Type = ForbiddenType
	return e
}

func (e *AppError) WithInternalError() *AppError {
	e.Status = InternalErrorStatus
	e.Type = InternalErrorType
	return e
}

func (e *AppError) WithTimeoutError() *AppError {
	e.Status = TimeoutStatus
	e.Type = TimeoutType
	return e
}

func (e *AppError) WithUnknownError() *AppError {
	e.Status = UnknownErrorStatus
	e.Type = UnknownErrorType
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
			if e.stack != nil {
				// fmt.Fprintln(s)
				for _, frame := range e.stack.Frames {
					fmt.Fprintf(s, "\n%s\t%s:%d", frame.Function, frame.AbsolutePath, frame.Lineno)
				}
			}

			if e.Cause != nil {
				herr, ok := e.Cause.(Error)
				if ok {
					fmt.Fprintf(s, "\n\n%+v\n", herr)
				} else {
					fmt.Fprintf(s, "\n%s", e.Cause.Error())
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

func Wrap(err error) Error {
	if err == nil {
		return nil
	}

	return New(err.Error())
}

func Convert(goerr error, err Error) Error {
	if goerr == nil {
		return nil
	}

	return err.WithCause(goerr)
}

func Simple(msg string) error {
	return errors.New(msg)
}

func New(msg string) *AppError {
	return &AppError{Status: InternalErrorStatus, Type: InternalErrorType, Code: msg, Message: msg, Input: nil, Cause: nil, stack: nil}
}

func Newf(format string, v ...interface{}) *AppError {
	msg := fmt.Sprintf(format, v...)
	return &AppError{Status: InternalErrorStatus, Type: InternalErrorType, Code: msg, Message: msg, Input: nil, Cause: nil, stack: nil}
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

func ParseErrorMessage(errStr string) (string, string) {
	es := strings.SplitN(errStr, delim, 2)
	if len(es) == 2 {
		return es[0], es[1]
	}

	return "", errStr
}

func ParseError(errStr string) *AppError {
	es := strings.SplitN(errStr, delim, 2)
	if len(es) == 2 {
		return UnknownError(es[0], es[1])
	}

	return UnknownError(errStr, errStr)
}

func WithCaller(err error) *AppError {
	if err == nil {
		return nil
	}

	herr, ok := err.(*AppError)
	if !ok {
		return New(err.Error()).WithCallerSkip(2)
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

func IsNotFound(err Error) bool {
	return err != nil && err.GetStatus() == 404
}

func IsInternalError(err Error) bool {
	return err != nil && err.GetStatus() == 500
}

func IsBadRequest(err Error) bool {
	return err != nil && err.GetStatus() == 400
}

func IsUnauthorized(err Error) bool {
	return err != nil && err.GetStatus() == 401
}

func IsForbidden(err Error) bool {
	return err != nil && err.GetStatus() == 403
}
