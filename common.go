package errors

func Is(err, target error) bool {
    u, ok := err.(Error)
    if !ok {
        if target == nil {
            return err == target
        }
        for {
            if err == target {
                return true
            }
            if x, ok := err.(interface{ Is(error) bool }); ok && x.Is(target) {
                return true
            }
            if err = Unwrap(err); err == nil {
                return false
            }
        }
    }

    return u.Is(target)
}

func Unwrap(err error) error {
    u, ok := err.(Error)
    if !ok {
        return nil
    }
    return u.Unwrap()
}

func Wrap(err error) Error {
    return &GenericError{
        Code:       GenericCode,
        Message:    err.Error(),
        stacktrace: NewStacktrace(1),
    }
}
