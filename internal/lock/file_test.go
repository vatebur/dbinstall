package lock

import (
	"errors"
	"testing"
)

func TestAcquireIsExclusive(t *testing.T) {
	directory := t.TempDir()
	first, err := Acquire(directory)
	if err != nil {
		t.Fatal(err)
	}
	defer first.Close()
	if _, err := Acquire(directory); !errors.Is(err, ErrHeld) {
		t.Fatalf("second lock error = %v", err)
	}
}
