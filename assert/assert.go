package assert

import "testing"

func Equal[V comparable](t *testing.T, got, expected V) {
	if got != expected {
		t.Errorf("%v != %v", got, expected)
	}
}

func IsNil(t *testing.T, got any) {
	if got != nil {
		t.Errorf("%v is not nil", got)
	}
}

func NoErr(t *testing.T, got error) {
	if got != nil {
		t.Errorf("Unexpected error: %e", got)
	}
}

func Err(t *testing.T, got error) {
	if got == nil {
		t.Errorf("Error expected, got nill")
	}
}

func IsNotNil(t *testing.T, got any) {
	if got == nil {
		t.Errorf("%v is nil", got)
	}
}
