package errors

import "fmt"

type ErrorDefinition struct {
    Code    string
    Type    string
    Message string
}

func Def(errType, code string, msg ...string) *ErrorDefinition {
    e := &ErrorDefinition{
        Code: code,
        Type: errType,
    }

    if len(msg) > 0 {
        e.Message = msg[0]
    }

    return e
}

func (e *ErrorDefinition) New(msg ...string) Error {
    if len(msg) > 0 {
        return &GenericError{
            Code:    e.Code,
            Message: msg[0],
            errType: e.Type,
            frame:   NewStacktrace(1),
        }
    }

    return &GenericError{
        Code:    e.Code,
        Message: e.Message,
        errType: e.Type,
        frame:   NewStacktrace(1),
    }
}

func (e *ErrorDefinition) Newf(format string, v ...interface{}) Error {
    return &GenericError{
        Code:    e.Code,
        Message: fmt.Sprintf(format, v...),
        errType: e.Type,
        frame:   NewStacktrace(1),
    }
}

func (e *ErrorDefinition) Is(err Error) bool {
    return err != nil && e.Code == err.GetCode()
}

func DefBadRequest(code string, msg ...string) *ErrorDefinition {
    e := &ErrorDefinition{
        Code: code,
        Type: BadRequestType,
    }
    if len(msg) > 0 {
        e.Message = msg[0]
    }

    return e
}

func DefUnauthorized(code string, msg ...string) *ErrorDefinition {
    e := &ErrorDefinition{
        Code: code,
        Type: UnauthorizedType,
    }
    if len(msg) > 0 {
        e.Message = msg[0]
    }

    return e
}

func DefForbidden(code string, msg ...string) *ErrorDefinition {
    e := &ErrorDefinition{
        Code: code,
        Type: ForbiddenType,
    }
    if len(msg) > 0 {
        e.Message = msg[0]
    }

    return e
}

func DefNotFound(code string, msg ...string) *ErrorDefinition {
    e := &ErrorDefinition{
        Code: code,
        Type: NotFoundType,
    }
    if len(msg) > 0 {
        e.Message = msg[0]
    }

    return e
}

func DefTimeout(code string, msg ...string) *ErrorDefinition {
    e := &ErrorDefinition{
        Code: code,
        Type: TimeoutType,
    }
    if len(msg) > 0 {
        e.Message = msg[0]
    }

    return e
}

func DefInternalError(code string, msg ...string) *ErrorDefinition {
    e := &ErrorDefinition{
        Code: code,
        Type: InternalErrorType,
    }
    if len(msg) > 0 {
        e.Message = msg[0]
    }

    return e
}

func DefNotImplement(code string, msg ...string) *ErrorDefinition {
    e := &ErrorDefinition{
        Code: code,
        Type: NotImplementType,
    }
    if len(msg) > 0 {
        e.Message = msg[0]
    }

    return e
}
