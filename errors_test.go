package errors

import (
	"fmt"
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

func TestFormat(t *testing.T) {
	err1 := InternalError("c1", "e1").WithCaller()
	err2 := InternalError("c2", "e2").WithCause(err1).WithCaller()
	err3 := InternalError("c3", "e3").WithCause(err2).WithCaller()

	fmt.Printf("%+v\n", err3)
}
