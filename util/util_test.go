package util

import (
	"testing"
)

func TestGenerateShaFromDescription(t *testing.T) {
	expected := "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
	description := "test"

	id, err := Id(description, false)

	if err != nil || id != expected {
		t.Fail()
	}

}

func TestReturnErrIfEmptyDescription(t *testing.T) {

	_, err := Id("", false)
	if err == nil {
		t.Fail()
	}

	_, err = Id(" \t", false)
	if err == nil {
		t.Fail()
	}
}

func TestReturnIdAsIsIfIsRawFlagSet(t *testing.T) {
	expected := "test"

	id, err := Id(expected, true)

	if err != nil || id != expected {
		t.Fail()
	}
}
