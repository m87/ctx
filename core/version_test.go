package core

import (
	"reflect"
	"testing"
)

func TestVersionSort(t *testing.T) {
	versions := []Version{
		{Major: 1, Minor: 0, Patch: 0},
		{Major: 1, Minor: 0, Patch: 1},
		{Major: 2, Minor: 0, Patch: 0},
		{Major: 1, Minor: 1, Patch: 0},
	}
	Sort(versions)

	expected := []Version{
		{Major: 1, Minor: 0, Patch: 0},
		{Major: 1, Minor: 0, Patch: 1},
		{Major: 1, Minor: 1, Patch: 0},
		{Major: 2, Minor: 0, Patch: 0},
	}

	if !reflect.DeepEqual(versions, expected) {
		t.Errorf("Expected %v, got %v", expected, versions)
	}
}
