package errors

import (
	"testing"
)

func TestParseError(t *testing.T) {
	err := InternalError("code1", "aaa: bbb")
	perr := ParseError(err.Error())

	if perr.Error() != "code1: aaa: bbb" {
		t.Error("not equal")
	}

	if perr.GetCode() != "code1" {
		t.Error("not equal")
	}

	if perr.Message != "aaa: bbb" {
		t.Error("not equal")
	}
}
