package assert

import (
	"fmt"
	"testing"
)

// Equal asserts that two values are equal
func Equal[T comparable](t *testing.T, expected, actual T, messages ...any) {
	if expected != actual {
		t.Error(fmt.Sprintf("expected %v but received %v", expected, actual), messages)
	}
}

// NoError asserts that there was no error
func NoError(t *testing.T, err error, messages ...any) {
	if err != nil {
		t.Error("received unexpected error", messages)
	}
}
