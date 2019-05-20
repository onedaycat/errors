package sentry

import (
    errs "errors"
    "os"
    "testing"

    "github.com/onedaycat/errors"
)

func x() errors.Error {
    xer := errs.New("database connectio lost")

    return errors.Wrap(errs.New("Unable to apply events")).WithCause(xer).WithInput(map[string]interface{}{
        "id":   1,
        "name": "tester",
    })
}

func TestCaptureAndWait(t *testing.T) {
    SetDSN(os.Getenv("SENTRY_DSN"))
    SetOptions(
        WithEnv("dev"),
        WithRelease("1.0.0"),
        WithServerName("Test Sentry"),
        WithServiceName("test1"),
        WithTags(Tags{
            {"tag1", "vtag1"},
        }),
        WithDefaultExtra(Extra{
            "dextra1": "dvextra1",
        }),
    )

    err := x()

    err = errors.InternalError("ErrUnableSomething", "Test Sentry Error").WithInput(map[string]interface{}{
        "input1": "vinput1",
    }).WithCause(err)

    p := NewPacket(err)
    p.SetUser(&User{
        ID: "tester",
    })

    CaptureAndWait(p)
}
