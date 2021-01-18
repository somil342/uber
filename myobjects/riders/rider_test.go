package riders

import "testing"

func TestNew(t *testing.T) {
	rider := New(1, "abhay")
	if rider == nil {
		t.Errorf("TestNew = nil; want non nil")
	}
}
