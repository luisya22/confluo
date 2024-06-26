package assert

import (
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}

func NotEqual[T comparable](t *testing.T, actual, notExpected T) {
	t.Helper()

	if actual == notExpected {
		t.Errorf("got: %v; doesn't want: %v", actual, notExpected)
	}
}

func StringContains(t *testing.T, actual, expectedSubstring string) {
	t.Helper()

	if !strings.Contains(actual, expectedSubstring) {
		t.Errorf("got: %v; expected to contain: %v", actual, expectedSubstring)
	}
}

func NilError(t *testing.T, actual error) {
	t.Helper()

	if actual != nil {
		t.Errorf("got: %v; expected: nil", actual)
	}
}

func Error(t *testing.T, actual error) {
	t.Helper()

	if actual == nil {
		t.Errorf("got: %v; expected: not nil", actual)
	}
}
