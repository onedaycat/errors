package errors

import (
    "errors"
    "fmt"
    "testing"

    "github.com/stretchr/testify/require"
)

func x() error {
    return InternalError("code2", "err2").WithCause(y())
}

func y() error {
    return InternalError("code1", "err1").WithCause(z())
}

func z() error {
    return errors.New("err0")
}

func TestStackTrace(t *testing.T) {
    err1 := InternalError("code1", "err1").WithPanic().WithInput(1)

    fmt.Printf("====\n%+v\n", err1)
    require.NotNil(t, err1.GetStacktrace())
    require.True(t, err1.(*GenericError).panic)
    require.Equal(t, []interface{}{1}, err1.GetAllInputs())
    require.True(t, err1.IsType(InternalErrorType))

    err2 := BadRequest("code2", "err2").WithInput(2).WithCause(err1)
    fmt.Printf("====\n%+v\n", err2)
    require.Equal(t, []interface{}{2, 1}, err2.GetAllInputs())
    require.True(t, err2.(*GenericError).panic)
    require.True(t, err2.IsType(InternalErrorType))
    require.True(t, err2.IsType(BadRequestType))

    err3 := NotFound("code3", "err3").WithInput(3).WithCause(err2)
    fmt.Printf("====\n%+v\n", err3)
    require.Equal(t, []interface{}{3, 2, 1}, err3.GetAllInputs())
    require.True(t, err3.(*GenericError).panic)
    require.True(t, err3.IsType(InternalErrorType))
    require.True(t, err3.IsType(BadRequestType))
    require.True(t, err3.IsType(NotFoundType))

    xxx := errors.New("xxx")
    err4 := Unauthorized("code4", "err4").WithInput(4).WithCause(xxx)
    fmt.Printf("====\n%+v\n", err4)
    require.Equal(t, []interface{}{4}, err4.GetAllInputs())
}

func TestPrintError(t *testing.T) {
    //err := New("code1", "err1").WithCaller()

    result := fmt.Sprintf("%+v\n", x())
    require.Contains(t, result, "code2: err2")
    require.Contains(t, result, "code1: err1")
    require.Contains(t, result, "err0")
    require.Contains(t, result, "x")
    require.Contains(t, result, "y")

    stackString := x().(Error).GetStacktrace().String()
    fmt.Println(stackString)
    require.Contains(t, stackString, "x")
    require.Contains(t, stackString, "y")
}

func TestUtils(t *testing.T) {
    err1 := x()
    err2 := x()
    require.True(t, err1.(Error).Is(err2))
    require.True(t, err1.(Error).Is(err2.(Error).Unwrap()))
    require.True(t, err1.(Error).Is(err2.(Error).Unwrap().(Error).Unwrap()))
    require.Nil(t, err2.(Error).Unwrap().(Error).Unwrap().(Error).Unwrap())
    require.True(t, Is(err1, err2))

    err0 := Unwrap(err1)
    require.EqualError(t, err0, "code1: err1")
    err0 = Unwrap(err0)
    require.EqualError(t, err0, "err0")

    wrapErr0 := Wrap(err0)
    _, ok := wrapErr0.(Error)
    require.True(t, ok)
}

func TestGetErrorWithCause(t *testing.T) {
    result := x().(Error).ErrorWithCause()
    fmt.Println(result)
    require.Contains(t, result, "code2: err2")
    require.Contains(t, result, "code1: err1")
    require.Contains(t, result, "err0")
}

func TestRootError(t *testing.T) {
    err := x().(Error).RootError()
    require.EqualError(t, err, "err0")

    err = InternalError("code1", "err1")
    require.EqualError(t, err.RootError(), "code1: err1")

    err = InternalError("code2", "err2").WithCause(err)
    require.EqualError(t, err.RootError(), "code1: err1")
}
