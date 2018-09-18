package sentry

import (
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
		WithTags(Tags{
			{"tag1", "vtag1"},
		}),
		WithDefaultExtra(Extra{
			"dextra1": "dvextra1",
		}),
	)

	err := errors.InternalError("1000", "Test Sentry Error").WithCaller().WithInput(errors.Input{
		"input1": "vinput1",
	})

	p := NewPacket(err)
	p.AddStackTrace(err.StackTrace())
	p.AddUser(&User{
		ID: "tester",
	})
	p.AddExtra(Extra{
		"input": err.Input,
	})

	CaptureAndWait(p)
}
