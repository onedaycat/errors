package sentry

import (
	errs "errors"
	"os"
	"testing"

	"github.com/onedaycat/errors"
)

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

	xer := errs.New("database connectio lost")

	der := errors.WithCaller(errs.New("gogo")).WithCause(xer)

	err := errors.InternalError("ErrUnableSomething", "Test Sentry Error").WithCaller().WithInput(errors.Input{
		"input1": "vinput1",
	}).WithCause(der)

	p := NewPacket(err)
	p.AddError(err)
	p.AddUser(&User{
		ID: "tester",
	})
	p.AddExtra(Extra{
		"input": err.Input,
	})

	CaptureAndWait(p)
}
