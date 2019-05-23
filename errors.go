package errors

import (
    "io"

    "golang.org/x/exp/errors/fmt"
)

const GenericCode = "Generic"

const (
    NoneType          = ""
    BadRequestType    = "BadRequest"
    UnauthorizedType  = "Unauthorized"
    ForbiddenType     = "Forbidden"
    NotFoundType      = "NotFound"
    TimeoutType       = "Timeout"
    InternalErrorType = "InternalError"
    NotImplementType  = "NotImplement"
)

type Error interface {
    Error() string
    ErrorWithCause() string
    Unwrap() error
    Is(err error) bool
    IsType(errType string) bool
    Format(s fmt.State, verb rune)
    RootError() Error
    IsPanic() bool

    GetCode() string
    GetMessage() string
    GetStacktrace() *Stacktrace
    GetAllInputs() []interface{}
    GetInput() interface{}
    GetType() string

    WithPanic() Error
    WithCause(err error) Error
    WithInput(input interface{}) Error
    WithMessage(msg string) Error
    WithMessagef(format string, args ...interface{}) Error
}

type GenericError struct {
    Code    string `json:"code,omitempty"`
    Message string `json:"message,omitempty"`
    errType string
    cause   Error
    frame   *Stacktrace
    panic   bool
    input   interface{}
}

func (e *GenericError) Error() string {
    if e.Code != "" && e.Code != GenericCode {
        return fmt.Sprintf("%s: %s", e.Code, e.Message)
    }

    return e.Message
}

func (e *GenericError) ErrorWithCause() string {
    s := fmt.Sprintf("%+v\n", e.Error())
    cause := e.cause
    for cause != nil {
        s += fmt.Sprintf("%+v\n", cause.Error())
        xcause := cause.Unwrap()
        if xcause == nil {
            break
        }
        cause = xcause.(Error)
    }

    return s
}

func (e *GenericError) Format(s fmt.State, verb rune) {
    switch verb {
    case 'v':
        _, _ = fmt.Fprintf(s, "%s\n", e.Error())

        if s.Flag('+') {
            if e.frame != nil {
                for _, frame := range e.frame.Frames {
                    _, _ = fmt.Fprintf(s, "%s\t%s:%d\n", frame.Function, frame.AbsolutePath, frame.Lineno)
                }
            }

            cause := e.cause
            for cause != nil {
                _, _ = fmt.Fprintf(s, "\n%s\n", cause.Error())
                if stack := cause.GetStacktrace(); stack != nil {
                    for _, frame := range stack.Frames {
                        _, _ = fmt.Fprintf(s, "%s\t%s:%d\n", frame.Function, frame.AbsolutePath, frame.Lineno)
                    }
                }

                xcause := cause.Unwrap()
                if xcause == nil {
                    break
                }
                cause = xcause.(Error)
            }
        }
    case 's':
        _, _ = io.WriteString(s, e.Error())
        cause := e.cause
        for cause != nil {
            _, _ = io.WriteString(s, e.Error())
            xcause := cause.Unwrap()
            if xcause == nil {
                break
            }
            cause = xcause.(Error)
        }
    default:
        _, _ = fmt.Fprintf(s, "%s", e.Error())
    }
}

func (e *GenericError) Unwrap() error {
    return e.cause
}

func (e *GenericError) WithCause(err error) Error {
    if err == nil {
        return e
    }

    cause, ok := err.(*GenericError)
    if !ok {
        e.cause = &GenericError{
            Code:    GenericCode,
            Message: err.Error(),
            errType: InternalErrorType,
            frame:   NewStacktrace(1),
        }

        return e
    }

    e.cause = cause
    e.frame.Frames = e.frame.Frames[len(e.frame.Frames)-1:]
    e.panic = cause.panic

    return e
}

func (e *GenericError) WithPanic() Error {
    e.panic = true

    return e
}

func (e *GenericError) WithInput(input interface{}) Error {
    e.input = input

    return e
}

func (e *GenericError) WithMessage(msg string) Error {
    e.Message = msg

    return e
}

func (e *GenericError) WithMessagef(format string, args ...interface{}) Error {
    e.Message = fmt.Sprintf(format, args...)

    return e
}

func (e *GenericError) GetCode() string {
    return e.Code
}

func (e *GenericError) GetStacktrace() *Stacktrace {
    return e.frame
}

func (e *GenericError) GetAllInputs() []interface{} {
    inputs := make([]interface{}, 0, 5)
    if e.input != nil {
        inputs = append(inputs, e.input)
    }

    cause := e.cause
    for cause != nil {
        cerr := cause.(*GenericError)
        if cerr.input != nil {
            inputs = append(inputs, cerr.input)
        }
        xcause := cerr.Unwrap()
        if xcause == nil {
            break
        }
        cause = xcause.(Error)
    }

    return inputs
}

func (e *GenericError) GetInput() interface{} {
    return e.input
}

func (e *GenericError) GetMessage() string {
    return e.Message
}

func (e *GenericError) GetType() string {
    return e.errType
}

func (e *GenericError) Is(err error) bool {
    xerr, ok := err.(Error)
    if !ok {
        return false
    }

    if e.Code == xerr.GetCode() {
        return true
    }

    cause := e.cause
    for cause != nil {
        cerr := cause.(*GenericError)
        if cerr.Code == xerr.GetCode() {
            return true
        }
        xcause := cerr.Unwrap()
        if xcause == nil {
            break
        }
        cause = xcause.(*GenericError)
    }

    return false
}

func (e *GenericError) IsType(errType string) bool {
    if e.errType == e.errType {
        return true
    }

    cause := e.cause
    for cause != nil {
        cerr := cause.(*GenericError)
        if cerr.errType == cerr.errType {
            return true
        }
        xcause := cerr.Unwrap()
        if xcause == nil {
            return false
        }
        cause = xcause.(Error)
    }

    return false
}

func (e *GenericError) RootError() Error {
    cause := e.cause
    for cause != nil {
        xcause := cause.Unwrap()
        if xcause == nil {
            return cause
        }
        cause = xcause.(Error)
    }

    return e
}

func (e *GenericError) IsPanic() bool {
    return e.panic
}

func New(msg string) Error {
    return &GenericError{
        Code:    GenericCode,
        Message: msg,
        frame:   NewStacktrace(1),
    }
}

func NewWithCode(code, msg string) Error {
    return &GenericError{
        Code:    code,
        Message: msg,
        frame:   NewStacktrace(1),
    }
}

func NewWithTypeAndCode(errType, code, msg string) Error {
    return &GenericError{
        Code:    code,
        Message: msg,
        errType: errType,
        frame:   NewStacktrace(1),
    }
}

func BadRequest(code, msg string) Error {
    return &GenericError{
        Code:    code,
        Message: msg,
        errType: BadRequestType,
        frame:   NewStacktrace(1),
    }
}

func Unauthorized(code, msg string) Error {
    return &GenericError{
        Code:    code,
        Message: msg,
        errType: UnauthorizedType,
        frame:   NewStacktrace(1),
    }
}

func Forbidden(code, msg string) Error {
    return &GenericError{
        Code:    code,
        Message: msg,
        errType: ForbiddenType,
        frame:   NewStacktrace(1),
    }
}

func NotFound(code, msg string) Error {
    return &GenericError{
        Code:    code,
        Message: msg,
        errType: NotFoundType,
        frame:   NewStacktrace(1),
    }
}

func Timeout(code, msg string) Error {
    return &GenericError{
        Code:    code,
        Message: msg,
        errType: TimeoutType,
        frame:   NewStacktrace(1),
    }
}

func InternalError(code, msg string) Error {
    return &GenericError{
        Code:    code,
        Message: msg,
        errType: InternalErrorType,
        frame:   NewStacktrace(1),
    }
}

func NotImplement(code, msg string) Error {
    return &GenericError{
        Code:    code,
        Message: msg,
        errType: NotImplementType,
        frame:   NewStacktrace(1),
    }
}
