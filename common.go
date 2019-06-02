package errors

import (
    "golang.org/x/exp/errors"
)

func As(err error, target interface{}) bool {
    return errors.As(err, target)
}

func Is(err, target error) bool {
    return errors.Is(err, target)
}

func Opaque(err error) error {
    return errors.Opaque(err)
}

func Unwrap(err error) error {
    return errors.Unwrap(err)
}

func Wrap(err error) Error {
    return &GenericError{
        Code:       GenericCode,
        Message:    err.Error(),
        stacktrace: NewStacktrace(1),
    }
}
