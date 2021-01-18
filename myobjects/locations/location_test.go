package locations

import (
	"testing"
)

func TestNew(t *testing.T) {
	l := New(1, 2)
	if l == nil {
		t.Errorf("TestNew = nil; want l{1,2}")
	}

	if l.Row != 1 || l.Col != 2 {
		t.Errorf("TestNew = %d,%d; want l{1,2}", l.Row, l.Col)
	}
}

func TestString(t *testing.T) {
	loc := New(1, 1)
	got := loc.String()
	want := "(1,1)"
	if got != want {
		t.Errorf("TestString = %s; want %s", got, want)
	}
}
