package errors

import (
	"errors"
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap/zapcore"
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
	cause   error
}

// Error error
func (e *AppError) Error() string { return e.Code + ": " + e.Message }

// Cause error
func (e *AppError) Cause() error { return e.cause }

// WithCause error
func (e *AppError) WithCause(err error) *AppError {
	e.cause = err
	return e
}

// WithInput add input
func (e *AppError) WithInput(input Input) *AppError {
	e.Input = input
	return e
}

func (e *AppError) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("error", e.Message)
	if e.cause != nil {
		enc.AddString("error_cause", e.cause.Error())
	}

	if e.Input != nil {
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		input, err := json.Marshal(e.Input)
		if err != nil {
			return err
		}

		enc.AddByteString("input", input)
	}

	return nil
}

func Wrap(err error, AppError *AppError) error {
	if err != nil {
		return nil
	}

	return AppError.WithCause(err)
}

func New(msg string) error {
	return errors.New(msg)
}

func Newf(format string, v ...interface{}) error {
	return errors.New(fmt.Sprintf(format, v...))
}

func Error(status int, code, msg string) *AppError {
	return &AppError{status, code, msg, nil, nil}
}

func Errorf(status int, code, format string, v ...interface{}) *AppError {
	return &AppError{status, code, fmt.Sprintf(format, v...), nil, nil}
}

func BadRequest(code, msg string) *AppError {
	return &AppError{BadRequestStatus, code, msg, nil, nil}
}

func BadRequestf(code, format string, v ...interface{}) *AppError {
	return &AppError{BadRequestStatus, code, fmt.Sprintf(format, v...), nil, nil}
}

func Unauthorized(code, msg string) *AppError {
	return &AppError{UnauthorizedStatus, code, msg, nil, nil}
}

func Unauthorizedf(code, format string, v ...interface{}) *AppError {
	return &AppError{UnauthorizedStatus, code, fmt.Sprintf(format, v...), nil, nil}
}

func Forbidden(code, msg string) *AppError {
	return &AppError{ForbiddenStatus, code, msg, nil, nil}
}

func Forbiddenf(code, format string, v ...interface{}) *AppError {
	return &AppError{ForbiddenStatus, code, fmt.Sprintf(format, v...), nil, nil}
}

func NotFound(code, msg string) *AppError {
	return &AppError{NotFoundStatus, code, msg, nil, nil}
}

func NotFoundf(code, format string, v ...interface{}) *AppError {
	return &AppError{NotFoundStatus, code, fmt.Sprintf(format, v...), nil, nil}
}

func InternalError(code, msg string) *AppError {
	return &AppError{InternalErrorStatus, code, msg, nil, nil}
}

func InternalErrorf(code, format string, v ...interface{}) *AppError {
	return &AppError{InternalErrorStatus, code, fmt.Sprintf(format, v...), nil, nil}
}

func Timeout(code, msg string) *AppError {
	return &AppError{TimeoutStatus, code, msg, nil, nil}
}

func Timeoutf(code, format string, v ...interface{}) *AppError {
	return &AppError{TimeoutStatus, code, fmt.Sprintf(format, v...), nil, nil}
}

func NotImplement(code, msg string) *AppError {
	return &AppError{NotFoundStatus, code, msg, nil, nil}
}

func NotImplementf(code, format string, v ...interface{}) *AppError {
	return &AppError{NotFoundStatus, code, fmt.Sprintf(format, v...), nil, nil}
}

func Unavailable(code, msg string) *AppError {
	return &AppError{UnavailableStatus, code, msg, nil, nil}
}

func Unavailablef(code, format string, v ...interface{}) *AppError {
	return &AppError{UnavailableStatus, code, fmt.Sprintf(format, v...), nil, nil}
}

func UnknownError(code, msg string) *AppError {
	return &AppError{UnknownErrorStatus, code, msg, nil, nil}
}

func UnknownErrorf(code, format string, v ...interface{}) *AppError {
	return &AppError{UnknownErrorStatus, code, fmt.Sprintf(format, v...), nil, nil}
}

func Cause(err error) error {
	herr, ok := err.(*AppError)
	if !ok {
		return err
	}

	return herr.Cause()
}

func FromError(err error) (*AppError, bool) {
	herr, ok := err.(*AppError)
	return herr, ok
}

func ErrStatus(err error) int {
	herr, ok := err.(*AppError)
	if !ok {
		return 0
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
