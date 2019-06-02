package errors

type JSONError struct {
    Code       string             `json:"code,omitempty"`
    Message    string             `json:"message,omitempty"`
    ErrType    string             `json:"errType,omitempty"`
    Cause      *JSONError         `json:"cause,omitempty"`
    Stacktrace []*StacktraceFrame `json:"stacktrace,omitempty"`
    Panic      bool               `json:"panic,omitempty"`
    Input      interface{}        `json:"input,omitempty"`
}

func (e *GenericError) JSON() *JSONError {
    jsonErr := &JSONError{
        Code:       e.Code,
        Message:    e.Message,
        ErrType:    e.errType,
        Panic:      e.panic,
        Stacktrace: e.stacktrace,
        Input:      e.input,
    }

    if e.cause != nil {
        jsonErr.Cause = e.cause.JSON()
    }

    return jsonErr
}

func ParseJSONError(jsonErr *JSONError) Error {
    ge := &GenericError{
        Code:       jsonErr.Code,
        Message:    jsonErr.Message,
        errType:    jsonErr.ErrType,
        panic:      jsonErr.Panic,
        input:      jsonErr.Input,
        stacktrace: jsonErr.Stacktrace,
    }

    if jsonErr.Cause != nil {
        ge.cause = ParseJSONError(jsonErr.Cause)
    }

    return ge
}
