package util

import (
	"testing"
)

func TestGenerateShaFromDescription(t *testing.T) {
	expected := "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
	description := "test"

	id, _ := Id(description, false)

	if id != expected {
		t.Fatalf("id: %s != %s", id, expected)
	}

}
