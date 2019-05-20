package errors

import (
    "fmt"
    "testing"

    "github.com/stretchr/testify/require"
)

func xx() error {
    def := DefInternalError("code1", "err1")
    return def.New()
}

func TestPrintDefError(t *testing.T) {
    result := fmt.Sprintf("%+v\n", xx())
    fmt.Println(result)
    require.Contains(t, result, "code1: err1")
    require.Contains(t, result, "xx")
}

func TestHttpStatus(t *testing.T) {
    require.Equal(t, 500, HttpStatus(DefInternalError("code1", "err1").New().GetType()))
    require.Equal(t, 404, HttpStatus(DefNotFound("code1", "err1").New().GetType()))
    require.Equal(t, 400, HttpStatus(DefBadRequest("code1", "err1").New().GetType()))
    require.Equal(t, 403, HttpStatus(DefForbidden("code1", "err1").New().GetType()))
    require.Equal(t, 501, HttpStatus(DefNotImplement("code1", "err1").New().GetType()))
    require.Equal(t, 441, HttpStatus(DefTimeout("code1", "err1").New().GetType()))
    require.Equal(t, 401, HttpStatus(DefUnauthorized("code1", "err1").New().GetType()))
}
